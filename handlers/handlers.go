package handlers

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"hammer-down-server/db"
	"hammer-down-server/utils"
	"hash/crc32"
	"io"
	"log"
	"net/http"
)

func decompressData(compressedData []byte) ([]byte, error) {
	// Create a reader for the compressed data
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer reader.Close()

	// Decompress the data
	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	return decompressedData.Bytes(), nil
}

func DetectCheat(w http.ResponseWriter, r *http.Request, database *db.DB) {

	log.Printf("Recv req from HD client")

	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	data, err := io.ReadAll(r.Body)
	var compLen int = len(data)
	if len(data) == 0 {
		log.Printf("Data is empty!")
		http.Error(w, "Data is empty", http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	data, err = decompressData(data)
	if err != nil {
		log.Printf("Failed to decompress data: %v", err)

		http.Error(w, "Failed to decompress page", http.StatusInternalServerError)
		return
	}
	log.Printf("Starting to analize page... Page comress ratio %d", len(data)/compLen)
	crc := crc32.ChecksumIEEE(data)
	fileHash := fmt.Sprintf("%08x", crc)

	cheat, err := database.FindCheatByHash(fileHash)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if cheat != nil {

		response := map[string]interface{}{
			"detected": true,
			"method":   "hash",
			"cheat":    cheat,
		}
		sendJSON(w, response)
		return
	}

	signatures, err := database.LoadAllSignatures()
	if err != nil {
		http.Error(w, "Error loading signatures", http.StatusInternalServerError)
		return
	}

	for _, sig := range signatures {
		matched, err := utils.MatchSignature(data, sig.SignaturePattern)
		if err != nil {
			http.Error(w, "Error matching signature", http.StatusInternalServerError)
			return
		}
		if !matched {
			continue
		}

		// Find the cheat info for this signature
		cheat, err := database.FindCheatByID(sig.CheatID)
		if err != nil || cheat == nil {
			http.Error(w, "Cheat not found for matched signature", http.StatusInternalServerError)
			return
		}

		// Optionally: insert the new hash into cheat_hashes for future quick detection
		if err := database.InsertCheatHash(sig.CheatID, fileHash, "Detected via signature and cached"); err != nil {
			log.Printf("Warning: failed to insert cheat hash cache: %v", err)
		}

		response := map[string]interface{}{
			"detected":   true,
			"method":     "signature",
			"cheat_name": cheat.Name,
		}
		sendJSON(w, response)
		return
	}

	sendJSON(w, map[string]interface{}{
		"detected": false,
		"message":  "No cheat detected",
	})
}

func sendJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
	io.Copy(w, bytes.NewReader(jsonBytes))
}

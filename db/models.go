package db

// Cheat represents a row from the `cheats` table.
type Cheat struct {
	CheatID   int    `json:"cheat_id"`
	Name      string `json:"cheat_name"`
	Desc      string `json:"cheat_description"`
	DateAdded string `json:"date_added"`
}

// Signature represents a row from the `signatures` table.
type Signature struct {
	SignatureID      int    `json:"signature_id"`
	CheatID          int    `json:"cheat_id"`
	SignaturePattern string `json:"signature_pattern"`
	Description      string `json:"description"`
	DateAdded        string `json:"date_added"`
}

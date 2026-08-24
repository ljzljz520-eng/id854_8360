package wardrobe

import "encoding/json"

func decodeLedgerJSON(data []byte, target any) error { return json.Unmarshal(data, target) }

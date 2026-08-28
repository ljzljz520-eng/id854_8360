package rehearsal

import "encoding/json"

func decodeRosterJSON(data []byte, target any) error { return json.Unmarshal(data, target) }

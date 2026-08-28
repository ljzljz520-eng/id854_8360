package show

import "encoding/json"

func decodeOrderJSON(data []byte, target any) error { return json.Unmarshal(data, target) }

package rehearsal

import "encoding/json"

func jsonDecode(data []byte, target any) error { return json.Unmarshal(data, target) }

package audit

import (
	"fmt"
	"sort"
	"theatrecontrol/internal/model"
)

func encodeEntry(entry model.AuditEntry) ([]byte, error)     { return marshal(entry) }
func decodeEntry(data []byte, entry *model.AuditEntry) error { return unmarshal(data, entry) }

func summarize(entries []model.AuditEntry) string {
	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		parts = append(parts, fmt.Sprintf("%s:%s", entry.Action, entry.EntityID))
	}
	sort.Strings(parts)
	return fmt.Sprintf("%d entries [%s]", len(parts), join(parts))
}

func join(values []string) string {
	result := ""
	for i, value := range values {
		if i > 0 {
			result += ","
		}
		result += value
	}
	return result
}

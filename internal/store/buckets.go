package store

var bucketNames = [][]byte{
	[]byte("roles"), []byte("rehearsals"), []byte("costumes"), []byte("allocations"),
	[]byte("performances"), []byte("ban_rules"), []byte("audit_entries"), []byte("assignments"), []byte("ticket_orders"),
}

func bucketFor(kind string) []byte {
	switch kind {
	case "role":
		return bucketNames[0]
	case "rehearsal":
		return bucketNames[1]
	case "costume":
		return bucketNames[2]
	case "allocation":
		return bucketNames[3]
	case "performance":
		return bucketNames[4]
	case "ban_rule":
		return bucketNames[5]
	case "audit":
		return bucketNames[6]
	case "assignment":
		return bucketNames[7]
	default:
		return bucketNames[8]
	}
}

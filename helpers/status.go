package helpers

func StatusString(status int) string {
	switch status {
	case 0:
		return "Waiting List"
	case 1:
		return "In Progress"
	case 2:
		return "Done"
	default:
		return "Deleted"
	}
}

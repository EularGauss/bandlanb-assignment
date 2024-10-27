package s3service

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// getContentType returns the MIME type based on the file extension
func getContentType(extension string) string {
	switch extension {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".bmp":
		return "image/bmp"
	case ".png":
		return "image/png"
	default:
		return ""
	}
}

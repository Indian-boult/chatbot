package model

// ImageInfo represents the metadata for an image uploaded by the user.
type ImageInfo struct {
	ID         int64  // Unique identifier for the image
	Identifier string // User-provided identifier for the image
	URL        string // URL where the image is stored
}

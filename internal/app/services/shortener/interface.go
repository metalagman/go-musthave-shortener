package shortener

// Store of the urls served by shortener
type Store interface {
	// WriteURL to storage
	WriteURL(url string) (string, error)
	// ReadURL from storage
	ReadURL(id string) (string, error)
}

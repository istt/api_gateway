package interfaces

// ObjectStorage provide basic function for upload and retrieve object based on ID
type ObjectStorage interface {

	// Put update the object with given key
	Put(key string, content []byte) error

	// Get return the value of the key
	Get(key string) ([]byte, error)

	// Del delete the value of the key
	Del(key string) error
}

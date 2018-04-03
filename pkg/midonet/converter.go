package midonet

type APIResource struct {
	PathForPost   string
	PathForPut    string
	PathForDelete string
	MediaType     string
	Body          interface{}
}

type Converter interface {
	// Convert
	// - if nil obj is given, only PathForDelete fields for the
	//   APIResource returned are valid.
	Convert(key string, obj interface{}, confing *Config) ([]*APIResource, error)
}

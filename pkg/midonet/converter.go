package midonet

type body interface {
	MediaType() string
}

type APIResource struct {
	PathForPost   string
	PathForPut    string
	PathForDelete string
	Body          body
}

type Converter interface {
	// Convert
	// - if nil obj is given, only PathForDelete fields for the
	//   APIResource returned are valid.
	Convert(key string, obj interface{}, confing *Config) ([]*APIResource, error)
}

package midonet

type APIResource interface {
	Path(string) string
	MediaType() string
}

type Converter interface {
	// Convert
	// - if nil obj is given, only PathForDelete fields for the
	//   APIResource returned are valid.
	Convert(key string, obj interface{}, config *Config) ([]APIResource, SubResourceMap, error)
}

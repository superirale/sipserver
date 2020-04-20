package uri

type URI struct {
	Uri     string
	Tags    map[string]string
	UriText string
}

func BuildURI(uri, uriText string, tags map[string]string) *URI {
	uriObj := new(URI)
	uriObj.Uri = uri
	uriObj.UriText = uriText
	uriObj.Tags = tags

	return uriObj
}
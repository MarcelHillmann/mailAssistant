package account

// AttachmentPromise covers a attachment
type AttachmentPromise struct {
	body   []byte
	params map[string]string
}

// GetContentType returns the content-type header
func (promise AttachmentPromise) GetContentType() string {
	if ret, ok := promise.params["content-type"]; ok {
		return ret
	}
	return ""
}

// GetContentDisposition returns the content-disposition header
func (promise AttachmentPromise) GetContentDisposition() string {
	if ret, ok := promise.params["content-disposition"]; ok {
		return ret
	}
	return ""
}

// GetFilename returns the file name
func (promise AttachmentPromise) GetFilename() string {
	if ret, ok := promise.params["filename"]; ok {
		return ret
	}
	return ""
}

// Body returns the attachment body as byte array
func (promise AttachmentPromise) Body() []byte {
	return promise.body
}

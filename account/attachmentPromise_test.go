package account

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAttachmentPromise(t *testing.T) {
	t.Run("GetContentDisposition", attachmentPromiseGetContentDisposition)
	t.Run("GetContentType", attachmentPromiseGetContentType)
	t.Run("GetFilename", attachmentPromiseGetFilename)
	t.Run("writer", attachmentPromiseWriter)
}

func attachmentPromiseGetContentDisposition(t *testing.T) {
	param := make(map[string]string)
	param["content-disposition"] = "test"
	promise := AttachmentPromise{make([]byte, 0), param}

	require.Equal(t, "test", promise.GetContentDisposition())
	require.Empty(t, promise.GetContentType())
	require.Empty(t, promise.GetFilename())
	require.Len(t, promise.Body(), 0)
}

func attachmentPromiseGetContentType(t *testing.T) {
	param := make(map[string]string)
	param["content-type"] = "test"
	promise := AttachmentPromise{make([]byte, 0), param}

	require.Equal(t, "test", promise.GetContentType())
	require.Empty(t, promise.GetContentDisposition())
	require.Empty(t, promise.GetFilename())
	require.Len(t, promise.Body(), 0)
}

func attachmentPromiseGetFilename(t *testing.T) {
	param := make(map[string]string)
	param["filename"] = "test"
	promise := AttachmentPromise{make([]byte, 0), param}

	require.Equal(t, "test", promise.GetFilename())
	require.Empty(t, promise.GetContentType())
	require.Empty(t, promise.GetContentDisposition())
	require.Len(t, promise.Body(), 0)
}

func attachmentPromiseWriter(t *testing.T) {
	body := []byte("hallo")
	promise := AttachmentPromise{body, make(map[string]string)}

	require.Len(t, promise.Body(), 5)
	require.Equal(t, []byte("hallo"), promise.Body())
	require.Empty(t, promise.GetFilename())
	require.Empty(t, promise.GetContentType())
	require.Empty(t, promise.GetContentDisposition())
}

package harmony

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type multipartPayload interface {
	json() ([]byte, error)
}

// multipartFromFiles generate a multipart body given a payload and some files.
// It returns the raw generated body along a header with the proper Content-Type value set.
func multipartFromFiles(payload multipartPayload, files ...File) ([]byte, http.Header, error) {
	// Underlying buffer the multipart body will be written to.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Send the endpoint parameters as JSON in a the "payload_json" part.
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", "application/json")
	pw, err := w.CreatePart(h)
	if err != nil {
		return nil, nil, err
	}

	bytes, err := payload.json()
	if err != nil {
		return nil, nil, err
	}
	if _, err = pw.Write(bytes); err != nil {
		return nil, nil, err
	}

	// Create a new part for each file.
	for i, f := range files {
		cd := fmt.Sprintf(`form-data; name="file%d"; filename="%s"`, i, f.Name)

		h = textproto.MIMEHeader{}
		h.Set("Content-Disposition", cd)
		h.Set("Content-Type", "application/octet-stream")

		pw, err = w.CreatePart(h)
		if err != nil {
			return nil, nil, err
		}

		if _, err = io.Copy(pw, f.Reader); err != nil {
			return nil, nil, err
		}
	}

	if err = w.Close(); err != nil {
		return nil, nil, err
	}

	header := http.Header{}
	header.Set("Content-Type", w.FormDataContentType())
	return b.Bytes(), header, nil
}

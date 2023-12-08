package handler

import (
	"bytes"
	"io"
	"net/http"
)

type CustomTransport struct {
	Transport http.RoundTripper
}

const message = "404 page not found"

func (c *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		resp = &http.Response{
			Status:        "Internal Server Error",
			StatusCode:    500,
			Proto:         req.Proto,
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          io.NopCloser(bytes.NewBufferString(message)),
			ContentLength: int64(len(message)),
			Request:       req,
			Header:        make(http.Header, 0),
		}

		return resp, nil

	}

	return resp, err
}

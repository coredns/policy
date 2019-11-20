package response

import (
	"github.com/miekg/dns"
)

// Reader implements ResponseWriter and exposes the message of the response.
type Reader struct {
	dns.ResponseWriter
	Msg   *dns.Msg
}

// NewReader returns a new Reader
func NewReader(w dns.ResponseWriter) *Reader {
	return &Reader{
		ResponseWriter: w,
		Msg:            nil,
	}
}

// WriteMsg overrides ResponseWriter.WriteMsg
func (r *Reader) WriteMsg(response *dns.Msg) error {
	r.Msg = response
	return r.ResponseWriter.WriteMsg(response)
}
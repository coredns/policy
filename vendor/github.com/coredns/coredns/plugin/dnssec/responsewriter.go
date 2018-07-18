package dnssec

import (
	"log"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// ResponseWriter sign the response on the fly.
type ResponseWriter struct {
	dns.ResponseWriter
	d Dnssec
}

// WriteMsg implements the dns.ResponseWriter interface.
func (d *ResponseWriter) WriteMsg(res *dns.Msg) error {
	// By definition we should sign anything that comes back, we should still figure out for
	// which zone it should be.
	state := request.Request{W: d.ResponseWriter, Req: res}

	zone := plugin.Zones(d.d.zones).Matches(state.Name())
	if zone == "" {
		return d.ResponseWriter.WriteMsg(res)
	}
	state.Zone = zone

	if state.Do() {
		res = d.d.Sign(state, time.Now().UTC())

		cacheSize.WithLabelValues("signature").Set(float64(d.d.cache.Len()))
	}
	state.SizeAndDo(res)

	return d.ResponseWriter.WriteMsg(res)
}

// Write implements the dns.ResponseWriter interface.
func (d *ResponseWriter) Write(buf []byte) (int, error) {
	log.Printf("[WARNING] Dnssec called with Write: not signing reply")
	n, err := d.ResponseWriter.Write(buf)
	return n, err
}

// Hijack implements the dns.ResponseWriter interface.
func (d *ResponseWriter) Hijack() { d.ResponseWriter.Hijack() }

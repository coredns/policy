package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/infobloxopen/themis/pdp"
)

const (
	queryCmd          = "query"
	missingStorageMsg = `"Server missing policy storage"`
)

type storageHandler struct {
	s *Server
}

func (handler *storageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		depth uint64
		err   error
	)
	path := strings.FieldsFunc(r.URL.Path, func(c rune) bool { return c == '/' })
	if len(path) == 0 || path[0] != queryCmd {
		http.Error(w, "404 page not found", 404)
	}

	// parse depth
	queryOpt := r.URL.Query()
	if depthOpt, ok := queryOpt["depth"]; ok {
		depthStr := depthOpt[0]
		depth, err = strconv.ParseUint(depthStr, 10, 31)
		if err != nil {
			http.Error(w, strconv.Quote(err.Error()), 400)
			return
		}
	}

	// sanity check
	root := handler.s.p
	if root == nil {
		http.Error(w, missingStorageMsg, 404)
		return
	}

	// parse path
	path = path[1:] // remove queryCmd
	target, err := root.GetAtPath(path)
	if err != nil {
		var errCode int
		if _, ok := err.(*pdp.PathNotFoundError); ok {
			errCode = 404
		} else {
			errCode = 500
		}
		http.Error(w, strconv.Quote(err.Error()), errCode)
		return
	}

	// dump
	if err = target.MarshalWithDepth(w, int(depth)); err != nil {
		http.Error(w, strconv.Quote(err.Error()), 500)
		return
	}
}

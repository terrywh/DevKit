package main

import (
	"net/http"
	"path/filepath"
	"strings"
)


type FileServer struct {
	Root string
	Mime map[string]string
}

func (fs FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	if ctype, ok := fs.Mime[filepath.Ext(r.URL.Path)]; ok {
		w.Header().Set("Content-Type", ctype)
	}
	http.ServeFile(w, r, filepath.Clean(filepath.Join(fs.Root, r.URL.Path)))
}
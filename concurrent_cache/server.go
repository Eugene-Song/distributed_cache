package cache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_cache/"

// HTTPPool is a struct for the pool of HTTP peers
type HttpPool struct {
	// peer's address, including the port
	self string

	// the base path for the cache service
	basePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handles the HTTP request
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check if the path is correct for cache service
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)

	}

	p.Log("%s %s", r.Method, r.URL.Path)

	// required path format: /<basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	// check if the path is correct foramtted
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write the data to the response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

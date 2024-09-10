package cache

import (
	"distributed_cache/consistenthash"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

const defaultBasePath = "/_cache/"
const defaultReplicas = 50

// HTTPPool is a struct for the pool of HTTP peers
type HttpPool struct {
	// peer's address, including the port
	self string

	// the base path for the cache service
	basePath string

	// getting the peers
	peers *consistenthash.ConsistentHashMap

	// protect the httpGetters
	mu sync.Mutex

	// map of the httpGetters
	httpGetters map[string]*httpGetter
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// helper function to log the server
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

// Set updates the pool's list of peers
func (p *HttpPool) SetHttpPool(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consistenthash.NewConsistentHashMap(defaultReplicas, nil)
	p.peers.AddNode(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer picks a peer according to the key
func (p *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// get peer from consistent hash
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}

	return nil, false
}

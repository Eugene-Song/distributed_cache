package cache

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// choose a peer to get the data
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// get the data from the peer
type PeerGetter interface {
	GetFromPeer(group string, key string) ([]byte, error)
}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) GetFromPeer(group string, key string) ([]byte, error) {
	// use queryescape to clean the group and key
	u := fmt.Sprintf("%v/%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	// bytes, err := ioutil.ReadAll(resp.Body)
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)

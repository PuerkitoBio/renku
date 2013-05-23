package cache

import (
	"bytes"
	"container/list"
	"io"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/ghost/handlers"
	"github.com/PuerkitoBio/purell"
)

type lruCache struct {
	sz int // a size of 0 means no limit, cache everything

	mu sync.Mutex // lock to protect the following fields
	l  *list.List
	m  map[string]*list.Element
}

type CacheableItem interface {
	Key() string
}

func newLRUCache(sz int) *lruCache {
	return &lruCache{
		sz: sz,
		l:  list.New(),
		m:  make(map[string]*list.Element, sz),
	}
}

func (ø *lruCache) get(k string) (CacheableItem, bool) {
	ø.mu.Lock()
	defer ø.mu.Unlock()
	e, ok := ø.m[k]
	if !ok {
		// Not in cache
		return nil, false
	}
	// Put back on top, this is the MRU
	ø.l.MoveToFront(e)
	return e.Value.(CacheableItem), true
}

func (ø *lruCache) set(ci CacheableItem) {
	ø.mu.Lock()
	defer ø.mu.Unlock()
	// Ensure the element does not already exist, avoid creating duplicates in the list
	k := ci.Key()
	if e, ok := ø.m[k]; ok {
		ø.l.MoveToFront(e)
		return
	}
	e := ø.l.PushFront(ci)
	ø.m[k] = e
	for ø.l.Len() > ø.sz && ø.sz > 0 {
		// The tail (LRU) must be dropped
		e := ø.l.Back()
		backCi := e.Value.(CacheableItem)
		delete(ø.m, backCi.Key())
	}
}

type responseCacheItem struct {
	buf  *bytes.Buffer
	hdr  http.Header
	nurl string
}

func (ø responseCacheItem) Key() string {
	return ø.nurl
}

type cacheWriter struct {
	http.ResponseWriter
	code  int
	multi io.Writer
}

func (ø *cacheWriter) WriteHeader(code int) {
	ø.code = code
	ø.ResponseWriter.WriteHeader(code)
}

func (ø *cacheWriter) Write(b []byte) (int, error) {
	return ø.multi.Write(b)
}

func copyHeader(dst, src http.Header) {

}

func LRUCacheHandler(h http.Handler, cacheSz int, normFlags purell.NormalizationFlags) http.Handler {
	// Create the cache
	c := newLRUCache(cacheSz)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := getCacheWriter(w); ok {
			// Self-awareness
			h.ServeHTTP(w, r)
			return
		}
		nUrl, err := purell.NormalizeURLString(r.URL.String(), normFlags)
		if err != nil {
			// Impossible - means the URL string passed as input could not be re-parsed as URL instance,
			// so just disengage the cache for this request.
			h.ServeHTTP(w, r)
			return
		}
		if ci, ok := c.get(nUrl); ok {
			// Return cached content
			item := ci.(*responseCacheItem)
			copyHeader(item.hdr, w.Header())
			w.Write(item.buf.Bytes())
			return
		}

		// Create the cache writer to store the response
		buf := bytes.NewBuffer(nil)
		cw := &cacheWriter{
			w,
			http.StatusOK,
			io.MultiWriter(w, buf),
		}

		// Call the wrapped handler with the cache writer
		h.ServeHTTP(cw, r)

		if cw.code >= 200 && cw.code < 300 {
			item := &responseCacheItem{
				buf:  buf,
				hdr:  cw.Header(),
				nurl: nUrl,
			}
			c.set(item)
		}
	})
}

func getCacheWriter(w http.ResponseWriter) (*cacheWriter, bool) {
	cw, ok := handlers.GetResponseWriter(w, func(tst http.ResponseWriter) bool {
		_, ok := tst.(*cacheWriter)
		return ok
	})
	if !ok {
		return nil, false
	}
	return cw.(*cacheWriter), true
}

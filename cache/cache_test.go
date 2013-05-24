package cache

import (
	"github.com/bmizerany/assert"
	"testing"
)

type cachedItem struct {
	k string
}

func (ø *cachedItem) Key() string {
	return ø.k
}

func TestCacheSet(t *testing.T) {
	c := newLRUCache(1)
	c.set(&cachedItem{"a"})
	assert.T(t, c.l.Len() == 1)
}

func TestCacheGet(t *testing.T) {
	c := newLRUCache(1)
	ci := &cachedItem{"a"}
	c.set(ci)

	ci2, ok := c.get("a")

	assert.T(t, ok)
	assert.Equal(t, ci, ci2)
}

func TestCacheDropSz(t *testing.T) {
	c := newLRUCache(1)
	ci1 := &cachedItem{"a"}
	ci2 := &cachedItem{"b"}
	c.set(ci1)
	c.set(ci2)

	_, ok1 := c.get("a")
	out2, ok2 := c.get("b")

	assert.T(t, !ok1)
	assert.T(t, ok2)
	assert.Equal(t, ci2, out2)
}

func TestCacheDropLRU(t *testing.T) {
	c := newLRUCache(3)
	ci1 := &cachedItem{"a"}
	ci2 := &cachedItem{"b"}
	ci3 := &cachedItem{"c"}
	ci4 := &cachedItem{"d"}
	c.set(ci1)
	c.set(ci2)
	c.set(ci3)
	// Get ci1, so that ci2 (b) is the LRU to be dropped
	c.get(ci1.Key())
	c.set(ci4)

	_, ok := c.get("b")
	assert.T(t, !ok)
	_, ok = c.get("a")
	assert.T(t, ok)
	_, ok = c.get("c")
	assert.T(t, ok)
	_, ok = c.get("d")
	assert.T(t, ok)
}

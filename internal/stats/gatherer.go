package stats

import (
	"sort"
	"sync"
)

// Gatherer is a counter of provided keys
//
// Its use is:
// - Notify a key hit using Gatberer.Hit(key)
// - Retrieve the different hits using Gatherer.OrderedValues()
// - Reset the hits using Gatherer.Reset()
type Gatherer struct {
	mutex    sync.Mutex // not using RWMutex to avoid len issues on OrderedValues
	registry map[string]int
}

// Count reprents the number of hits a key encountered.
type Count struct {
	Key string `json:"key"`
	Hit int    `json:"hit"`
}

// NewGatherer will spawn a Gatherer instance.
func NewGatherer() *Gatherer {
	return &Gatherer{registry: make(map[string]int)}
}

// Hit acknowledges a key hit.
func (g *Gatherer) Hit(key string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	newValue := g.registry[key] + 1
	g.registry[key] = newValue
}

// OrderedValues gathers the hit keys as a slice of Count.
func (g *Gatherer) OrderedValues() []Count {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	values := make([]Count, 0, len(g.registry))
	for key, value := range g.registry {
		values = append(values, Count{Key: key, Hit: value})
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Hit > values[j].Hit
	})
	return values
}

// Reset trashes all previous hits.
func (g *Gatherer) Reset() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.registry = make(map[string]int)
}

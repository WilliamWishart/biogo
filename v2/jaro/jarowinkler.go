// https://github.com/hbollon/go-edlib/blob/master/jaro.go
package jaro

import "sync"

// JaroSimilarity returns a similarity index (between 0 and 1) using the Jaro distance algorithm.
// Only transposition operations are allowed.
func JaroSimilarity(a, b string) float32 {
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}
	if a == b {
		if len(a) == 0 {
			return 0.0
		}
		return 1.0
	}

	runeA := []rune(a)
	runeB := []rune(b)
	nA := len(runeA)
	nB := len(runeB)

	maxDist := Max(nA, nB)/2 - 1
	str1Table := make([]int8, nA)
	str2Table := make([]int8, nB)
	match := 0

	for i := 0; i < nA; i++ {
		start := Max(0, i-maxDist)
		end := Min(nB, i+maxDist+1)
		for j := start; j < end; j++ {
			if runeA[i] == runeB[j] && str2Table[j] == 0 {
				str1Table[i] = 1
				str2Table[j] = 1
				match++
				break
			}
		}
	}
	if match == 0 {
		return 0.0
	}
	t := 0
	p := 0
	for i := 0; i < nA; i++ {
		if str1Table[i] == 1 {
			for str2Table[p] == 0 {
				p++
			}
			if runeA[i] != runeB[p] {
				t++
			}
			p++
		}
	}
	t /= 2

	m := float32(match)
	return (m/float32(nA) + m/float32(nB) + (m-float32(t))/m) / 3.0
}

// JaroWinklerSimilarity returns a similarity index (between 0 and 1) using the Jaro-Winkler distance algorithm.
// It uses Jaro similarity and then looks for a common prefix (length <= 4).
func JaroWinklerSimilarity(a, b string) float32 {
	jaroSim := JaroSimilarity(a, b)

	if jaroSim != 0.0 && jaroSim != 1.0 {
		runeA := []rune(a)
		runeB := []rune(b)
		nA := len(runeA)
		nB := len(runeB)
		prefix := 0
		for i := 0; i < Min(nA, nB); i++ {
			if runeA[i] == runeB[i] {
				prefix++
			} else {
				break
			}
		}
		prefix = Min(prefix, 4)
		return jaroSim + 0.1*float32(prefix)*(1-jaroSim)
	}
	return jaroSim
}

// Jaro is a struct for caching Jaro similarity results.
type Jaro struct {
	cache map[[2]string]float32
	mu    sync.RWMutex
}

// NewJaro creates a new Jaro struct with an internal cache.
func NewJaro() *Jaro {
	return &Jaro{cache: make(map[[2]string]float32)}
}

// Similarity returns the Jaro similarity between two strings, using a cache for efficiency.
func (j *Jaro) Similarity(a, b string) float32 {
	key := [2]string{a, b}
	j.mu.RLock()
	if val, ok := j.cache[key]; ok {
		j.mu.RUnlock()
		return val
	}
	j.mu.RUnlock()
	val := JaroSimilarity(a, b)
	j.mu.Lock()
	j.cache[key] = val
	j.mu.Unlock()
	return val
}

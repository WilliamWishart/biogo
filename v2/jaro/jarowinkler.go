//https://github.com/hbollon/go-edlib/blob/master/jaro.go
package jaro

// JaroSimilarity return a similarity index (between 0 and 1)
// It use Jaro distance algorithm and allow only transposition operation
func JaroSimilarity(str1, str2 string) float32 {
	if len(str1) == 0 || len(str2) == 0 {
		return 0.0
	}
	if str1 == str2 {
		if len(str1) == 0 {
			return 0.0
		}
		return 1.0
	}

	runeStr1 := []rune(str1)
	runeStr2 := []rune(str2)
	runeStr1len := len(runeStr1)
	runeStr2len := len(runeStr2)

	maxDist := Max(runeStr1len, runeStr2len)/2 - 1
	str1Table := make([]int8, runeStr1len)
	str2Table := make([]int8, runeStr2len)
	match := 0

	for i := 0; i < runeStr1len; i++ {
		start := Max(0, i-maxDist)
		end := Min(runeStr2len, i+maxDist+1)
		for j := start; j < end; j++ {
			if runeStr1[i] == runeStr2[j] && str2Table[j] == 0 {
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
	for i := 0; i < runeStr1len; i++ {
		if str1Table[i] == 1 {
			for str2Table[p] == 0 {
				p++
			}
			if runeStr1[i] != runeStr2[p] {
				t++
			}
			p++
		}
	}
	t /= 2

	m := float32(match)
	return (m/float32(runeStr1len) + m/float32(runeStr2len) + (m-float32(t))/m) / 3.0
}

// JaroWinklerSimilarity return a similarity index (between 0 and 1)
// Use Jaro similarity and after look for a common prefix (length <= 4)
func JaroWinklerSimilarity(str1, str2 string) float32 {
	// Get Jaro similarity index between str1 and str2
	jaroSim := JaroSimilarity(str1, str2)

	if jaroSim != 0.0 && jaroSim != 1.0 {
		// Convert string parameters to rune arrays to be compatible with non-ASCII
		runeStr1 := []rune(str1)
		runeStr2 := []rune(str2)

		// Get and store length of these strings
		runeStr1len := len(runeStr1)
		runeStr2len := len(runeStr2)

		var prefix int

		// Find length of the common prefix
		for i := 0; i < Min(runeStr1len, runeStr2len); i++ {
			if runeStr1[i] == runeStr2[i] {
				prefix++
			} else {
				break
			}
		}

		// Normalized prefix count with Winkler's constraint
		// (prefix length must be inferior or equal to 4)
		prefix = Min(prefix, 4)

		// Return calculated Jaro-Winkler similarity index
		return jaroSim + 0.1*float32(prefix)*(1-jaroSim)
	}

	return jaroSim
}

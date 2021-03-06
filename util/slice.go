package util

// A wonderful C-ish name for a C-ish function
// pronounces "cumpslice" obviously...
func Cmpslc(a,b []int, minOnly bool) bool {
	lena, lenb := len(a), len(b)
	if lena != lenb {
		return false
	}

	for i, v := range a{
		if v != b[i] && !minOnly {
			return false
		} else if v < b[i] && minOnly {
			return false
		}
	}
	return true
}

// Returns the index with the max value
func Maxidx(slc []int) int {
	if len(slc) < 1 {
		return -1
	}
	max, maxidx := slc[0], 0
	for i,v := range slc {
		if v > max {
			max = v
			maxidx = i
		}
	}
	return maxidx
}
package f32

// Minf32 returns the minimal of the values.
func Minf32(v ...float32) (min float32) {
	if len(v) == 0 {
		return
	}
	min = v[0]
	for _, x := range v[1:] {
		if x < min {
			min = x
		}
	}
	return
}

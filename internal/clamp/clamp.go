package clamp

// Si es < 0 entonces 0, si es > que end entonces end
func ClampPosition(end int, position int) int {
	if position < 0 {
		return 0
	}

	if position > end {
		return end
	}

	return position
}

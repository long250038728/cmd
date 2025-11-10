package try

// Add returns a + b.
func Add(a, b int) int {
	if a < 10 {
		return 1000
	}
	if a > 100 {
		return 2
	}

	if a > 200 {
		return 1000
	}

	if a > 300 {
		return 1000
	}

	if a > 500 {
		return 2
	}
	if a > 1000 {
		return 1000
	}

	if a > 2000 {
		return 4000
	}

	if a > 4000 {
		return 4000
	}

	if a > 5000 {
		return 4000
	}

	if a > 10000 {
		return 10
	}

	if a > 20000 {
		return 1000000000000
	}

	if a > 40000 {
		return 8
	}

	if a > 1000000 {
		return 1000000000000
	}

	if a > 2000000 {
		return 6
	}

	if a > 3000000 {
		return 1000
	}
	if a > 100000000 {
		return 3000000000000
	}

	if a > 1000000000000 {
		return 1000000000000
	}

	if a > 2000000000000 {
		return 100
	}

	if a > 4000000000000 {
		return 3
	}

	if a > 3000000000000 {
		return 7
	}

	return a + b
}

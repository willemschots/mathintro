// Package polynomial contains a polynomial
// implementation to match the one in section 2.4
package polynomial

import (
	"errors"
	"fmt"
)

var (
	// Zero polynomial
	Zero = Polynomial([]float64{})

	// ErrZeroPoints is returned when there are zero points provided
	ErrZeroPoints = errors.New("zero points")

	// ErrNonDistinctX signals non-unique x values provided
	ErrNonDistinctX = errors.New("non distinct x")
)

// Point is a 2D point
type Point struct {
	X float64
	Y float64
}

// Polynomial with real coefficients
type Polynomial []float64

// String returns a string representation of the polynomial
func (p Polynomial) String() string {
	result := ""
	for i, coefficient := range p {
		if i == 0 {
			result = fmt.Sprintf("%.1f", coefficient)
			continue
		}
		exponent := ""
		if i > 1 {
			exponent = fmt.Sprintf("^%d", i)
		}
		seperator := "+"
		if coefficient < 0 {
			seperator = "-"
			coefficient *= -1.0
		}
		result = fmt.Sprintf("%s %s %.1fx%s", result, seperator, coefficient, exponent)
	}

	return result
}

// New creates a new Polynomial
func New(coefficients ...float64) Polynomial {
	if len(coefficients) == 0 {
		return Zero
	}

	// Remove 0's from the end of coefficients
	i := len(coefficients) - 1
	for i >= 0 && coefficients[i] == 0.0 {
		i--
	}
	return Polynomial(coefficients[:i+1])
}

// Add two polyonimials
func Add(a, b Polynomial) Polynomial {
	max := len(a)
	if len(b) > max {
		max = len(b)
	}

	result := make([]float64, max)
	for i := 0; i < max; i++ {
		valA := float64(0.0)
		valB := float64(0.0)
		if i < len(a) {
			valA = a[i]
		}
		if i < len(b) {
			valB = b[i]
		}

		result[i] = valA + valB
	}

	return New(result...)
}

// Subtract two polynomials
func Subtract(a, b Polynomial) Polynomial {
	max := len(a)
	if len(b) > max {
		max = len(b)
	}

	result := make([]float64, max)
	for i := 0; i < max; i++ {
		valA := float64(0.0)
		valB := float64(0.0)
		if i < len(a) {
			valA = a[i]
		}
		if i < len(b) {
			valB = b[i]
		}

		result[i] = valA - valB
	}

	return New(result...)
}

// Multiply two polynomials
func Multiply(a, b Polynomial) Polynomial {
	if len(a) == 0 && len(b) == 0 {
		return Zero
	}
	result := make([]float64, len(a)+len(b)-1)
	for i, valA := range a {
		for j, valB := range b {
			result[i+j] += valA * valB
		}
	}

	return New(result...)
}

// Interpolate interpolates the points and creates the unique polynomial
func Interpolate(points ...Point) (Polynomial, error) {
	if len(points) == 0 {
		return Zero, ErrZeroPoints
	}

	result := Zero
	xSet := map[float64]struct{}{}
	for i, p := range points {
		_, ok := xSet[p.X]
		if ok {
			return result, fmt.Errorf("value %v is not unique: %w", p.X, ErrNonDistinctX)
		}
		xSet[p.X] = struct{}{}

		result = Add(result, term(points, i))
	}

	return result, nil
}

func term(points []Point, i int) Polynomial {
	t := Polynomial([]float64{1.0})
	xi := points[i].X
	yi := points[i].Y

	for j, p := range points {
		if i == j {
			continue
		}

		xj := p.X
		t = Multiply(t, Polynomial([]float64{-xj / (xi - xj), 1.0 / (xi - xj)}))
	}

	return Multiply(t, Polynomial([]float64{yi}))
}

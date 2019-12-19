package polynomial_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/willemschots/mathintro/polynomial"
)

func TestNew(t *testing.T) {
	tests := []struct {
		in     []float64
		result polynomial.Polynomial
	}{
		{[]float64{}, polynomial.Zero},
		{[]float64{1.0, 2.0, 3.0}, polynomial.Polynomial([]float64{1.0, 2.0, 3.0})},
		{[]float64{1.0, 2.0, 0.0, 3.0, 0.0, 0.0, 0.0}, polynomial.Polynomial([]float64{1.0, 2.0, 0.0, 3.0})},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("in %v, out: %v", tc.in, tc.result), func(t *testing.T) {
			result := polynomial.New(tc.in...)
			assertPolynomial(t, tc.result, result)
		})
	}
}
func TestAdd(t *testing.T) {
	tests := []struct {
		a polynomial.Polynomial
		b polynomial.Polynomial
		c polynomial.Polynomial
	}{
		{
			polynomial.Zero,
			polynomial.Zero,
			polynomial.Zero,
		},
		{
			polynomial.Zero,
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
		},
		{
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0, 4.0}),
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
			polynomial.Polynomial([]float64{2.0, 4.0, 6.0, 4.0}),
		},
		{
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0, 4.0}),
			polynomial.Polynomial([]float64{2.0, 4.0, 6.0, 4.0}),
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v + %v = %v", tc.a, tc.b, tc.c), func(t *testing.T) {
			c := polynomial.Add(tc.a, tc.b)
			assertPolynomial(t, tc.c, c)
		})
	}
}
func TestSubtract(t *testing.T) {
	tests := []struct {
		a polynomial.Polynomial
		b polynomial.Polynomial
		c polynomial.Polynomial
	}{
		{
			polynomial.Zero,
			polynomial.Zero,
			polynomial.Zero,
		},
		{
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
			polynomial.Zero,
		},
		{
			polynomial.Polynomial([]float64{2.0, 4.0, 6.0}),
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
		},
		{
			polynomial.Zero,
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
			polynomial.Polynomial([]float64{-1.0, -2.0, -3.0}),
		},
		{
			polynomial.Polynomial([]float64{2.0, 4.0, 6.0}),
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0, 4.0}),
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0, -4.0}),
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v - %v = %v", tc.a, tc.b, tc.c), func(t *testing.T) {
			c := polynomial.Subtract(tc.a, tc.b)
			assertPolynomial(t, tc.c, c)
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		a polynomial.Polynomial
		b polynomial.Polynomial
		c polynomial.Polynomial
	}{
		{
			polynomial.Zero,
			polynomial.Zero,
			polynomial.Zero,
		},
		{
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0}),
			polynomial.Polynomial([]float64{4.0, 5.0, 6.0}),
			polynomial.Polynomial([]float64{4.0, 13.0, 28.0, 27.0, 18.0}),
		},
		{
			polynomial.Polynomial([]float64{1.0, 2.0, 3.0, -4.0}),
			polynomial.Polynomial([]float64{4.0, 5.0, 6.0}),
			polynomial.Polynomial([]float64{4.0, 13.0, 28.0, 11.0, -2.0, -24.0}),
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v * %v = %v", tc.a, tc.b, tc.c), func(t *testing.T) {
			c := polynomial.Multiply(tc.a, tc.b)
			assertPolynomial(t, tc.c, c)
		})
	}
}

func TestInterPolate(t *testing.T) {
	tests := []struct {
		input  []polynomial.Point
		result polynomial.Polynomial
		err    error
	}{
		{
			[]polynomial.Point{},
			polynomial.Zero,
			polynomial.ErrZeroPoints,
		},
		{
			[]polynomial.Point{
				polynomial.Point{1.0, 2.0},
				polynomial.Point{1.0, 3.0},
			},
			polynomial.Zero,
			polynomial.ErrNonDistinctX,
		},
		{
			[]polynomial.Point{
				polynomial.Point{1.0, 2.0},
			},
			polynomial.New(2.0),
			nil,
		},
		{
			[]polynomial.Point{
				polynomial.Point{1.0, 2.0},
				polynomial.Point{2.0, 3.0},
			},
			polynomial.New(1.0, 1.0),
			nil,
		},
		{
			[]polynomial.Point{
				polynomial.Point{1.0, 2.0},
				polynomial.Point{2.0, 5.0},
				polynomial.Point{3.0, 2.0},
			},
			polynomial.New(-7.0, 12.0, -3.0),
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Interpolate(%v) = %s", tc.input, tc.result.String()), func(t *testing.T) {

			result, err := polynomial.Interpolate(tc.input...)
			// Negative cases
			if tc.err != nil {
				unwrap := errors.Unwrap(err)
				if unwrap != nil {
					err = unwrap
				}
				if err != tc.err {
					t.Errorf("expected error %v but got %v", tc.err, err)
				}
				return
			}

			// Positive cases
			if err != nil {
				t.Errorf("expected nil error but got %v", err)
				return
			}

			assertPolynomial(t, tc.result, result)
		})
	}
}

func assertPolynomial(t *testing.T, expect, actual polynomial.Polynomial) {
	if len(expect) != len(actual) {
		t.Errorf("expected number of coefficients to be %v but got %v", len(expect), len(actual))
		return
	}

	for i, coefficient := range expect {
		if coefficient != actual[i] {
			t.Errorf("expected coefficient %d to be %v but got %v", i, coefficient, actual[i])
		}
	}
}

package pdp

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Test Integer Equal a == b
func TestIntegerEqual(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b int64
	}{
		{
			a: 1, b: 1,
		},
		{
			a: 0, b: 0,
		},
		{
			a: -1, b: -1,
		},
		{
			a: 1, b: 0,
		},
		{
			a: 0, b: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Integer Equal %d == %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeIntegerValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionIntegerEqual(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a == tc.b)
			res, err := v.boolean()
			if err != nil {
				t.Errorf("Expect boolean result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%t', but got '%t'", expect, res)
			}
		})
	}
}

// Test Integer Greater a > b
func TestIntegerGreater(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b int64
	}{
		{
			a: 1, b: 0,
		},
		{
			a: 0, b: -1,
		},
		{
			a: 1, b: -1,
		},
		{
			a: 0, b: 0,
		},
		{
			a: -1, b: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Integer Greater %d > %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeIntegerValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionIntegerGreater(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a > tc.b)
			res, err := v.boolean()
			if err != nil {
				t.Errorf("Expect boolean result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%t', but got '%t'", expect, res)
			}
		})
	}
}

// Test Integer Add a + b
func TestIntegerAdd(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b int64
	}{
		{
			a: 1, b: 1,
		},
		{
			a: 0, b: 0,
		},
		{
			a: -1, b: -1,
		},
		{
			a: 1, b: 0,
		},
		{
			a: -1, b: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Integer Add %d + %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeIntegerValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionIntegerAdd(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a + tc.b)
			res, err := v.integer()
			if err != nil {
				t.Errorf("Expect integer result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%d', but got '%d'", expect, res)
			}
		})
	}
}

// Test Integer Substract a - b
func TestIntegerSubtract(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b int64
	}{
		{
			a: 1, b: 1,
		},
		{
			a: 0, b: 0,
		},
		{
			a: -1, b: -1,
		},
		{
			a: 1, b: -1,
		},
		{
			a: 0, b: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Integer Subtract %d - %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeIntegerValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionIntegerSubtract(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a - tc.b)
			res, err := v.integer()
			if err != nil {
				t.Errorf("Expect integer result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%d', but got '%d'", expect, res)
			}
		})
	}
}

// Test Integer Multiply a * b
func TestIntegerMultiply(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b int64
	}{
		{
			a: 9, b: 1,
		},
		{
			a: 11, b: 0,
		},
		{
			a: -3, b: -4,
		},
		{
			a: 2, b: -2,
		},
		{
			a: -23, b: 97,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Integer Multiply %d * %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeIntegerValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionIntegerMultiply(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a * tc.b)
			res, err := v.integer()
			if err != nil {
				t.Errorf("Expect integer result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%d', but got '%d'", expect, res)
			}
		})
	}
}

// Test Integer Divide a / b
func TestIntegerDivide(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b int64
		err  string
	}{
		{
			a: 9, b: 1,
		},
		{
			a: 0, b: 1,
		},
		{
			a: 12, b: 2,
		},
		{
			a: 12, b: 5,
		},
		{
			a: -20, b: 6,
		},
		{
			a: 85, b: -13,
		},
		{
			a: 27, b: 0, err: "Integer divisor has a value of 0",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Integer Divide %d / %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeIntegerValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionIntegerDivide(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				if tc.err == "" {
					t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				} else if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("Expect Calculate() returns error contains '%s', but got '%s'", tc.err, err)
				}
				return
			} else if tc.err != "" {
				t.Errorf("Expect Calculate() returns error contains '%s', but got nil", tc.err)
				return
			}

			var expect int64
			if tc.b != 0 {
				expect = (tc.a / tc.b)
			}
			res, err := v.integer()
			if err != nil {
				t.Errorf("Expect integer result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%d', but got '%d'", expect, res)
			}
		})
	}
}

// Test Integer Range: range min max val
func TestIntegerRange(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		min, max, val int64
		expect        string
	}{
		{
			min: 1, max: 5, val: 0,
			expect: "Below",
		},
		{
			min: 1, max: 5, val: 7,
			expect: "Above",
		},
		{
			min: 1, max: 5, val: 3,
			expect: "Within",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Integer Range: range %d %d %d", tc.min, tc.max, tc.val), func(t *testing.T) {
			min := MakeIntegerValue(tc.min)
			max := MakeIntegerValue(tc.max)
			val := MakeIntegerValue(tc.val)
			e := makeFunctionIntegerRange(min, max, val)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			res, err := v.str()
			if err != nil {
				t.Errorf("Expect string result with no error, but got '%s'", err)
			} else if res != tc.expect {
				t.Errorf("Expect result '%s', but got '%s'", tc.expect, res)
			}
		})
	}
}

// Test Float Greater a > b
func TestFloatGreater(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b float64
	}{
		{
			a: 1.0, b: 0.9,
		},
		{
			a: 12.3, b: -12.3,
		},
		{
			a: -38.23, b: 38.23,
		},
		{
			a: 0.0, b: -0.4927,
		},
		{
			a: 0.8, b: 0.9,
		},
		{
			a: 0.0, b: 9735.23,
		},
		{
			a: 7.26491E+11, b: 7.26490E+11,
		},
		{
			a: -8.36591E-33, b: -5.36591E-33,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float Greater %G > %G", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeFloatValue(tc.b)
			e := makeFunctionFloatGreater(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a > tc.b)
			res, err := v.boolean()
			if err != nil {
				t.Errorf("Expect boolean result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%t', but got '%t'", expect, res)
			}
		})
	}
}

// Test Float Add a + b
func TestFloatAdd(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b float64
	}{
		{
			a: 72.2, b: 29.3,
		},
		{
			a: 0.0, b: 0.0,
		},
		{
			a: -38.23, b: -50.34,
		},
		{
			a: 0.0, b: -0.4927,
		},
		{
			a: 2.8769E+10, b: 5.7362E+11,
		},
		{
			a: 7.26491E+11, b: -7.26490E+11,
		},
		{
			a: -8.36591E-33, b: -5.36591E-33,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float Add %G + %G", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeFloatValue(tc.b)
			e := makeFunctionFloatAdd(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a + tc.b)
			res, err := v.float()
			if err != nil {
				t.Errorf("Expect float result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%G', but got '%G'", expect, res)
			}
		})
	}
}

// Test Float Subtract a - b
func TestFloatSubtract(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b float64
	}{
		{
			a: 72.0, b: 29.0,
		},
		{
			a: 0.0, b: 0.0,
		},
		{
			a: -38.23, b: -50.34,
		},
		{
			a: 0.0, b: -0.4927,
		},
		{
			a: 2.8769E+10, b: 5.7362E+11,
		},
		{
			a: 7.26491E+11, b: -7.26490E+11,
		},
		{
			a: -8.36591E-33, b: -5.36591E-33,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float Subtract %G - %G", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeFloatValue(tc.b)
			e := makeFunctionFloatSubtract(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a - tc.b)
			res, err := v.float()
			if err != nil {
				t.Errorf("Expect float result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%G', but got '%G'", expect, res)
			}
		})
	}
}

// Test Float Multiply a * b
func TestFloatMultiply(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b float64
		err  string
	}{
		{
			a: 72.0, b: 29.0,
		},
		{
			a: 0.0, b: 0.0,
		},
		{
			a: -38.23, b: -50.34,
		},
		{
			a: 0.0, b: -0.4927,
		},
		{
			a: 2.8769E+10, b: 5.7362E+11,
		},
		{
			a: 7.26491E+11, b: -7.26490E+11,
		},
		{
			a: -8.36591E-33, b: -5.36591E-33,
		},
		{
			a: 1.9E+200, b: 5.3E+233,
			err: "Float result has a value of Inf",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float Multiply %g * %g", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeFloatValue(tc.b)
			e := makeFunctionFloatMultiply(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				if tc.err == "" {
					t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				} else if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("Expect Calculate() returns error contains '%s', but got '%s'", tc.err, err)
				}
				return
			} else if tc.err != "" {
				t.Errorf("Expect Calculate() returns error contains '%s', but got nil", tc.err)
				return
			}

			expect := (tc.a * tc.b)
			res, err := v.float()
			if err != nil {
				t.Errorf("Expect float result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%g', but got '%g'", expect, res)
			}
		})
	}
}

// Test Float Divide a / b
func TestFloatDivide(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a, b float64
		err  string
	}{
		{
			a: 72.0, b: 29.0,
		},
		{
			a: 0.0, b: 0.0,
			err: "Float divisor has a value of 0",
		},
		{
			a: -38.23, b: -50.34,
		},
		{
			a: 0.0, b: -0.4927,
		},
		{
			a: 2.8769E+10, b: 5.7362E+11,
		},
		{
			a: 7.26491E+11, b: -7.26490E+11,
		},
		{
			a: -8.36591E-33, b: -5.36591E-33,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float Divide %g / %g", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeFloatValue(tc.b)
			e := makeFunctionFloatDivide(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				if tc.err == "" {
					t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				} else if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("Expect Calculate() returns error contains '%s', but got '%s'", tc.err, err)
				}
				return
			} else if tc.err != "" {
				t.Errorf("Expect Calculate() returns error contains '%s', but got nil", tc.err)
				return
			}

			expect := (tc.a / tc.b)
			res, err := v.float()
			if err != nil {
				t.Errorf("Expect float result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%g', but got '%g'", expect, res)
			}
		})
	}
}

// Test Float-Integer Equal a == b
func TestFloatIntegerEqual(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a float64
		b int64
	}{
		{
			a: 1., b: 1,
		},
		{
			a: 0., b: 0,
		},
		{
			a: -1., b: -1,
		},
		{
			a: 12345.0, b: 12345,
		},
		{
			a: -6789.0, b: -6789,
		},
		{
			a: 3.1, b: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float-Integer Equal %G == %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionFloatEqual(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a == float64(tc.b))
			res, err := v.boolean()
			if err != nil {
				t.Errorf("Expect boolean result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%t', but got '%t'", expect, res)
			}
		})
	}
}

// Test Float-Integer Greater a > b
func TestFloatIntegerGreater(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a float64
		b int64
	}{
		{
			a: 1., b: 0,
		},
		{
			a: 0., b: 0,
		},
		{
			a: -134.5, b: -134,
		},
		{
			a: 12345.1, b: 12345,
		},
		{
			a: -6789.0, b: -6789,
		},
		{
			a: 3.4, b: 4,
		},
		{
			a: -4., b: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float-Integer Greater %G > %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionFloatGreater(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a > float64(tc.b))
			res, err := v.boolean()
			if err != nil {
				t.Errorf("Expect boolean result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%t', but got '%t'", expect, res)
			}
		})
	}
}

// Test Float-Integer Add a + b
func TestFloatIntegerAdd(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a float64
		b int64
	}{
		{
			a: 72.2, b: 29,
		},
		{
			a: 0.0, b: 0,
		},
		{
			a: -38.23, b: -50,
		},
		{
			a: 8694.3, b: -4927,
		},
		{
			a: 2.8769E+10, b: 2372766369,
		},
		{
			a: -7.26491E+11, b: -939284,
		},
		{
			a: -8.36591E-33, b: 536591,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float-Integer Add %G + %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionFloatAdd(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a + float64(tc.b))
			res, err := v.float()
			if err != nil {
				t.Errorf("Expect float result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%G', but got '%G'", expect, res)
			}
		})
	}
}

// Test Float-Integer Subtract a - b
func TestFloatIntegerSubtract(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a float64
		b int64
	}{
		{
			a: 72.2, b: 29,
		},
		{
			a: 0.0, b: 0,
		},
		{
			a: -38.23, b: -50,
		},
		{
			a: 8694.3, b: -4927,
		},
		{
			a: 2.8769E+10, b: 2372766369,
		},
		{
			a: -7.26491E+11, b: -939284,
		},
		{
			a: -8.36591E-33, b: 536591,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float-Integer Subtract %G - %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionFloatSubtract(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			expect := (tc.a - float64(tc.b))
			res, err := v.float()
			if err != nil {
				t.Errorf("Expect float result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%G', but got '%G'", expect, res)
			}
		})
	}
}

// Test Float-Integer Multiply a * b
func TestFloatIntegerMultiply(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a   float64
		b   int64
		err string
	}{
		{
			a: 72.2, b: 29,
		},
		{
			a: 0.0, b: 0,
		},
		{
			a: -38.23, b: -50,
		},
		{
			a: 8694.3, b: -4927,
		},
		{
			a: 2.8769E+10, b: 2372766369,
		},
		{
			a: -7.26491E+11, b: -939284,
		},
		{
			a: -8.36591E-33, b: 536591,
		},
		{
			a: 1E+308, b: 3, err: "Float result has a value of Inf",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float-Integer Multiply %G * %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionFloatMultiply(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				if tc.err == "" {
					t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				} else if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("Expect Calculate() returns error contains '%s', but got '%s'", tc.err, err)
				}
				return
			} else if tc.err != "" {
				t.Errorf("Expect Calculate() returns error contains '%s', but got nil", tc.err)
				return
			}

			expect := (tc.a * float64(tc.b))
			res, err := v.float()
			if err != nil {
				t.Errorf("Expect float result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%G', but got '%G'", expect, res)
			}
		})
	}
}

// Test Float-Integer Divide a / b
func TestFloatIntegerDivide(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		a   float64
		b   int64
		err string
	}{
		{
			a: 72.0, b: 29,
		},
		{
			a: 187837.83, b: 0,
			err: "Float divisor has a value of 0",
		},
		{
			a: -38.23, b: -50,
		},
		{
			a: 0.0, b: -4927,
		},
		{
			a: 2.8769E+10, b: 7362,
		},
		{
			a: 7.26491E+11, b: -26490,
		},
		{
			a: -8.36591E-33, b: -36591,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float-Integer Divide %g / %d", tc.a, tc.b), func(t *testing.T) {
			a := MakeFloatValue(tc.a)
			b := MakeIntegerValue(tc.b)
			e := makeFunctionFloatDivide(a, b)

			v, err := e.Calculate(ctx)
			if err != nil {
				if tc.err == "" {
					t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				} else if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("Expect Calculate() returns error contains '%s', but got '%s'", tc.err, err)
				}
				return
			} else if tc.err != "" {
				t.Errorf("Expect Calculate() returns error contains '%s', but got nil", tc.err)
				return
			}

			expect := (tc.a / float64(tc.b))
			res, err := v.float()
			if err != nil {
				t.Errorf("Expect float result with no error, but got '%s'", err)
			} else if res != expect {
				t.Errorf("Expect result '%g', but got '%g'", expect, res)
			}
		})
	}
}

// Test Float Range: range min max val
func TestFloatRange(t *testing.T) {
	ctx, err := NewContext(nil, 0, nil)
	if err != nil {
		t.Fatalf("Expected context but got error %s", err)
	}

	testCases := []struct {
		min, max, val float64
		expect        string
	}{
		{
			min: 1.0, max: 5.0, val: 0.1,
			expect: "Below",
		},
		{
			min: 1.0, max: 5.0, val: 11.1,
			expect: "Above",
		},
		{
			min: 1.0, max: 5.0, val: 3.3,
			expect: "Within",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Float Range: range %G %G %G", tc.min, tc.max, tc.val), func(t *testing.T) {
			min := MakeFloatValue(tc.min)
			max := MakeFloatValue(tc.max)
			val := MakeFloatValue(tc.val)
			e := makeFunctionFloatRange(min, max, val)

			v, err := e.Calculate(ctx)
			if err != nil {
				t.Errorf("Expect Calculate() returns no error, but got '%s'", err)
				return
			}

			res, err := v.str()
			if err != nil {
				t.Errorf("Expect string result with no error, but got '%s'", err)
			} else if res != tc.expect {
				t.Errorf("Expect result '%s', but got '%s'", tc.expect, res)
			}
		})
	}
}

func TestFunctionArgumentValidators(t *testing.T) {
	type testCase struct {
		argTypes         Signature
		expectFunc       reflect.Type
		expectResultType Type
	}

	type testSuite struct {
		funcName  string
		testCases []testCase
	}

	testSuites := []testSuite{
		{
			funcName: "add",
			testCases: []testCase{
				{
					argTypes:         MakeSignature(TypeInteger, TypeInteger),
					expectFunc:       reflect.TypeOf(functionIntegerAdd{}),
					expectResultType: TypeInteger,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatAdd{}),
					expectResultType: TypeFloat,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatAdd{}),
					expectResultType: TypeFloat,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatAdd{}),
					expectResultType: TypeFloat,
				},
			},
		},
		{
			funcName: "subtract",
			testCases: []testCase{
				{
					argTypes:         MakeSignature(TypeInteger, TypeInteger),
					expectFunc:       reflect.TypeOf(functionIntegerSubtract{}),
					expectResultType: TypeInteger,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatSubtract{}),
					expectResultType: TypeFloat,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatSubtract{}),
					expectResultType: TypeFloat,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatSubtract{}),
					expectResultType: TypeFloat,
				},
			},
		},
		{
			funcName: "multiply",
			testCases: []testCase{
				{
					argTypes:         MakeSignature(TypeInteger, TypeInteger),
					expectFunc:       reflect.TypeOf(functionIntegerMultiply{}),
					expectResultType: TypeInteger,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatMultiply{}),
					expectResultType: TypeFloat,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatMultiply{}),
					expectResultType: TypeFloat,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatMultiply{}),
					expectResultType: TypeFloat,
				},
			},
		},
		{
			funcName: "divide",
			testCases: []testCase{
				{
					argTypes:         MakeSignature(TypeInteger, TypeInteger),
					expectFunc:       reflect.TypeOf(functionIntegerDivide{}),
					expectResultType: TypeInteger,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatDivide{}),
					expectResultType: TypeFloat,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatDivide{}),
					expectResultType: TypeFloat,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatDivide{}),
					expectResultType: TypeFloat,
				},
			},
		},
		{
			funcName: "equal",
			testCases: []testCase{
				{
					argTypes:         MakeSignature(TypeInteger, TypeInteger),
					expectFunc:       reflect.TypeOf(functionIntegerEqual{}),
					expectResultType: TypeBoolean,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatEqual{}),
					expectResultType: TypeBoolean,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatEqual{}),
					expectResultType: TypeBoolean,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatEqual{}),
					expectResultType: TypeBoolean,
				},
			},
		},
		{
			funcName: "greater",
			testCases: []testCase{
				{
					argTypes:         MakeSignature(TypeInteger, TypeInteger),
					expectFunc:       reflect.TypeOf(functionIntegerGreater{}),
					expectResultType: TypeBoolean,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatGreater{}),
					expectResultType: TypeBoolean,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatGreater{}),
					expectResultType: TypeBoolean,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatGreater{}),
					expectResultType: TypeBoolean,
				},
			},
		},
		{
			funcName: "range",
			testCases: []testCase{
				{
					argTypes:         MakeSignature(TypeInteger, TypeInteger, TypeInteger),
					expectFunc:       reflect.TypeOf(functionIntegerRange{}),
					expectResultType: TypeString,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeInteger, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatRange{}),
					expectResultType: TypeString,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeFloat, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatRange{}),
					expectResultType: TypeString,
				},
				{
					argTypes:         MakeSignature(TypeInteger, TypeFloat, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatRange{}),
					expectResultType: TypeString,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeInteger, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatRange{}),
					expectResultType: TypeString,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeInteger, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatRange{}),
					expectResultType: TypeString,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeFloat, TypeInteger),
					expectFunc:       reflect.TypeOf(functionFloatRange{}),
					expectResultType: TypeString,
				},
				{
					argTypes:         MakeSignature(TypeFloat, TypeFloat, TypeFloat),
					expectFunc:       reflect.TypeOf(functionFloatRange{}),
					expectResultType: TypeString,
				},
			},
		},
	}

	for _, ts := range testSuites {
		t.Run(fmt.Sprintf("Validators for \"%s\"", ts.funcName), func(t *testing.T) {
			for _, tc := range ts.testCases {
				t.Run(fmt.Sprintf("Argument Types: %v", tc.argTypes), func(t *testing.T) {

					args := make([]Expression, len(tc.argTypes))
					for i, argType := range tc.argTypes {
						expr, err := MakeValueFromString(argType, "1")
						if err != nil {
							t.Fatalf("Unexpected error making expression of type %s", argType)
							return
						}
						args[i] = expr
					}

					for _, fav := range FunctionArgumentValidators[ts.funcName] {
						funcMaker := fav(args)
						if funcMaker != nil {
							numFunc := funcMaker(args)
							if tc.expectFunc != reflect.TypeOf(numFunc) {
								t.Errorf("Expected numeric function is '%v' but got '%v'",
									tc.expectFunc, reflect.TypeOf(numFunc))
							}
							if tc.expectResultType != numFunc.GetResultType() {
								t.Errorf("Expected result type is '%s' but got '%s'",
									tc.expectResultType, numFunc.GetResultType())
							}
							return
						}
					}
					t.Error("Expect to find function maker but none found")
				})
			}
		})
	}
}

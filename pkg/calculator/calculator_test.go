package calculator

import (
	"reflect"
	"testing"
)

func TestIsNumber(t *testing.T) {
	testCases := []struct {
		input string
		want  bool
	}{
		{"1", true},
		{"-1", true},
		{"1.0", true},
		{"-1.0", true},
		{"-1.00", true},
		{"1.23", true},
		{"-1.23", true},
		{"1.23e3", true},
		{"-1.23e3", true},
		{"a", false},
		{"1a", false},
		{"a1", false},
	}

	for _, tc := range testCases {
		got := isNumber(tc.input)
		if got != tc.want {
			t.Errorf("IsNumber(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestIsOperator(t *testing.T) {
	testCases := []struct {
		input rune
		want  bool
	}{
		{'+', true},
		{'-', true},
		{'*', true},
		{'/', true},
		{'(', false},
		{')', false},
		{'1', false},
		{'a', false},
	}

	for _, tc := range testCases {
		got := isOperator(tc.input)
		if got != tc.want {
			t.Errorf("IsOperator(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestTokenize(t *testing.T) {
	testCases := []struct {
		input string
		want  []string
	}{
		{"1+2*3", []string{"1", "+", "2", "*", "3"}},
		{"(1+2)*3", []string{"(", "1", "+", "2", ")", "*", "3"}},
		{"-1+2", []string{"-1", "*", "1", "+", "2"}},
		{"2*-3", []string{"2", "*", "-1", "*", "3"}},
		{"2*-3.14", []string{"2", "*", "-1", "*", "3.14"}},
	}

	for _, tc := range testCases {
		got, err := tokenize(tc.input)
		if err != nil {
			t.Errorf("Tokenize(%q) returned unexpected error: %v", tc.input, err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Tokenize(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestTokenizePositive(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected []string
	}{
		{
			name:     "Simple expression",
			expr:     "1 + 2",
			expected: []string{"1", "+", "2"},
		},
		{
			name:     "Expression with negative number",
			expr:     "-1 + 2",
			expected: []string{"-1", "*", "1", "+", "2"},
		},
		{
			name:     "Expression with brackets",
			expr:     "(-1 + 2) * 3",
			expected: []string{"(", "-1", "*", "1", "+", "2", ")", "*", "3"},
		},
		{
			name:     "Expression with decimal numbers",
			expr:     "1.2 - 0.5 * 3",
			expected: []string{"1.2", "-", "0.5", "*", "3"},
		},
		{
			name:     "Complex expression",
			expr:     "(-2.3 + 4.5) / 1.1 - 2 * 3.3",
			expected: []string{"(", "-1", "*", "2.3", "+", "4.5", ")", "/", "1.1", "-", "2", "*", "3.3"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := tokenize(test.expr)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, res)
			}
		})
	}
}

func TestTokenizeNegative(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{
			name: "Decimal dot not after number",
			expr: ". + 1",
		},
		{
			name: "More than one operators in a row",
			expr: "1.234 + 1 -* 2",
		},
		{
			name: "More than two minus in a row",
			expr: "1.234 + 1 --- 2",
		},
		{
			name: "Unterminated bracket at start",
			expr: "(2 + 3",
		},
		{
			name: "Unterminated brackets",
			expr: "(2 + ((((((((3)))",
		},
		{
			name: "Unterminated bracket at end",
			expr: "2 + 3)",
		},
		{
			name: "Mismatch operator",
			expr: "* 1",
		},
		{
			name: "Operator after unary minus",
			expr: "-* 1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := tokenize(test.expr)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestPrecedence(t *testing.T) {
	testCases := []struct {
		input string
		want  int
	}{
		{"+", 1},
		{"-", 1},
		{"*", 2},
		{"/", 2},
		{"", 0},
		{"a", 0},
	}

	for _, tc := range testCases {
		got := precedence(tc.input)
		if got != tc.want {
			t.Errorf("Precedence(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestToRPN(t *testing.T) {
	testCases := []struct {
		input []string
		want  []string
	}{
		{[]string{"1", "+", "2"}, []string{"1", "2", "+"}},
		{[]string{"1", "+", "2", "*", "3"}, []string{"1", "2", "3", "*", "+"}},
		{[]string{"(", "1", "+", "2", ")", "*", "3"}, []string{"1", "2", "+", "3", "*"}},
	}

	for _, tc := range testCases {
		got, err := toRPN(tc.input)
		if err != nil {
			t.Errorf("ToRPN(%v) returned unexpected error: %v", tc.input, err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("ToRPN(%v) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestCalculateRPN(t *testing.T) {
	testCases := []struct {
		input []string
		want  float64
	}{
		{[]string{"1", "2", "+"}, 3},
		{[]string{"1", "2", "3", "*", "+"}, 7},
		{[]string{"1", "2", "+", "3", "*"}, 9},
		{[]string{"3", "2", "1", "*", "+"}, 5},
	}

	for _, tc := range testCases {
		got, err := calculateRPN(tc.input)
		if err != nil {
			t.Errorf("calculateRPN(%v) returned unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("calculateRPN(%v) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestEvaluate(t *testing.T) {
	testCases := []struct {
		input string
		want  float64
	}{
		{"1 + 2", 3},
		{"1.5 + 2", 3.5},
		{"1.5 - 2", -0.5},
		{"2 * 3", 6},
		{"6 / 3", 2},
		{"2 + 3 * 4", 14},
		{"(2 + 3) * 4", 20},
		{"2 + 3 * (4 / 2)", 8},
		{"(2 + 3) * (4 + 3)", 35},
		{"-1 + 2", 1},
		{"2 * -3", -6},
		{"2.2 * -3.4", -7.48},
	}

	for _, tc := range testCases {
		got, err := Evaluate(tc.input)
		if err != nil {
			t.Errorf("Evaluate(%q) returned unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("Evaluate(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestEvaluateErrors(t *testing.T) {
	testCases := []struct {
		input   string
		wantErr bool
	}{
		{"1.0.0 + 2", true},
		{"1 + a", true},
		{"1 +", true},
		{"1 / 0", true},
		{"1 + (2 * (3", true},
		{"1 + 2) * 3", true},
		{"1.1.1 + 2 * 3", true},
	}

	for _, tc := range testCases {
		_, err := Evaluate(tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("Evaluate(%q) error = %v, want Error %v", tc.input, err, tc.wantErr)
		}
	}
}

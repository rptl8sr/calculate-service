package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Evaluate takes a mathematical expression as a string and returns the result or an error.
func Evaluate(expr string) (float64, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}

	rpn, err := toRPN(tokens)
	if err != nil {
		return 0, err
	}

	return calculateRPN(rpn)
}

type tokenType string

const (
	Empty        tokenType = ""
	Number       tokenType = "number"
	Operator     tokenType = "operator"
	Dot          tokenType = "."
	UnaryMinus   tokenType = "unaryMinus"
	Add          tokenType = "+"
	Sub          tokenType = "-"
	Multi        tokenType = "*"
	Div          tokenType = "/"
	BracketLeft  tokenType = "("
	BracketRight tokenType = ")"
)

// tokenize converts the input string into a slice of tokens.
func tokenize(input string) ([]string, error) {
	var tokens []string
	var brackets int
	var currToken strings.Builder

	prevTokenType := Empty

	for i, r := range input {
		switch {
		case unicode.IsSpace(r):
			continue
		case unicode.IsDigit(r):
			if prevTokenType == UnaryMinus {
				tokens = append(tokens, "-1", "*")
			}
			currToken.WriteRune(r)
			prevTokenType = Number
			continue
		case tokenType(r) == Dot:
			switch prevTokenType {
			case Number:
				currToken.WriteRune(r)
				prevTokenType = Dot
				continue
			default:
				return nil, NewCalcError(ErrInsufficientValues, fmt.Sprintf("decimal dot delimiter not after number, position %d: %c", i, r))
			}
		case tokenType(r) == BracketLeft:
			switch prevTokenType {
			case Empty, BracketLeft, Operator:
				if currToken.Len() > 0 {
					tokens = append(tokens, currToken.String())
					currToken.Reset()
				}
			case UnaryMinus:
				tokens = append(tokens, "-1", "*")
			default:
				return nil, NewCalcError(ErrMismatchedParentheses, fmt.Sprintf("position %d: %c", i, r))
			}
			tokens = append(tokens, string(r))
			brackets++
			prevTokenType = BracketLeft
		case tokenType(r) == BracketRight:
			switch prevTokenType {
			case Number, BracketRight:
				if currToken.Len() > 0 {
					tokens = append(tokens, currToken.String())
					currToken.Reset()
				}
			default:
				return nil, NewCalcError(ErrMismatchedParentheses, fmt.Sprintf("position %d: %c", i, r))
			}
			tokens = append(tokens, string(r))
			brackets--
			if brackets < 0 {
				return nil, NewCalcError(ErrMismatchedParentheses, fmt.Sprintf("position %d: %c", i, r))
			}
			prevTokenType = BracketRight
		case tokenType(r) == Sub:
			switch prevTokenType {
			case Empty, BracketLeft, Operator:
				prevTokenType = UnaryMinus
				continue
			case Number, BracketRight:
				if currToken.Len() > 0 {
					tokens = append(tokens, currToken.String())
				}
				tokens = append(tokens, string(r))
				currToken.Reset()
				prevTokenType = Operator
			default:
				return nil, NewCalcError(ErrTooManyValues, fmt.Sprintf("position %d: %c", i, r))
			}
		case strings.ContainsRune("+*/", r):
			switch prevTokenType {
			case Number, BracketRight:
				if currToken.Len() > 0 {
					tokens = append(tokens, currToken.String())
				}
				tokens = append(tokens, string(r))
				currToken.Reset()
				prevTokenType = Operator
			default:
				return nil, NewCalcError(ErrMismatchOperator, fmt.Sprintf("position %d: %c", i, r))
			}
		default:
			return nil, NewCalcError(ErrInvalidCharacter, fmt.Sprintf("position %d: %c", i, r))
		}
	}

	if currToken.Len() > 0 {
		tokens = append(tokens, currToken.String())
	}

	if brackets != 0 {
		return nil, NewCalcError(ErrMismatchedParentheses, "unterminated last parentheses' group")
	}

	return tokens, nil
}

// isOperator checks if a rune is an arithmetic operator.
func isOperator(ch rune) bool {
	return tokenType(ch) == Add || tokenType(ch) == Sub || tokenType(ch) == Multi || tokenType(ch) == Div
}

// precedence returns the precedence of an operator.
func precedence(op string) int {
	switch tokenType(op) {
	case Add, Sub:
		return 1
	case Multi, Div:
		return 2
	}
	return 0
}

// toRPN converts a list of tokens to Reverse Polish Notation using the Shunting Yard algorithm.
func toRPN(tokens []string) ([]string, error) {
	var output []string
	var operators []string

	for i, token := range tokens {
		switch {
		case isNumber(token):
			output = append(output, token)
		case tokenType(token) == BracketLeft:
			operators = append(operators, token)
		case tokenType(token) == BracketRight:
			for len(operators) > 0 && tokenType(operators[len(operators)-1]) != BracketLeft {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, NewCalcError(ErrMismatchedParentheses, fmt.Sprintf("position %d: %s", i, token))
			}
			operators = operators[:len(operators)-1]
		case isOperator(rune(token[0])):
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		default:
			return nil, NewCalcError(ErrInsufficientValues, fmt.Sprintf("position %d: %s", i, token))
		}
	}

	for len(operators) > 0 {
		if tokenType(operators[len(operators)-1]) == BracketLeft {
			return nil, NewCalcError(ErrMismatchedParentheses, "")
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

// isNumber checks if a string represents a number.
func isNumber(s string) bool {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	// Check for large numbers
	if num > 1e308 || num < -1e308 {
		return false
	}
	return true
}

// calculateRPN calculates the result of an expression in Reverse Polish Notation.
func calculateRPN(rpn []string) (float64, error) {
	var stack []float64

	for i, token := range rpn {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else if isOperator(rune(token[0])) {
			if len(stack) < 2 {
				return 0, NewCalcError(ErrInsufficientValues, fmt.Sprintf("position %d: %s", i, token))
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch tokenType(token) {
			case Add:
				result = a + b
			case Sub:
				result = a - b
			case Multi:
				result = a * b
			case Div:
				if b == 0 {
					return 0, NewCalcError(ErrDivisionByZero, "")
				}
				result = a / b
			}
			// Check for large numbers after operation
			if result > 1e308 || result < -1e308 {
				return 0, NewCalcError(ErrTooLargeNumber, fmt.Sprintf("%f", result))
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, NewCalcError(ErrTooManyValues, "")
	}

	if stack[0] == 0 {
		return 0, nil
	}

	return stack[0], nil
}

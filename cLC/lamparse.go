package main

import (
	"errors"
	"strings"

	"github.com/ElecProg/LamCalc"
)

// parseString turns the input into a Term
func parseString(expr string, globals map[string]LamCalc.Abst) (LamCalc.Term, error) {
	// Backslash is a notation for lambda
	expr = strings.Replace(expr, "\\", "λ", -1)

	return furtherParseString([]rune(expr), map[string]LamCalc.Var{}, globals)
}

func furtherParseString(expr []rune, boundVars map[string]LamCalc.Var, globals map[string]LamCalc.Abst) (LamCalc.Term, error) {
	// Clean string
	expr = []rune(strings.TrimSpace(string(expr)))

	if len(expr) == 0 {
		return nil, errors.New("empty expression")

	} else if expr[0] == 'λ' {
		i := 0

		if len(expr) < 2 {
			return nil, errors.New("no local variable specified in abstraction")
		}

		i++

		// Create copy of boundVars where every index is one higher
		oldVars := boundVars
		boundVars = map[string]LamCalc.Var{}

		// First increment the index of each bound variable
		for variable := range oldVars {
			boundVars[variable] = oldVars[variable] + 1
		}

		// Now get the name of the currently bound variable
		avar := ""

		for ; i < len(expr) && expr[i] != '.'; i++ {
			avar += string(expr[i])
		}

		if !isValidVariableName(avar) {
			return nil, errors.New("invalid variable name '" + avar + "'")

		} else if i >= len(expr)-1 {
			return nil, errors.New("abstraction body not started")
		}

		i++ // Skip the .

		// Add the abstraction variable to the boundvars map
		boundVars[avar] = 0

		// Start of abstraction body, the rest of the expression is part of it
		part, err := furtherParseString(expr[i:], boundVars, globals)

		if err != nil {
			return nil, err
		}

		return LamCalc.Abst{part}, nil
	}

	term := LamCalc.Appl{}

	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case 'λ':
			// Start of abstraction, the rest of the expression is part of it
			part, err := furtherParseString(expr[i:], boundVars, globals)
			i = len(expr)

			if err != nil {
				return nil, err
			}

			term[1] = part

		case '(':
			var cterm LamCalc.Term

			i++
			starte := i

			nBrack := 0
			for ; i < len(expr); i++ {
				if expr[i] == ')' {
					if nBrack == 0 {
						break
					} else {
						nBrack--
					}
				} else if expr[i] == '(' {
					nBrack++
				}
			} // After this loop i points at the closing bracket

			cterm, err := furtherParseString(expr[starte:i], boundVars, globals)

			if err != nil {
				return nil, err
			}

			term[1] = cterm

		case '\t':
			// Skip tabs

		case ' ':
			// Skip spaces

		default:
			// A variable
			cvar := ""

			for ; i < len(expr); i++ {
				if expr[i] == '(' {
					// End of var, take another look at the character later on
					i--
					break
				} else if expr[i] == ' ' {
					// Space end of var
					break
				} else {
					cvar += string(expr[i])
				}
			}

			cindex, ok := boundVars[cvar]

			if ok {
				term[1] = cindex
			} else {
				cfnc, ok := globals[cvar]
				if ok {
					term[1] = cfnc
				} else {
					return nil, errors.New("'" + cvar + "' not yet defined")
				}
			}
		}

		// If the Appl is full: encapsulate it in a new one
		if term[1] != nil {
			if term[0] == nil { // First term wasn't added (happens after first element)
				term = LamCalc.Appl{term[1]}

			} else {
				term = LamCalc.Appl{term}
			}
		}
	}

	// We build the Appl so that there is always one empty spot on the right,
	// thus we only return the first element
	return term[0], nil
}

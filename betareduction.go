package LamCalc

// substitute replaces index by sub
func (lx LamExpr) substitute(index int, sub interface{}) LamExpr {
	nw := LamExpr{}

	for _, term := range lx {
		switch term := term.(type) {
		case int:
			if term == index {
				nw = append(nw, sub)

			} else {
				nw = append(nw, term)
			}

		case LamExpr:
			nw = append(nw, term.substitute(index, sub))

		case LamFunc:
			nw = append(nw, term.substitute(index+1, shiftIndex(1, 0, sub)))
		}
	}

	return nw
}

// substitute replaces index by sub
func (lf LamFunc) substitute(index int, sub interface{}) LamFunc {
	return LamFunc(LamExpr(lf).substitute(index, sub))
}

// betaReduce replaces index 0 by sub and returns a LamExpr
func (lf LamFunc) betaReduce(sub interface{}) LamExpr {
	return shiftIndex(-1, 1, LamExpr(lf).substitute(0, shiftIndex(1, 0, sub))).(LamExpr)
}
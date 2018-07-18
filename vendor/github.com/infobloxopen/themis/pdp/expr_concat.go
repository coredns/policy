package pdp

type functionConcat struct {
	args []Expression
}

func makeFunctionConcat(args []Expression) Expression {
	return functionConcat{
		args: args,
	}
}

func (f functionConcat) GetResultType() Type {
	return TypeListOfStrings
}

func (f functionConcat) describe() string {
	return "concat"
}

// Calculate implements Expression interface and returns calculated value.
func (f functionConcat) Calculate(ctx *Context) (AttributeValue, error) {
	var err error
	out := []string{}
	missCount := 0

	for i, arg := range f.args {
		out, err = appendConcatArg(out, arg, ctx)
		if err != nil {
			if _, ok := err.(*MissingValueError); !ok {
				return UndefinedValue, bindError(bindErrorf(err, "%d", i+1), f.describe())
			}

			missCount++
		}
	}

	if missCount >= len(f.args) {
		return UndefinedValue, bindError(newMissingValueError(), f.describe())
	}

	return MakeListOfStringsValue(out), nil
}

func functionConcatValidator(args []Expression) functionMaker {
	if len(args) <= 0 {
		return nil
	}

	for _, arg := range args {
		t := arg.GetResultType()
		if t != TypeString && t != TypeSetOfStrings && t != TypeListOfStrings {
			if _, ok := t.(*FlagsType); !ok {
				return nil
			}
		}
	}

	return makeFunctionConcat
}

func appendConcatArg(s []string, arg Expression, ctx *Context) ([]string, error) {
	t := arg.GetResultType()
	switch t {
	case TypeString:
		v, err := ctx.calculateStringExpression(arg)
		if err != nil {
			return s, err
		}

		return append(s, v), nil

	case TypeSetOfStrings:
		v, err := ctx.calculateSetOfStringsExpression(arg)
		if err != nil {
			return s, err
		}

		return append(s, SortSetOfStrings(v)...), nil

	case TypeListOfStrings:
		v, err := ctx.calculateListOfStringsExpression(arg)
		if err != nil {
			return s, err
		}

		return append(s, v...), nil
	}

	if t, ok := t.(*FlagsType); ok {
		var n uint64
		switch t.c {
		case 8:
			n8, err := ctx.calculateFlags8Expression(arg)
			if err != nil {
				return s, err
			}

			n = uint64(n8)

		case 16:
			n16, err := ctx.calculateFlags16Expression(arg)
			if err != nil {
				return s, err
			}

			n = uint64(n16)

		case 32:
			n32, err := ctx.calculateFlags32Expression(arg)
			if err != nil {
				return s, err
			}

			n = uint64(n32)

		case 64:
			n64, err := ctx.calculateFlags64Expression(arg)
			if err != nil {
				return s, err
			}

			n = n64
		}

		for i := 0; i < len(t.b); i++ {
			if n&(1<<uint(i)) != 0 {
				s = append(s, t.b[i])
			}
		}

		return s, nil
	}

	return s, newConcatTypeError(t)
}

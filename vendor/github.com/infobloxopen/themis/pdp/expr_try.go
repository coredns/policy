package pdp

type functionTry struct {
	args []Expression
}

func makeFunctionTry(args []Expression) Expression {
	return functionTry{
		args: args,
	}
}

func (f functionTry) GetResultType() Type {
	if len(f.args) > 0 {
		return f.args[0].GetResultType()
	}

	return TypeUndefined
}

func (f functionTry) describe() string {
	return "try"
}

func (f functionTry) Calculate(ctx *Context) (AttributeValue, error) {
	var (
		v   AttributeValue
		err error
	)

	for _, arg := range f.args {
		v, err = arg.Calculate(ctx)
		if err == nil {
			return v, nil
		}
	}

	return UndefinedValue, bindError(err, f.describe())
}

func functionTryValidator(args []Expression) functionMaker {
	if len(args) <= 0 {
		return nil
	}

	t := args[0].GetResultType()
	for _, arg := range args[1:] {
		if arg.GetResultType() != t {
			return nil
		}
	}

	return makeFunctionTry
}

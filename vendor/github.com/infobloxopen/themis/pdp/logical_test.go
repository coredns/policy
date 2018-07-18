package pdp

import "testing"

func TestBooleanNot(t *testing.T) {
	args := []Expression{MakeBooleanValue(true)}

	maker := functionBooleanNotValidator(args)
	if maker == nil {
		t.Error("Expected makeFunctionBooleanNot but got nil")
	}

	e := maker(args)

	rt := e.GetResultType()
	if rt != TypeBoolean {
		t.Errorf("Expected %q type but got %q", TypeBoolean, rt)
	}

	v, err := e.Calculate(nil)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		expected := "false"
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
		} else if s != expected {
			t.Errorf("Expected %q but got %q", expected, s)
		}
	}
}

func TestBooleanOr(t *testing.T) {
	args := []Expression{MakeBooleanValue(false), MakeBooleanValue(false), MakeBooleanValue(true)}

	maker := functionBooleanOrValidator(args)
	if maker == nil {
		t.Error("Expected makeFunctionBooleanNot but got nil")
	}

	e := maker(args)

	rt := e.GetResultType()
	if rt != TypeBoolean {
		t.Errorf("Expected %q type but got %q", TypeBoolean, rt)
	}

	v, err := e.Calculate(nil)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		expected := "true"
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
		} else if s != expected {
			t.Errorf("Expected %q but got %q", expected, s)
		}
	}

	args = []Expression{MakeBooleanValue(false), MakeBooleanValue(false), MakeBooleanValue(false)}
	e = maker(args)
	v, err = e.Calculate(nil)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		expected := "false"
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
		} else if s != expected {
			t.Errorf("Expected %q but got %q", expected, s)
		}
	}
}

func TestBooleanAnd(t *testing.T) {
	args := []Expression{MakeBooleanValue(true), MakeBooleanValue(true), MakeBooleanValue(true)}

	maker := functionBooleanAndValidator(args)
	if maker == nil {
		t.Error("Expected makeFunctionBooleanNot but got nil")
	}

	e := maker(args)

	rt := e.GetResultType()
	if rt != TypeBoolean {
		t.Errorf("Expected %q type but got %q", TypeBoolean, rt)
	}

	v, err := e.Calculate(nil)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		expected := "true"
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
		} else if s != expected {
			t.Errorf("Expected %q but got %q", expected, s)
		}
	}

	args = []Expression{MakeBooleanValue(true), MakeBooleanValue(true), MakeBooleanValue(false)}
	e = maker(args)
	v, err = e.Calculate(nil)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		expected := "false"
		s, err := v.Serialize()
		if err != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
		} else if s != expected {
			t.Errorf("Expected %q but got %q", expected, s)
		}
	}
}

package pdp

import (
	"sort"
	"strings"
	"testing"
)

func TestBuiltinTypes(t *testing.T) {
	for k, bt := range BuiltinTypes {
		if k != bt.GetKey() {
			t.Errorf("expected %q key but got %q", k, bt.GetKey())
		}

		if len(bt.String()) <= 0 {
			t.Errorf("exepcted some human readable name for type %q but got empty string", k)
		}

		if !bt.Match(bt) {
			t.Errorf("expected that %q matches itself", bt)
		}
	}

	if TypeBoolean.Match(TypeString) {
		t.Errorf("expected that %q doesn't match %q", TypeBoolean, TypeString)
	}
}

func TestFlagsType(t *testing.T) {
	flags8Name := "8Flags"
	flags8 := []string{
		"f00", "f01", "f02", "f03", "f04", "f05", "f06",
	}
	ft8, err := NewFlagsType(flags8Name, flags8...)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		if ft8, ok := ft8.(*FlagsType); ok {
			n := ft8.String()
			if n != flags8Name {
				t.Errorf("Expected %q as type name but got %q", flags8Name, n)
			}

			k := ft8.GetKey()
			if k != strings.ToLower(flags8Name) {
				t.Errorf("Expected %q as type key but got %q", strings.ToLower(flags8Name), k)
			}

			c := ft8.Capacity()
			if c != 8 {
				t.Errorf("Expected 8 bit as type capacity but got %d", c)
			}

			// Flags names use octal system
			b := ft8.GetFlagBit("f06")
			if b != 06 {
				t.Errorf("Expected 006 as bit number for %q but got %03o", "f06", b)
			}

			b = ft8.GetFlagBit("f07")
			if b != -1 {
				t.Errorf("Expected no flag %q (-1) but got %03o", "f07", b)
			}

			assertMapStringIntKeys(ft8.f, flags8, "flags8 index", t)
			assertStrings(ft8.b, flags8, "flags8 names", t)
		} else {
			t.Errorf("Expected *FlagsType but got %T", ft8)
		}

		if !ft8.Match(ft8) {
			t.Errorf("expected that %q matches itself", ft8)
		}

		if ft8.Match(TypeBoolean) {
			t.Errorf("expected that %q doesn't match %q", ft8, TypeBoolean)
		}
	}

	oft86, err := NewFlagsType("Other86Flags", flags8...)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if !oft86.Match(ft8) {
		t.Errorf("expected that %q matches %q", oft86, ft8)
	}

	oft87, err := NewFlagsType("Other87Flags",
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
	)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if oft87.Match(ft8) {
		t.Errorf("expected that %q doesn't match %q", oft87, ft8)
	}

	flags16Name := "16Flags"
	flags16 := []string{
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16",
	}
	ft16, err := NewFlagsType(flags16Name, flags16...)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if ft16, ok := ft16.(*FlagsType); ok {
		n := ft16.String()
		if n != flags16Name {
			t.Errorf("Expected %q as type name but got %q", flags16Name, n)
		}

		k := ft16.GetKey()
		if k != strings.ToLower(flags16Name) {
			t.Errorf("Expected %q as type key but got %q", strings.ToLower(flags16Name), k)
		}

		c := ft16.Capacity()
		if c != 16 {
			t.Errorf("Expected 16 bit as type capacity but got %d", c)
		}

		// Flags names use octal system
		b := ft16.GetFlagBit("f16")
		if b != 016 {
			t.Errorf("Expected 016 as bit number for %q but got %03o", "f16", b)
		}

		b = ft16.GetFlagBit("f17")
		if b != -1 {
			t.Errorf("Expected no flag %q (-1) but got %03o", "f17", b)
		}

		assertMapStringIntKeys(ft16.f, flags16, "flags16 index", t)
		assertStrings(ft16.b, flags16, "flags16 names", t)
	} else {
		t.Errorf("Expected *FlagsType but got %T", ft16)
	}

	flags32Name := "32Flags"
	flags32 := []string{
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36",
	}
	ft32, err := NewFlagsType(flags32Name, flags32...)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if ft32, ok := ft32.(*FlagsType); ok {
		n := ft32.String()
		if n != flags32Name {
			t.Errorf("Expected %q as type name but got %q", flags32Name, n)
		}

		k := ft32.GetKey()
		if k != strings.ToLower(flags32Name) {
			t.Errorf("Expected %q as type key but got %q", strings.ToLower(flags32Name), k)
		}

		c := ft32.Capacity()
		if c != 32 {
			t.Errorf("Expected 32 bit as type capacity but got %d", c)
		}

		// Flags names use octal system
		b := ft32.GetFlagBit("f36")
		if b != 036 {
			t.Errorf("Expected 036 as bit number for %q but got %03o", "f36", b)
		}

		b = ft32.GetFlagBit("f37")
		if b != -1 {
			t.Errorf("Expected no flag %q (-1) but got %03o", "f37", b)
		}

		assertMapStringIntKeys(ft32.f, flags32, "flags32 index", t)
		assertStrings(ft32.b, flags32, "flags32 names", t)
	} else {
		t.Errorf("Expected *FlagsType but got %T", ft32)
	}

	flags64Name := "64Flags"
	flags64 := []string{
		"f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
		"f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
		"f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
		"f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
		"f70", "f71", "f72", "f73", "f74", "f75", "f76",
	}
	ft64, err := NewFlagsType(flags64Name, flags64...)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else if ft64, ok := ft64.(*FlagsType); ok {
		n := ft64.String()
		if n != flags64Name {
			t.Errorf("Expected %q as type name but got %q", flags64Name, n)
		}

		k := ft64.GetKey()
		if k != strings.ToLower(flags64Name) {
			t.Errorf("Expected %q as type key but got %q", strings.ToLower(flags64Name), k)
		}

		c := ft64.Capacity()
		if c != 64 {
			t.Errorf("Expected 64 bit as type capacity but got %d", c)
		}

		// Flags names use octal system
		b := ft64.GetFlagBit("f76")
		if b != 076 {
			t.Errorf("Expected 076 as bit number for %q but got %03o", "f76", b)
		}

		b = ft64.GetFlagBit("f77")
		if b != -1 {
			t.Errorf("Expected no flag %q (-1) but got %03o", "f77", b)
		}

		assertMapStringIntKeys(ft64.f, flags64, "flags64 index", t)
		assertStrings(ft64.b, flags64, "flags64 names", t)
	} else {
		t.Errorf("Expected *FlagsType but got %T", ft64)
	}

	ftBool, err := NewFlagsType("Boolean", "false", "true")
	if err == nil {
		t.Errorf("Expected error but got type %q", ftBool)
	} else if _, ok := err.(*duplicatesBuiltinTypeError); !ok {
		t.Errorf("Expected *duplicatesBuiltinTypeError error but got %T: %s", err, err)
	}

	ftEmpty, err := NewFlagsType("Empty")
	if err == nil {
		t.Errorf("Expected error but got type %q", ftEmpty)
	} else if _, ok := err.(*noFlagsDefinedError); !ok {
		t.Errorf("Expected *noFlagsDefinedError error but got %T: %s", err, err)
	}

	ft128, err := NewFlagsType("128Flags",
		"f000", "f001", "f002", "f003", "f004", "f005", "f006", "f007",
		"f010", "f011", "f012", "f013", "f014", "f015", "f016", "f017",
		"f020", "f021", "f022", "f023", "f024", "f025", "f026", "f027",
		"f030", "f031", "f032", "f033", "f034", "f035", "f036", "f037",
		"f040", "f041", "f042", "f043", "f044", "f045", "f046", "f047",
		"f050", "f051", "f052", "f053", "f054", "f055", "f056", "f057",
		"f060", "f061", "f062", "f063", "f064", "f065", "f066", "f067",
		"f070", "f071", "f072", "f073", "f074", "f075", "f076", "f077",
		"f100", "f101", "f102", "f103", "f104", "f105", "f106", "f107",
		"f110", "f111", "f112", "f113", "f114", "f115", "f116", "f117",
		"f120", "f121", "f122", "f123", "f124", "f125", "f126", "f127",
		"f130", "f131", "f132", "f133", "f134", "f135", "f136", "f137",
		"f140", "f141", "f142", "f143", "f144", "f145", "f146", "f147",
		"f150", "f151", "f152", "f153", "f154", "f155", "f156", "f157",
		"f160", "f161", "f162", "f163", "f164", "f165", "f166", "f167",
		"f170", "f171", "f172", "f173", "f174", "f175", "f176", "f177",
	)
	if err == nil {
		t.Errorf("Expected error but got type %q", ft128)
	} else if _, ok := err.(*tooManyFlagsDefinedError); !ok {
		t.Errorf("Expected *tooManyFlagsDefinedError error but got %T: %s", err, err)
	}

	ftDup, err := NewFlagsType("Dup",
		"f00", "f01", "dup", "f03", "f04", "f05", "f06", "f07",
		"f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
		"f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
		"f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
		"f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
		"f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
		"f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
		"f70", "f71", "f72", "f73", "f74", "DUP", "f76", "f77",
	)
	if err == nil {
		t.Errorf("Expected error but got type %q", ftDup)
	} else if _, ok := err.(*duplicateFlagName); !ok {
		t.Errorf("Expected *duplicateFlagName error but got %T: %s", err, err)
	}
}

func TestSignature(t *testing.T) {
	sign := MakeSignature(
		TypeUndefined,
		TypeBoolean,
		TypeString,
		TypeInteger,
		TypeFloat,
		TypeAddress,
		TypeNetwork,
		TypeDomain,
		TypeSetOfStrings,
		TypeSetOfNetworks,
		TypeSetOfDomains,
		TypeListOfStrings,
	)

	e := "\"Undefined\"/" +
		"\"Boolean\"/" +
		"\"String\"/" +
		"\"Integer\"/" +
		"\"Float\"/" +
		"\"Address\"/" +
		"\"Network\"/" +
		"\"Domain\"/" +
		"\"Set of Strings\"/" +
		"\"Set of Networks\"/" +
		"\"Set of Domains\"/" +
		"\"List of Strings\""
	s := sign.String()
	if s != e {
		t.Errorf("expected %s signature but got %s", e, s)
	}

	sign = MakeSignature()
	e = "empty"
	s = sign.String()
	if s != e {
		t.Errorf("expected %s signature but got %s", e, s)
	}
}

func TestTypeSet(t *testing.T) {
	set := makeTypeSet(
		TypeUndefined,
		TypeBoolean,
		TypeString,
		TypeAddress,
		TypeNetwork,
		TypeDomain,
		TypeSetOfStrings,
		TypeSetOfNetworks,
		TypeSetOfDomains,
		TypeListOfStrings,
	)

	e := "\"Address\", " +
		"\"Boolean\", " +
		"\"Domain\", " +
		"\"List of Strings\", " +
		"\"Network\", " +
		"\"Set of Domains\", " +
		"\"Set of Networks\", " +
		"\"Set of Strings\", " +
		"\"String\", " +
		"\"Undefined\""
	s := set.String()
	if s != e {
		t.Errorf("expected %s signature but got %s", e, s)
	}

	if !set.Contains(TypeAddress) {
		t.Errorf("expected %q in the set but it isn't here", TypeAddress)
	}

	if set.Contains(TypeInteger) {
		t.Errorf("expected %q to be not in the set but it is here", TypeInteger)
	}

	set = makeTypeSet()
	e = "empty"
	s = set.String()
	if s != e {
		t.Errorf("expected %s signature but got %s", e, s)
	}
}

func assertMapStringIntKeys(v map[string]int, e []string, desc string, t *testing.T) {
	vs := make([]string, len(v))
	i := 0
	for s := range v {
		vs[i] = s
		i++
	}
	sort.Strings(vs)

	assertStrings(vs, e, desc, t)
}

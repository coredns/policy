package pdp

// AttributeDesignator represents an expression which result is corresponding
// attribute value from request context.
type AttributeDesignator struct {
	a Attribute
}

// MakeAttributeDesignator creates designator expression instance for given
// attribute.
func MakeAttributeDesignator(a Attribute) AttributeDesignator {
	return AttributeDesignator{a}
}

// MakeDesignator creates designator expression instance for given attribute id
// and type.
func MakeDesignator(id string, t Type) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, t))
}

// MakeBooleanDesignator creates boolean designator expression instance for
// given attribute id.
func MakeBooleanDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeBoolean))
}

// MakeStringDesignator creates boolean designator expression instance for given
// attribute id.
func MakeStringDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeString))
}

// MakeIntegerDesignator creates boolean designator expression instance for
// given attribute id.
func MakeIntegerDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeInteger))
}

// MakeFloatDesignator creates boolean designator expression instance for given
// attribute id.
func MakeFloatDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeFloat))
}

// MakeAddressDesignator creates boolean designator expression instance for
// given attribute id.
func MakeAddressDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeAddress))
}

// MakeNetworkDesignator creates boolean designator expression instance for
// given attribute id.
func MakeNetworkDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeNetwork))
}

// MakeDomainDesignator creates boolean designator expression instance for given
// attribute id.
func MakeDomainDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeDomain))
}

// MakeSetOfStringsDesignator creates boolean designator expression instance for
// given attribute id.
func MakeSetOfStringsDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeSetOfStrings))
}

// MakeSetOfNetworksDesignator creates boolean designator expression instance
// for given attribute id.
func MakeSetOfNetworksDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeSetOfNetworks))
}

// MakeSetOfDomainsDesignator creates boolean designator expression instance for
// given attribute id.
func MakeSetOfDomainsDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeSetOfDomains))
}

// MakeListOfStringsDesignator creates boolean designator expression instance
// for given attribute id.
func MakeListOfStringsDesignator(id string) AttributeDesignator {
	return MakeAttributeDesignator(MakeAttribute(id, TypeListOfStrings))
}

// GetID returns ID of wrapped attribute.
func (d AttributeDesignator) GetID() string {
	return d.a.id
}

// GetResultType returns type of wrapped attribute (implements Expression
// interface).
func (d AttributeDesignator) GetResultType() Type {
	return d.a.t
}

// Calculate implements Expression interface and returns calculated value
func (d AttributeDesignator) Calculate(ctx *Context) (AttributeValue, error) {
	return ctx.getAttribute(d.a)
}

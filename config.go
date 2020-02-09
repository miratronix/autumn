package autumn

import "unicode"

// config defines the configuration structure for autumn
type config struct {
	tagName             string
	leafNameMethod      string
	postConstructMethod string
	preDestroyMethod    string
}

// NewConfig creates a new configuration object
func NewConfig() *config {
	return &config{
		tagName:             "autumn",
		leafNameMethod:      "GetLeafName",
		postConstructMethod: "PostConstruct",
		preDestroyMethod:    "PreDestroy",
	}
}

// TagName sets the autumn tag name
func (c *config) TagName(tag string) *config {
	if len(tag) == 0 {
		panic("The tag name cannot be empty")
	}
	c.tagName = tag
	return c
}

// LeafNameMethod sets the method name for getting the leaf name
func (c *config) LeafNameMethod(method string) *config {
	c.ensurePublicMethod(method)
	c.leafNameMethod = method
	return c
}

// PostConstructMethod sets the method name for post construct calls
func (c *config) PostConstructMethod(method string) *config {
	c.ensurePublicMethod(method)
	c.postConstructMethod = method
	return c
}

// PreDestroyMethod sets the method name for pre destroy calls
func (c *config) PreDestroyMethod(method string) *config {
	c.ensurePublicMethod(method)
	c.preDestroyMethod = method
	return c
}

// ensurePublicMethod ensures the supplied method name is public
func (c *config) ensurePublicMethod(method string) {

	// Make sure it's not an empty string
	if len(method) == 0 {
		panic("The method name cannot be empty")
	}

	// Convert to a slice of runes
	runes := []rune(method)

	// Make sure the first character is uppercase
	if !unicode.IsUpper(runes[0]) {
		panic("The method name must be public (start with an uppercase character)")
	}
}

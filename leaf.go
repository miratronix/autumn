package autumn

import (
	"reflect"
)

// leaf describes a single injected class
type leaf struct {
	structureType    reflect.Type
	structureValue   reflect.Value
	structureElement reflect.Value

	name          string
	postConstruct reflect.Value
	preDestroy    reflect.Value

	unresolvedDependencies map[string]reflect.Value
	resolvedDependencies   map[string]reflect.Value
}

// newLeaf constructs a new leaf, using the structure name as the name
func newLeaf(config *config, structurePointer interface{}) *leaf {
	leaf := &leaf{
		structureType:    getStructureType(structurePointer),
		structureValue:   getStructureValue(structurePointer),
		structureElement: getStructureElement(structurePointer),
	}

	leaf.initializeName(config.leafNameMethod)
	leaf.initializeDependencies(config.tagName)
	leaf.initializePostConstruct(config.postConstructMethod)
	leaf.initializePreDestroy(config.preDestroyMethod)

	return leaf
}

// newNamedLeaf constructs a new leaf with the specified name
func newNamedLeaf(config *config, name string, structurePointer interface{}) *leaf {
	leaf := &leaf{
		structureType:    getStructureType(structurePointer),
		structureValue:   getStructureValue(structurePointer),
		structureElement: getStructureElement(structurePointer),
		name:             name,
	}

	leaf.initializeDependencies(config.tagName)
	leaf.initializePostConstruct(config.postConstructMethod)
	leaf.initializePreDestroy(config.preDestroyMethod)

	return leaf
}

// initializeName initializes the name for the leaf
func (l *leaf) initializeName(getNameMethod string) {

	method := l.structureValue.MethodByName(getNameMethod)
	if !method.IsValid() {
		l.name = l.structureType.String()
		return
	}

	if method.Type().NumIn() != 0 {
		panic(l.structureType.String() + " - " + getNameMethod + " must not take any parameters")
	} else if method.Type().NumOut() != 1 {
		panic(l.structureType.String() + " - " + getNameMethod + " must return exactly one parameter")
	} else if method.Type().Out(0).Kind() != reflect.String {
		panic(l.structureType.String() + " - " + getNameMethod + " must return a string")
	}

	l.name = method.Call([]reflect.Value{})[0].String()
}

// initializeDependencies reads in structure tags to find dependencies
func (l *leaf) initializeDependencies(tagName string) {
	l.unresolvedDependencies = map[string]reflect.Value{}
	l.resolvedDependencies = map[string]reflect.Value{}

	for i := 0; i < l.structureType.NumField(); i++ {
		field := l.structureType.Field(i)
		dep := field.Tag.Get(tagName)
		if len(dep) != 0 {
			l.unresolvedDependencies[dep] = l.structureElement.FieldByName(field.Name)
		}
	}
}

// initializePostConstruct initializes the post construct function for the leaf, panicking if it's invalid
func (l *leaf) initializePostConstruct(postConstructMethod string) {
	l.postConstruct = l.structureValue.MethodByName(postConstructMethod)
	if !l.postConstruct.IsValid() {
		return
	}

	if l.postConstruct.Type().NumIn() != 0 {
		panic(l.structureType.String() + " - " + postConstructMethod + " must not take any parameters")
	} else if l.postConstruct.Type().NumOut() != 0 {
		panic(l.structureType.String() + " - " + postConstructMethod + " must not return any parameters")
	}
}

// initializePreDestroy initializes the pre destroy function for the leaf, panicking if it's invalid
func (l *leaf) initializePreDestroy(preDestroyMethod string) {
	l.preDestroy = l.structureValue.MethodByName(preDestroyMethod)
	if !l.preDestroy.IsValid() {
		return
	}

	if l.preDestroy.Type().NumIn() != 0 {
		panic(l.structureType.String() + " - " + preDestroyMethod + " must not take any parameters")
	} else if l.preDestroy.Type().NumOut() != 0 {
		panic(l.structureType.String() + " - " + preDestroyMethod + " must not return any parameters")
	}
}

// resolveDependencies resolves dependencies for the leaf using the supplied tree
func (l *leaf) resolveDependencies(tree *Tree) {
	for name := range l.unresolvedDependencies {
		dep := tree.GetLeaf(name)
		if dep != nil {
			l.setDependency(tree, name, dep)
		}
	}
}

// setDependency sets a dependency in the leaf
func (l *leaf) setDependency(tree *Tree, name string, leaf *leaf) {
	if !l.unresolvedDependencies[name].IsValid() {
		panic("Can't set dependency " + name + "in leaf" + l.name)
	}

	// Set the dependency and move it to "resolved"
	l.unresolvedDependencies[name].Set(leaf.structureValue)
	l.resolvedDependencies[name] = l.unresolvedDependencies[name]
	delete(l.unresolvedDependencies, name)
}

// dependenciesResolved determines if dependencies have been resolved
func (l *leaf) dependenciesResolved() bool {
	return len(l.unresolvedDependencies) == 0
}

// callPostConstruct calls the leaf's PostConstruct method if it has one
func (l *leaf) callPostConstruct() {
	if !l.postConstruct.IsValid() {
		return
	}
	l.postConstruct.Call([]reflect.Value{})
}

// callPreDestroy calls the leaf's PreDestroy method if it has one
func (l *leaf) callPreDestroy() {
	if !l.preDestroy.IsValid() {
		return
	}
	l.preDestroy.Call([]reflect.Value{})
}

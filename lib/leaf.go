package lib

import (
	"reflect"
)

const (
	tagName             = "autumn"
	getNameMethod       = "GetLeafName"
	postConstructMethod = "PostConstruct"
	preDestroyMethod    = "PreDestroy"
)

// Leaf describes a single injected class
type Leaf struct {
	structureType    reflect.Type
	structureValue   reflect.Value
	structureElement reflect.Value
	structurePointer uintptr

	name          string
	aliases       map[string]struct{}
	value         interface{}
	postConstruct reflect.Value
	preDestroy    reflect.Value

	unresolvedDependencies map[string]reflect.Value
	resolvedDependencies   map[string]reflect.Value
}

// NewLeaf constructs a new leaf, using the structure name as the name
func NewLeaf(structPtr interface{}) *Leaf {
	leaf := &Leaf{
		structureType:    getStructureType(structPtr),
		structureValue:   getStructureValue(structPtr),
		structureElement: getStructureElement(structPtr),
		structurePointer: getStructurePointer(structPtr),
		aliases:          make(map[string]struct{}),
	}

	leaf.initializeName()
	leaf.initializeDependencies()
	leaf.initializePostConstruct()
	leaf.initializePreDestroy()

	leaf.addAlias(leaf.name)

	return leaf
}

// NewNamedLeaf constructs a new leaf with the specified name
func NewNamedLeaf(name string, structPtr interface{}) *Leaf {
	leaf := &Leaf{
		structureType:    getStructureType(structPtr),
		structureValue:   getStructureValue(structPtr),
		structureElement: getStructureElement(structPtr),
		structurePointer: getStructurePointer(structPtr),
		name:             name,
		aliases:          make(map[string]struct{}),
	}

	leaf.initializeDependencies()
	leaf.initializePostConstruct()
	leaf.initializePreDestroy()

	leaf.addAlias(leaf.name)

	return leaf
}

// addAlias adds the supplied name to the leaf's alias' list
func (l *Leaf) addAlias(name string) {
	l.aliases[name] = struct{}{}
}

// hasAlias checks if the leaf has the supplied name as an alias
func (l *Leaf) hasAlias(name string) bool {
	_, ok := l.aliases[name]
	return ok
}

// initializeName initializes the name for the leaf
func (l *Leaf) initializeName() {

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

// initializeDependencies read in structure tags to find dependencies
func (l *Leaf) initializeDependencies() {
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
func (l *Leaf) initializePostConstruct() {
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
func (l *Leaf) initializePreDestroy() {
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
func (l *Leaf) resolveDependencies(tree *Tree) {
	for name := range l.unresolvedDependencies {
		dep := tree.GetLeaf(name)
		if dep != nil {
			l.setDependency(tree, name, dep)
		}
	}

	// All resolved
	if l.dependenciesResolved() {
		l.callPostConstruct()
	}
}

// setDependency sets a dependency in the leaf
func (l *Leaf) setDependency(tree *Tree, name string, leaf *Leaf) {
	if !l.unresolvedDependencies[name].IsValid() {
		panic("Can't set dependency " + name + "in leaf" + l.name)
	}

	// Set the dependency and move it to "resolved"
	l.unresolvedDependencies[name].Set(leaf.structureValue)
	l.resolvedDependencies[name] = l.unresolvedDependencies[name]
	delete(l.unresolvedDependencies, name)
}

// dependenciesResolved determines if dependencies have been resolved
func (l *Leaf) dependenciesResolved() bool {
	return len(l.unresolvedDependencies) == 0
}

// callPostConstruct calls the leaf's PostConstruct method if it has one
func (l *Leaf) callPostConstruct() {
	if !l.postConstruct.IsValid() {
		return
	}
	l.postConstruct.Call([]reflect.Value{})
}

// callPreDestroy calls the leaf's PreDestroy method if it has one
func (l *Leaf) callPreDestroy() {
	if !l.preDestroy.IsValid() {
		return
	}
	l.preDestroy.Call([]reflect.Value{})
}

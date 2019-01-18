package lib

import "reflect"

// isStructurePointer determines if the supplied value is a structure pointer
func isStructurePointer(data interface{}) bool {
	isPtr := reflect.ValueOf(data).Kind() == reflect.Ptr
	if !isPtr {
		return false
	}
	return reflect.ValueOf(data).Elem().Kind() == reflect.Struct
}

// getStructureType gets the type of the supplied structure
func getStructureType(data interface{}) reflect.Type {
	return reflect.ValueOf(data).Elem().Type()
}

// getStructureValue gets the reflection value of the supplied structure, which is usually a pointer to the structure
func getStructureValue(data interface{}) reflect.Value {
	return reflect.ValueOf(data)
}

// getStructureElement dereferences the structure pointer and gets a reflection value for the underlying structure
func getStructureElement(data interface{}) reflect.Value {
	return reflect.ValueOf(data).Elem()
}

// getStructurePointer gets the pointer of the supplied structure
func getStructurePointer(data interface{}) uintptr {
	return reflect.ValueOf(data).Pointer()
}

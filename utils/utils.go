package utils

import (
	"math"
	"reflect"
	"strings"
)

func MaxValue[T any]() T {
	var v T
	t := reflect.TypeOf(v)
	zero := reflect.Zero(t)
	max := reflect.New(t).Elem()
	max.Set(zero)

	switch t.Kind() {
	case reflect.Int:
		max.SetInt(math.MaxInt)
	case reflect.Int8:
		max.SetInt(math.MaxInt8)
	case reflect.Int16:
		max.SetInt(math.MaxInt16)
	case reflect.Int32:
		max.SetInt(math.MaxInt32)
	case reflect.Int64:
		max.SetInt(math.MaxInt64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		max.SetUint(^zero.Uint())
	case reflect.Float32:
		max.SetFloat(math.MaxFloat32)
	case reflect.Float64:
		max.SetFloat(math.MaxFloat64)
	case reflect.String:
		// Building Max value of string would be very space-consuming. This satisfies majority map use cases
		max.SetString(strings.Repeat(string(rune(math.MaxInt32)), 255))
	default:
		panic("trying to retrieve max value of unsupported type")
	}

	return max.Interface().(T)
}

func MinValue[T any]() T {
	var v T
	t := reflect.TypeOf(v)
	zero := reflect.Zero(t)
	max := reflect.New(t).Elem()
	max.Set(zero)

	switch t.Kind() {
	case reflect.Int:
		max.SetInt(math.MinInt)
	case reflect.Int8:
		max.SetInt(math.MinInt8)
	case reflect.Int16:
		max.SetInt(math.MinInt16)
	case reflect.Int32:
		max.SetInt(math.MinInt32)
	case reflect.Int64:
		max.SetInt(math.MinInt64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		max.SetUint(zero.Uint())
	case reflect.Float32:
		max.SetFloat(-math.MaxFloat32)
	case reflect.Float64:
		max.SetFloat(-math.MaxFloat64)
	case reflect.String:
		max.SetString("") //Assuming it's the minimal possible string.
	default:
		panic("trying to retrieve max value of unsupported type")
	}

	return max.Interface().(T)
}

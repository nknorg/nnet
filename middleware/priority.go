package middleware

import (
	"reflect"
	"sort"
)

func getPriority(v reflect.Value) int32 {
	return int32(v.FieldByName("Priority").Int())
}

// Sort sorts an array/slice of middleware
func Sort(middlewares interface{}) {
	sort.SliceStable(middlewares, func(i int, j int) bool {
		s := reflect.ValueOf(middlewares)
		if s.Kind() != reflect.Slice {
			panic("Sort input is not a slice")
		}
		return getPriority(s.Index(i)) > getPriority(s.Index(j))
	})
}

package reflectutil

import "reflect"

// AssertNoZeroFields panics if v has any zero fields.
func AssertNoZeroFields(v any, assertionFail func(field string)) {
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		panic("reflectutil: expected struct")
	}

	for i := 0; i < rt.NumField(); i++ {
		fv := rv.Field(i)
		if fv.IsZero() {
			if assertionFail != nil {
				assertionFail(rt.Field(i).Name)
			} else {
				panic("reflectutil: zero field " + rt.Field(i).Name + " in " + rt.String())
			}
		}
	}
}

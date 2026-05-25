package luautil


import  (
	"github.com/yuin/gopher-lua"
	"reflect"
)

func requiredString(L *lua.LState, tbl *lua.LTable, key string) string {
    val := tbl.RawGetString(key)

    if val.Type() == lua.LTNil {
        L.RaiseError("missing required field: '%s'", key)
    }
    if val.Type() != lua.LTString {
        L.RaiseError("field '%s' must be a string, got %s", key, val.Type())
    }
    if val.String() == "" {
        L.RaiseError("field '%s' cannot be empty", key)
    }

    return val.String()
}

func optionalString(L *lua.LState, tbl *lua.LTable, key string, fallback string) string {

    val := tbl.RawGetString(key)

    if val.Type() == lua.LTNil {
        return fallback
    }

	if val.Type() != lua.LTString {
        L.RaiseError("field '%s' must be a string, got %s", key, val.Type())
    }


    return val.String()
}



func allowedKeysFromStruct(v any) []string {
    t := reflect.TypeOf(v)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    var keys []string
    for i := range t.NumField() {
        tag := t.Field(i).Tag.Get("lua")
        if tag != "" {
            keys = append(keys, tag)
        }
    }
    return keys
}


func checkNoExtraKeys(L *lua.LState, tbl *lua.LTable, allowed []string) {
    set := make(map[string]bool, len(allowed))
    for _, k := range allowed {
        set[k] = true
    }
    tbl.ForEach(func(key lua.LValue, _ lua.LValue) {
        if key.Type() != lua.LTString {
            L.RaiseError("unexpected non-string key")
        }
        if !set[key.String()] {
            L.RaiseError("unknown field: '%s'", key.String())
        }
    })
}
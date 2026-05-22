package luautil

import  "github.com/yuin/gopher-lua"

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
package parce


import (
	"fmt"
	"strings"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"reflect"

 )



 
func validateStringFields(target any, body *hclsyntax.Body) hcl.Diagnostics {
    var diags hcl.Diagnostics

    v := reflect.ValueOf(target)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    t := v.Type()

    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
        if !f.IsExported() || f.Tag.Get("validate") != "nonempty" {
            continue
        }

        // "mirror,optional" -> "mirror"
        name, _, _ := strings.Cut(f.Tag.Get("hcl"), ",")
        rng := attrRangeOf(body, name)

        switch fv := v.Field(i).Interface().(type) {
        case string:
            diags = append(diags, checkStringField(fv, name, rng)...)
        case *string:
            diags = append(diags, checkStringField(fv, name, rng)...)
        default:
            panic(fmt.Sprintf("field %s: validate:\"nonempty\" only supports string or *string", f.Name))
        }
    }
    return diags
}




func checkStringField[T string | *string](value T, field string, rng hcl.Range) hcl.Diagnostics {



    switch v := any(value).(type) {

    case *string:
        if v == nil {
            return nil
        }
        if strings.TrimSpace(*v) == "" {
            return hcl.Diagnostics{optionalEmptyErr(field, rng)}
        }
    case string:
        if strings.TrimSpace(v) == "" {
            return hcl.Diagnostics{requiredErr(field, rng)}
        }
    }
    return nil
}








func optionalEmptyErr(field string, rng hcl.Range) *hcl.Diagnostic {

	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  fmt.Sprintf("%s must not be empty if provided", field),
		Detail:   fmt.Sprintf(
			"omit %s entirely to use the default, or provide a non-empty value",
			field,
		),
		Subject: rng.Ptr(),
	}
}



func requiredErr(field string, rng hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  fmt.Sprintf("%s must not be empty", field),
		Subject:  rng.Ptr(),
	}
}





func checkRequiredString(value, field string, attrRange hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if strings.TrimSpace(value) == "" {
		diags = append(diags, requiredErr(field, attrRange))
		return diags
	}
	return diags
}




func checkOptionalString(value *string, field string, attrRange hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if value == nil {
		return diags // reduent check but it does not hurt 
	}

	if strings.TrimSpace(*value) == "" {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("%s must not be empty if provided", field),
			Detail:   fmt.Sprintf("omit %s entirely to use the default, or provide a non-empty value", field),
			Subject:  attrRange.Ptr(),
		})
		return diags
	}

	return diags
}








func attrRangeOf(body *hclsyntax.Body, name string) hcl.Range {
    if attr, ok := body.Attributes[name]; ok {
        return attr.Expr.Range()
    }
    return body.EndRange
}




func derefOr(p *string, def string) string {
	if p == nil {
		return def
	}
	return *p
}


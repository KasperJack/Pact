package parce


import (
	"fmt"
	"strings"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"reflect"
 )








func checkField(value any, field string, rng hcl.Range) hcl.Diagnostics {
    v := reflect.ValueOf(value)

    if v.Kind() == reflect.Ptr {
        if v.IsNil() {
            return nil // not provided that's fine
        }
        v = v.Elem() // dereference to get the actual string
        if strings.TrimSpace(v.String()) == "" {
            return hcl.Diagnostics{&hcl.Diagnostic{
                Severity: hcl.DiagError,
                Summary:  fmt.Sprintf("%s must not be empty if provided", field),
                Detail:   fmt.Sprintf("omit %s entirely to use the default, or provide a non-empty value", field),
                Subject:  rng.Ptr(),
            }}
        }
        return nil
    }

    if strings.TrimSpace(v.String()) == "" {
        return hcl.Diagnostics{requiredErr(field, rng)}
    }
    return nil
}


func requiredErr(field string, rng hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  fmt.Sprintf("%s must not be empty", field),
		Subject:  rng.Ptr(),
	}
}


func checkRequired(value, field string, attrRange hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if strings.TrimSpace(value) == "" {
		diags = append(diags, requiredErr(field, attrRange))
		return diags
	}
	return diags
}

func checkOptional(value *string, field string, attrRange hcl.Range) hcl.Diagnostics {
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


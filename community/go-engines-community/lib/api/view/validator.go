package view

import (
	"github.com/go-playground/validator/v10"
)

func ValidateEditPositionRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditPositionRequest)

	if len(r.Items) > 0 {
		exists := make(map[string]bool, len(r.Items))
		existsView := make(map[string]bool, len(r.Items))
		for _, item := range r.Items {
			if exists[item.ID] {
				sl.ReportError(r.Items, "Items", "Item", "has_duplicates", "")
				return
			}

			exists[item.ID] = true

			for _, view := range item.Views {
				if existsView[view] {
					sl.ReportError(r.Items, "Items", "Item", "has_duplicates", "")
					return
				}

				existsView[view] = true
			}
		}
	}
}

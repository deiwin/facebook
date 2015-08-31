package model

import "net/url"

// Formable represents objects that can be encoded into HTTP form values
type Formable interface {
	AsForm() url.Values
}

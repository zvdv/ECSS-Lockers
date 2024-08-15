package router

import (
	"testing"
)

func TestFormValidate(t *testing.T) {
	validEmail := []string{
		"foobar@uvic.ca",
		"Goobarba123z@uvic.ca",
		"Goobarba_123z@uvic.ca",
	}

	for _, email := range validEmail {
		if !uvicEmailValidator(email) {
			t.Fail()
		}
	}

	invalidEmails := []string{
		"foobar@uvic.caa",
		"Goobarba123z@uvic.com",
		"Goobarba_123z@gmail.uk",
	}

	for _, email := range invalidEmails {
		if uvicEmailValidator(email) {
			t.Fail()
		}
	}
}

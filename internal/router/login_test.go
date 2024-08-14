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

	invalidEmail := []string{
		"foobar@uvic.c",
		"Goobarba123z@uvic.a",
		"Goobarba_123z@gmail.ca",
	}

	for _, email := range invalidEmail {
		if uvicEmailValidator(email) {
			t.Fail()
		}
	}
}

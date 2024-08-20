package email_test

import (
	"testing"

	"github.com/zvdv/ECSS-Lockers/internal/email"
)

func TestFormValidate(t *testing.T) {
	validEmail := []string{
		"foobar@uvic.ca",
		"Goobarba123z@uvic.ca",
		"Goobarba_123z@uvic.ca",
	}

	for _, addr := range validEmail {
		if !email.ValidUVicEmail(addr) {
			t.Fatal(addr)
		}
	}

	invalidEmails := []string{
		"foobar@uvic.caa",
		"Goobarba123z@uvic.com",
		"Goobarba_123z@gmail.uk",
	}

	for _, addr := range invalidEmails {
		if email.ValidUVicEmail(addr) {
			t.Fatal(addr)
		}
	}
}

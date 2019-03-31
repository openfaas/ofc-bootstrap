package types

import (
	"testing"
)

func TestFormatCustomTemplates_None(t *testing.T) {

	d := Deployment{
		CustomTemplate: []string{},
	}

	want := ""
	got := d.FormatCustomTemplates()
	if got != want {
		t.Errorf("FormatCustomTemplates want: %s, got %s", want, got)
		t.Fail()
	}
}

func TestFormatCustomTemplates_Single(t *testing.T) {

	d := Deployment{
		CustomTemplate: []string{"https://w.com/repo1"},
	}

	want := d.CustomTemplate[0]
	got := d.FormatCustomTemplates()
	if got != want {
		t.Errorf("FormatCustomTemplates want: %s, got %s", want, got)
		t.Fail()
	}
}

func TestFormatCustomTemplates_Two(t *testing.T) {

	d := Deployment{
		CustomTemplate: []string{
			"https://w.com/repo1",
			"https://w.com/repo2",
		},
	}

	want := d.CustomTemplate[0] + ", " + d.CustomTemplate[1]
	got := d.FormatCustomTemplates()
	if got != want {
		t.Errorf("FormatCustomTemplates want: %s, got %s", want, got)
		t.Fail()
	}
}

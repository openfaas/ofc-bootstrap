package validators

import (
	"fmt"
	"testing"
)

func Test_ValidateRegistryAuthNoCredStore(t *testing.T) {
	file := []byte(fmt.Sprintf("{ \"auths\": { \"%s\": {\"auth\": \"%s\"} } } ",
		"https://index.docker.io/v1/",
		"Zm9vCg=="))

	got := ValidateRegistryAuth("docker.io/some-user", file)
	if got != nil {
		t.Errorf("error want: %s, got %s", "nil", got)
		t.Fail()
	}
}

func Test_ValidateRegistryAuthCredStoreNoAuth(t *testing.T) {
	file := []byte(fmt.Sprintf("{ \"auths\": { \"%s\": {} }, \"credsStore\": \"something\" } ",
		"https://index.docker.io/v1/"))
	got := ValidateRegistryAuth("docker.io/some-user", file)
	if got == nil {
		t.Errorf("error was nil.")
		t.Fail()
	}
}

func Test_ValidateRegistryAuthNoCredStoreNoAuth(t *testing.T) {
	file := []byte(fmt.Sprintf("{ \"auths\": { \"%s\": {} } } ",
		"index.docker.io/index.html"))
	got := ValidateRegistryAuth("docker.io/some-user", file)
	if got == nil {
		t.Errorf("error was nil.")
		t.Fail()
	}
}

func Test_ValidateRegistryAuthNoValidEndpoint(t *testing.T) {
	file := []byte(fmt.Sprintf("{ \"auths\": { \"%s\": {} } } ",
		""))
	got := ValidateRegistryAuth("docker.io/some-user", file)
	if got == nil {
		t.Errorf("error was nil.")
		t.Fail()
	}
}

func Test_Validate_CorrectDetails(t *testing.T) {
	data := createDockerConfigJson("https://index.docker.io/v1/", "Zm9vOmJhcgo=")

	got := validate(data, "docker.io")
	if got != nil {
		t.Errorf("error want no error, got %q", got)
		t.Fail()
	}
}

func Test_Validate_CorrectDetailsNotDockerRegistry(t *testing.T) {
	data := createDockerConfigJson("https://my.other.registry/auth/v22", "Zm9vOmJhcgo=")

	got := validate(data, "https://my.other.registry/auth/v22")
	if got != nil {
		t.Errorf("error want no error, got %q", got)
		t.Fail()
	}
}

func Test_Validate_NoEntryForDomain(t *testing.T) {
	data := createDockerConfigJson("http://notdocker.io/some-user", "dsonosc")

	got := validate(data, "docker.io")
	if got == nil {
		t.Errorf("error wanted error, got %q", got)
		t.Fail()
	}
}

func Test_Validate_EntryDoesNotContainAuth(t *testing.T) {
	data := createDockerConfigJson("https://index.docker.io/v1/", "Zm9vOmJhcgo=")

	got := validate(data, "notdocker.io")
	if got == nil {
		t.Errorf("error want no error, got %q", got)
		t.Fail()
	}
}

func Test_Validate_EntryDoesNotContainAuthNotDocker(t *testing.T) {
	data := createDockerConfigJson("https://index.myregistry.io/v1/", "Zm9vOmJhcgo=")

	got := validate(data, "https://index.not.my.registry.io/v1/")
	if got == nil {
		t.Errorf("error want no error, got %q", got)
		t.Fail()
	}
}

func Test_Validate_EntryEmptyAuthString(t *testing.T) {
	data := createDockerConfigJson("https://index.docker.io/v1/", "Zm9vOmJhcgo=")

	got := validate(data, "docker.io")
	if got != nil {
		t.Errorf("error want no error, got %q", got)
		t.Fail()
	}
}

func Test_Validate_NonBase64AuthString(t *testing.T) {
	data := createDockerConfigJson("https://index.docker.io/v1/", "ds:\\/onosc")

	got := validate(data, "docker.io")
	if got == nil {
		t.Errorf("want error, got %q", got)
		t.Fail()
	}
}

func createDockerConfigJson(endpoint string, authString string) *DockerConfigJson {

	conf := DockerConfigJson{
		AuthConfigs: map[string]AuthConfig{
			endpoint: {
				Auth: authString,
			},
		},
	}

	return &conf
}

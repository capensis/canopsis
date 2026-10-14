package user

import (
	"slices"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateCreateRequest(sl validator.StructLevel)
}

type baseValidator struct {
	availableAuthSources      map[string]bool
	availableAuthSourcesNames []string
}

func NewValidator(secConfig security.Config) Validator {
	availableAuthSources := make(map[string]bool)
	availableAuthSourcesNames := make([]string, 0, len(secConfig.Security.AuthProviders))

	for _, source := range secConfig.Security.AuthProviders {
		switch source {
		case security.SourceOauth2:
			for oauth2source := range secConfig.Security.OAuth2.Providers {
				availableAuthSources[oauth2source] = true
				availableAuthSourcesNames = append(availableAuthSourcesNames, oauth2source)
			}
		case security.SourceSaml, security.SourceLdap, security.SourceCas:
			availableAuthSources[source] = true
			availableAuthSourcesNames = append(availableAuthSourcesNames, source)
		}
	}

	// necessary for tests to check error.
	slices.Sort(availableAuthSourcesNames)

	return &baseValidator{
		availableAuthSources:      availableAuthSources,
		availableAuthSourcesNames: availableAuthSourcesNames,
	}
}

func (v *baseValidator) ValidateCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	if r.Source != "" {
		if !v.availableAuthSources[r.Source] {
			sl.ReportError(r.Source, "Source", "Source", "oneoforempty", strings.Join(v.availableAuthSourcesNames, ", "))
		}
	}
}

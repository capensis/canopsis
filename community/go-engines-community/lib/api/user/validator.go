package user

import (
	"slices"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateCreateRequest(sl validator.StructLevel)
	ValidateUpdateRequest(sl validator.StructLevel)
	ValidatePatchRequest(sl validator.StructLevel)
	ValidateBulkUpdateRequestItem(sl validator.StructLevel)
	ValidateBulkPatchRequestItem(sl validator.StructLevel)
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

func (v *baseValidator) ValidateBulkUpdateRequestItem(sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkUpdateRequestItem)

	v.validatePassword(sl, r.EditRequest, r.ID)
}

func (v *baseValidator) ValidateBulkPatchRequestItem(sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkPatchRequestItem)

	v.validatePatchEditRequest(sl, r.ID, r.PatchEditRequest)
}

func (v *baseValidator) ValidateCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)

	// Validate source and external_id
	if r.Source == "" && r.ExternalID != "" {
		sl.ReportError(r.Source, "Source", "Source", "required_with", "ExternalID")
	}

	if r.ExternalID == "" && r.Source != "" {
		sl.ReportError(r.ExternalID, "ExternalID", "ExternalID", "required_with", "Source")
	}

	// Validate password
	if r.Source != "" {
		if r.Password != "" {
			sl.ReportError(r.Source, "Source", "Source", "required_not_both", "Password")
		}

		if !v.availableAuthSources[r.Source] {
			sl.ReportError(r.Source, "Source", "Source", "oneoforempty", strings.Join(v.availableAuthSourcesNames, ", "))
		}
	} else {
		v.validatePassword(sl, r.EditRequest, "")
	}
}

func (v *baseValidator) ValidateUpdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)

	v.validatePassword(sl, r.EditRequest, r.ID)
}

func (v *baseValidator) ValidatePatchRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(PatchRequest)

	v.validatePatchEditRequest(sl, r.ID, r.PatchEditRequest)
}

func (v *baseValidator) validatePatchEditRequest(sl validator.StructLevel, id string, r PatchEditRequest) {
	// to avoid code duplication, use validation functions for EditRequest
	editReq := EditRequest{}
	if r.Password != nil {
		editReq.Password = *r.Password
	}

	v.validatePassword(sl, editReq, id)
}

func (v *baseValidator) validatePassword(sl validator.StructLevel, r EditRequest, id string) {
	if r.Password == "" {
		if id == "" {
			sl.ReportError(r.Password, "Password", "Password", "required", "")
		}
	} else {
		if len(r.Password) < password.MinLength {
			sl.ReportError(r.Password, "Password", "Password", "min", strconv.Itoa(password.MinLength))
		}
		if len(r.Password) > password.MaxLength {
			sl.ReportError(r.Password, "Password", "Password", "max", strconv.Itoa(password.MaxLength))
		}
	}
}

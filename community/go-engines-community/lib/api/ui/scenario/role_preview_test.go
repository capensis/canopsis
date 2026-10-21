package scenario

import (
	"context"
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type previewRoleReader struct {
	rolesByUser map[string][]string
	users       []string
	err         error
}

func (r *previewRoleReader) GetRoleNames(_ context.Context, userID string) ([]string, error) {
	r.users = append(r.users, userID)
	if r.err != nil {
		return nil, r.err
	}
	return r.rolesByUser[userID], nil
}

func TestPrepareTemplateAdditionalDataUsesRequestLocalUserCache(t *testing.T) {
	reader := &previewRoleReader{rolesByUser: map[string][]string{
		"event-user-a": {"Operator"},
		"event-user-b": {"Support"},
	}}
	s := &store{userRoleReader: reader}
	cache := make(map[string][]string)

	userA, err := s.prepareTemplateAdditionalData(context.Background(), types.AdditionalData{User: "event-user-a"}, cache)
	if err != nil {
		t.Fatal(err)
	}
	userB, err := s.prepareTemplateAdditionalData(context.Background(), types.AdditionalData{User: "event-user-b"}, cache)
	if err != nil {
		t.Fatal(err)
	}
	userA.UserRoles[0] = "changed"
	reader.rolesByUser["event-user-a"][0] = "Renamed"
	cachedUserA, err := s.prepareTemplateAdditionalData(context.Background(), types.AdditionalData{User: "event-user-a"}, cache)
	if err != nil {
		t.Fatal(err)
	}

	if userB.UserRoles[0] != "Support" || cachedUserA.UserRoles[0] != "Operator" {
		t.Fatalf("role data leaked between actions or was not detached: userB=%#v cachedUserA=%#v", userB.UserRoles, cachedUserA.UserRoles)
	}
	if len(reader.users) != 2 || reader.users[0] != "event-user-a" || reader.users[1] != "event-user-b" {
		t.Fatalf("expected one read per represented user, got %#v", reader.users)
	}
}

func TestPrepareTemplateAdditionalDataClearsRolesForEmptyUser(t *testing.T) {
	reader := &previewRoleReader{}
	s := &store{userRoleReader: reader}
	got, err := s.prepareTemplateAdditionalData(context.Background(), types.AdditionalData{}, make(map[string][]string))
	if err != nil {
		t.Fatal(err)
	}
	if got.UserRoles == nil || len(got.UserRoles) != 0 || len(reader.users) != 0 {
		t.Fatalf("expected no lookup and a non-nil empty slice, data=%#v users=%#v", got, reader.users)
	}
}

func TestPrepareTemplateAdditionalDataPropagatesLookupError(t *testing.T) {
	wantErr := errors.New("role lookup failed")
	reader := &previewRoleReader{err: wantErr}
	s := &store{userRoleReader: reader}
	_, err := s.prepareTemplateAdditionalData(context.Background(), types.AdditionalData{User: "event-user"}, make(map[string][]string))
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped role lookup error, got %v", err)
	}
}

type invalidAuthorTemplateValidator struct{}

func (invalidAuthorTemplateValidator) Validate(string, any) (bool, *validator.ErrReport, string, error) {
	return false, nil, "", nil
}

func TestInvalidNonForwardAuthorClearsPreviousForwardedUser(t *testing.T) {
	forwardAuthor := true
	reader := &previewRoleReader{rolesByUser: map[string][]string{"event-user": {"Operator"}}}
	s := &store{userRoleReader: reader, tplValidator: invalidAuthorTemplateValidator{}}

	_, additionalData, err := s.validateAuthorTpl(
		TemplateActionRequest{Parameters: TemplateActionParameters{ForwardAuthor: &forwardAuthor}},
		types.Event{Author: "event author", UserID: "event-user"},
		types.AdditionalData{},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if additionalData.User != "event-user" {
		t.Fatalf("expected forwarded event user, got %q", additionalData.User)
	}

	_, additionalData, err = s.validateAuthorTpl(
		TemplateActionRequest{Parameters: TemplateActionParameters{Author: "{{ invalid author }}"}},
		types.Event{},
		additionalData,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if additionalData.User != "" {
		t.Fatalf("invalid non-forward author retained previous user %q", additionalData.User)
	}

	prepared, err := s.prepareTemplateAdditionalData(context.Background(), additionalData, make(map[string][]string))
	if err != nil {
		t.Fatal(err)
	}
	if prepared.UserRoles == nil || len(prepared.UserRoles) != 0 || len(reader.users) != 0 {
		t.Fatalf("stale roles were resolved after invalid author: data=%#v lookups=%#v", prepared.UserRoles, reader.users)
	}
}

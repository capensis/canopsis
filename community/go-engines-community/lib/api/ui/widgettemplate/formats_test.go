package widgettemplate

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/binding"
	uiresponse "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	uivalidation "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/gin-gonic/gin"
	ginbinding "github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/rs/zerolog"
)

func TestTicketFormatsBindResponses(t *testing.T) {
	requestValidator := validator.New()
	requestValidator.SetTagName("binding")
	if err := requestValidator.RegisterValidation("notblank", validators.NotBlank); err != nil {
		t.Fatalf("RegisterValidation(notblank) error = %v", err)
	}
	oldBindingValidator := ginbinding.Validator
	ginbinding.Validator = ticketFormatsTestBindingValidator{validate: requestValidator}
	defer func() { ginbinding.Validator = oldBindingValidator }()

	trans, err := libvalidator.NewTranslator(requestValidator, uivalidation.InvalidIDChars)
	if err != nil {
		t.Fatalf("NewTranslator() error = %v", err)
	}
	responder := uiresponse.NewErrorResponder(validation.NewErrorTranslator(trans, zerolog.Nop()), zerolog.Nop())

	tooManyFormats, err := json.Marshal(map[string]any{
		"title":   "Ticket formats",
		"type":    view.WidgetTemplateTypeTicketFormats,
		"formats": make([]string, view.TicketFormatsMaxCount+1),
	})
	if err != nil {
		t.Fatal(err)
	}
	tooLongFormat, err := json.Marshal(map[string]any{
		"title":   "Ticket formats",
		"type":    view.WidgetTemplateTypeTicketFormats,
		"formats": []string{strings.Repeat("é", view.TicketFormatMaxLength+1)},
	})
	if err != nil {
		t.Fatal(err)
	}
	validBoundaryPayloads := make([][]byte, 2)
	for i, formats := range [][]string{
		make([]string, view.TicketFormatsMaxCount),
		{strings.Repeat("é", view.TicketFormatMaxLength)},
	} {
		payload, err := json.Marshal(map[string]any{
			"title":   "Ticket formats",
			"type":    view.WidgetTemplateTypeTicketFormats,
			"formats": formats,
		})
		if err != nil {
			t.Fatal(err)
		}
		validBoundaryPayloads[i] = payload
	}
	for i, payload := range validBoundaryPayloads {
		t.Run("valid boundary "+strconv.Itoa(i), func(t *testing.T) {
			if status, _, err := bindTicketFormatsForTest(t, responder, string(payload), "en"); status != 200 || err != nil {
				t.Fatalf("Bind() status = %d, want 200", status)
			}
		})
	}

	tests := []struct {
		name      string
		payload   string
		field     string
		wantTag   string
		wantParam string
	}{
		{name: "null list", payload: `{"title":"Ticket formats","type":"ticket_formats","formats":null}`, field: "formats"},
		{name: "non-array list", payload: `{"title":"Ticket formats","type":"ticket_formats","formats":"INC"}`, field: "formats"},
		{name: "null item", payload: `{"title":"Ticket formats","type":"ticket_formats","formats":[null]}`, field: "formats.0"},
		{name: "number item", payload: `{"title":"Ticket formats","type":"ticket_formats","formats":[1]}`, field: "formats.0"},
		{name: "object item", payload: `{"title":"Ticket formats","type":"ticket_formats","formats":[{}]}`, field: "formats.0"},
		{name: "later duplicate after malformed item", payload: `{"title":"Ticket formats","type":"ticket_formats","formats":["valid",null],"Formats":[]}`, field: "formats.1"},
		{name: "escaped case variant key", payload: `{"title":"Ticket formats","type":"ticket_formats","format\u0073":[null]}`, field: "formats.0"},
		{name: "too many entries", payload: string(tooManyFormats), field: "formats", wantTag: "max", wantParam: "100"},
		{name: "too many unicode characters", payload: string(tooLongFormat), field: "formats.0", wantTag: "max", wantParam: "255"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantTag != "" {
				var request EditRequest
				if err := json.Unmarshal([]byte(tc.payload), &request); err != nil {
					t.Fatalf("json.Unmarshal() error = %v", err)
				}
				validationErr := requestValidator.Struct(request)
				var validationErrs validator.ValidationErrors
				if !errors.As(validationErr, &validationErrs) || len(validationErrs) != 1 {
					t.Fatalf("validator error = %v, want one validation error", validationErr)
				}
				fieldErr := validationErrs[0]
				if fieldErr.Tag() != tc.wantTag || fieldErr.Param() != tc.wantParam {
					t.Errorf("tag/param = %q/%q, want %q/%q", fieldErr.Tag(), fieldErr.Param(), tc.wantTag, tc.wantParam)
				}
				gotField := validation.NewError(validationErrs, request).TransformNamespace(fieldErr.StructNamespace())
				if gotField != tc.field {
					t.Errorf("validation field = %q, want %q", gotField, tc.field)
				}
			}
			for _, locale := range []string{"en", "fr"} {
				status, responseErrors, bindErr := bindTicketFormatsForTest(t, responder, tc.payload, locale)
				if status != 400 {
					t.Fatalf("HTTP status = %d, want 400", status)
				}
				var validationErr *validation.Error
				if !errors.As(bindErr, &validationErr) {
					t.Fatalf("Bind() error = %T %v, want translated validation error", bindErr, bindErr)
				}
				if len(responseErrors) != 1 {
					t.Fatalf("HTTP errors = %#v, want one error", responseErrors)
				}
				message, ok := responseErrors[tc.field]
				if !ok || strings.TrimSpace(message) == "" {
					t.Fatalf("HTTP errors = %#v, want non-empty message for %q", responseErrors, tc.field)
				}
				if raw := strings.TrimSpace(validationErr.Error()); raw != "" && strings.Contains(message, raw) {
					t.Errorf("HTTP error contains untranslated validator message %q", raw)
				}
				for _, internal := range []string{"EditRequest", "Field validation", "string_array", "on the '"} {
					if strings.Contains(message, internal) {
						t.Errorf("HTTP error leaked validator internals: %q", message)
					}
				}
			}
		})
	}
}

type ticketFormatsTestBindingValidator struct {
	validate *validator.Validate
}

func (v ticketFormatsTestBindingValidator) ValidateStruct(value any) error {
	return v.validate.Struct(value)
}

func (v ticketFormatsTestBindingValidator) Engine() any {
	return v.validate
}

func bindTicketFormatsForTest(t *testing.T, responder uiresponse.ErrorResponder, payload, locale string) (int, map[string]string, error) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequestWithContext(t.Context(), "POST", "/internal/widget-templates", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", gin.MIMEJSON)
	authctx.SetLocale(c, locale)

	var request EditRequest
	err := binding.Bind(c, &request)
	if err == nil {
		return http.StatusOK, nil, nil
	}
	responder.Respond(c, err)
	var body struct {
		Errors map[string]string `json:"errors"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal(response) error = %v; body: %s", err, w.Body.String())
	}
	return w.Code, body.Errors, err
}

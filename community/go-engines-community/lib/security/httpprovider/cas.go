package httpprovider

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	libhttp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

// casProvider implements CAS authentication.
type casProvider struct {
	client       libhttp.Doer
	config       security.CasConfig
	userProvider security.UserProvider
}

// NewCasProvider creates new provider.
func NewCasProvider(
	client libhttp.Doer,
	config security.CasConfig,
	userProvider security.UserProvider,
) security.HttpProvider {
	return &casProvider{
		client:       client,
		config:       config,
		userProvider: userProvider,
	}
}

type casResponse struct { //nolint:musttag
	XMLName xml.Name `xml:"http://www.yale.edu/tp/cas serviceResponse"`

	Failure struct {
		XMLName xml.Name `xml:"authenticationFailure"`
		Code    string   `xml:"code,attr"`
		Message string   `xml:",innerxml"`
	}
	Success struct {
		XMLName xml.Name `xml:"authenticationSuccess"`
		User    string   `xml:"user"`
	}
}

func (p *casProvider) Auth(request *http.Request) (*security.User, error, bool) {
	ticket := request.URL.Query().Get(security.QueryParamCasTicket)
	if ticket == "" {
		return nil, nil, false
	}

	service := request.URL.Query().Get(security.QueryParamCasService)
	if service == "" {
		return nil, nil, false
	}

	// Add request query (except ticket) to service
	serviceUrl, err := url.Parse(service)
	if err != nil {
		return nil, nil, false
	}
	serviceQuery := serviceUrl.Query()
	query := request.URL.Query()
	query.Del(security.QueryParamCasTicket)
	for k, vals := range query {
		for _, v := range vals {
			serviceQuery.Add(k, v)
		}
	}
	serviceUrl.RawQuery = serviceQuery.Encode()
	service = serviceUrl.String()

	username, err := p.validateTicket(request.Context(), ticket, service)
	if err != nil || username == "" {
		return nil, err, true
	}

	user, err := p.saveUser(request.Context(), username)
	return user, err, true
}

// validateTicket calls CAS server to validate ticket.
func (p *casProvider) validateTicket(
	ctx context.Context,
	ticket, service string,
) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.config.ValidateUrl, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Set("ticket", ticket)
	q.Set("service", service)
	req.URL.RawQuery = q.Encode()
	res, err := p.client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", err
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var casRes casResponse
	err = xml.Unmarshal(buf, &casRes)
	if err != nil {
		return "", err
	}

	if casRes.Failure.Code != "" {
		return "", nil
	}

	return casRes.Success.User, nil
}

// saveUser adds user data to storage.
func (p *casProvider) saveUser(ctx context.Context, username string) (*security.User, error) {
	user, err := p.userProvider.FindByExternalSource(ctx, username, security.SourceCas)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %w", err)
	}

	if user == nil {
		user = &security.User{
			Name:       username,
			Roles:      []string{p.config.DefaultRole},
			IsEnabled:  true,
			ExternalID: username,
			Source:     security.SourceCas,
		}
		err = p.userProvider.Save(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("cannot save user: %w", err)
		}
	} else if !user.IsEnabled {
		return nil, nil
	}

	return user, nil
}

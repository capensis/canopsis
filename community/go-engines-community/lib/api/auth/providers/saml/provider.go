package saml

import (
	"bytes"
	"cmp"
	"compress/flate"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/roleprovider"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"github.com/beevik/etree"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	saml2 "github.com/russellhaering/gosaml2"
	samltypes "github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

const MetadataReqTimeout = time.Second * 15

const (
	BindingRedirect = "redirect"
	BindingPOST     = "post"
)

type Provider interface {
	SamlMetadataHandler() gin.HandlerFunc
	SamlAuthHandler() gin.HandlerFunc
	SamlAcsHandler() gin.HandlerFunc
	SamlSloHandler() gin.HandlerFunc
}

type provider struct {
	samlSP             *saml2.SAMLServiceProvider
	userProvider       security.UserProvider
	roleProvider       security.RoleProvider
	sessionStore       libsession.Store
	enforcer           security.Enforcer
	config             security.SamlConfig
	tokenService       apisecurity.TokenService
	maintenanceAdapter config.MaintenanceAdapter
	logger             zerolog.Logger
}

func NewProvider(
	ctx context.Context,
	userProvider security.UserProvider,
	roleValidator security.RoleProvider,
	sessionStore libsession.Store,
	enforcer security.Enforcer,
	config security.Config,
	tokenService apisecurity.TokenService,
	maintenanceAdapter config.MaintenanceAdapter,
	logger zerolog.Logger,
) (Provider, error) {
	if config.Security.Saml.IdpMetadataUrl != "" && config.Security.Saml.IdpMetadataXml != "" {
		return nil, errors.New("should provide only idp metadata url or xml, not both")
	}

	if config.Security.Saml.CanopsisSSOBinding != BindingRedirect && config.Security.Saml.CanopsisSSOBinding != BindingPOST {
		return nil, errors.New("wrong canopsis_sso_binding value, should be post or redirect")
	}

	if config.Security.Saml.CanopsisACSBinding != BindingRedirect && config.Security.Saml.CanopsisACSBinding != BindingPOST {
		return nil, errors.New("wrong canopsis_acs_binding value, should be post or redirect")
	}

	keyPair, err := tls.LoadX509KeyPair(config.Security.Saml.X509Cert, config.Security.Saml.X509Key)
	if err != nil {
		return nil, err
	}

	idpMetadata := &samltypes.EntityDescriptor{}
	if config.Security.Saml.IdpMetadataUrl != "" {
		dt, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			return nil, errors.New("unknown type of http.DefaultTransport")
		}

		tr := dt.Clone()
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: config.Security.Saml.InsecureSkipVerify} //nolint:gosec

		hc := &http.Client{Timeout: MetadataReqTimeout, Transport: tr}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.Security.Saml.IdpMetadataUrl, nil)
		if err != nil {
			return nil, err
		}

		res, err := hc.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		rawMetadata, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = xml.Unmarshal(rawMetadata, idpMetadata)
		if err != nil {
			return nil, err
		}
	}

	if config.Security.Saml.IdpMetadataXml != "" {
		rawMetadata, err := os.ReadFile(config.Security.Saml.IdpMetadataXml)
		if err != nil {
			return nil, err
		}

		err = xml.Unmarshal(rawMetadata, idpMetadata)
		if err != nil {
			return nil, err
		}
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range idpMetadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				panic(fmt.Errorf("metadata certificate(%d) must not be empty", idx))
			}
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				return nil, err
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				return nil, err
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	ssoLocation := ""
	sloLocation := ""

	if len(idpMetadata.IDPSSODescriptor.SingleSignOnServices) > 0 {
		ssoLocation = idpMetadata.IDPSSODescriptor.SingleSignOnServices[0].Location
	}

	if len(idpMetadata.IDPSSODescriptor.SingleLogoutServices) > 0 {
		sloLocation = idpMetadata.IDPSSODescriptor.SingleLogoutServices[0].Location
	}

	return &provider{
		samlSP: &saml2.SAMLServiceProvider{
			IdentityProviderSSOURL:         ssoLocation,
			IdentityProviderSLOURL:         sloLocation,
			IdentityProviderIssuer:         idpMetadata.EntityID,
			AssertionConsumerServiceURL:    fmt.Sprintf("%s/%s", config.Security.Saml.CanopsisSamlUrl, "acs"),
			ServiceProviderSLOURL:          fmt.Sprintf("%s/%s", config.Security.Saml.CanopsisSamlUrl, "slo"),
			ServiceProviderIssuer:          fmt.Sprintf("%s/%s", config.Security.Saml.CanopsisSamlUrl, "metadata"),
			SignAuthnRequests:              config.Security.Saml.SignAuthRequest,
			AudienceURI:                    fmt.Sprintf("%s/%s", config.Security.Saml.CanopsisSamlUrl, "metadata"),
			IDPCertificateStore:            &certStore,
			SPKeyStore:                     dsig.TLSCertKeyStore(keyPair),
			NameIdFormat:                   config.Security.Saml.NameIdFormat,
			SkipSignatureValidation:        config.Security.Saml.SkipSignatureValidation,
			SignAuthnRequestsCanonicalizer: dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList(""),
		},
		userProvider:       userProvider,
		roleProvider:       roleValidator,
		sessionStore:       sessionStore,
		enforcer:           enforcer,
		config:             config.Security.Saml,
		tokenService:       tokenService,
		maintenanceAdapter: maintenanceAdapter,
		logger:             logger,
	}, nil
}

func (p *provider) SamlMetadataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		meta, err := p.samlSP.MetadataWithSLO(0)
		if err != nil {
			p.logger.Err(err).Msg("samlMetadataHandler: MetadataWithSLO error")
			panic(err)
		}

		if len(meta.SPSSODescriptor.SingleLogoutServices) > 0 && p.config.CanopsisSSOBinding == BindingRedirect {
			meta.SPSSODescriptor.SingleLogoutServices[0].Binding = saml2.BindingHttpRedirect
		}

		if len(meta.SPSSODescriptor.AssertionConsumerServices) > 0 {
			if p.config.CanopsisACSBinding == BindingRedirect {
				meta.SPSSODescriptor.AssertionConsumerServices[0].Binding = saml2.BindingHttpRedirect
			}

			if p.config.ACSIndex != nil {
				meta.SPSSODescriptor.AssertionConsumerServices[0].Index = *p.config.ACSIndex
			} else {
				meta.SPSSODescriptor.AssertionConsumerServices[0].Index = 1
			}
		}

		c.XML(http.StatusOK, meta)
	}
}

type samlLoginRequest struct {
	RelayState string `form:"relayState" binding:"required,url"`
}

func (p *provider) SamlAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if p.samlSP.IdentityProviderSSOURL == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		request := samlLoginRequest{}

		if err := c.ShouldBindQuery(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		if p.config.CanopsisSSOBinding == BindingRedirect {
			var authRequest *etree.Document
			var err error

			if p.config.SignAuthRequest {
				authRequest, err = p.samlSP.BuildAuthRequestDocument()
			} else {
				authRequest, err = p.samlSP.BuildAuthRequestDocumentNoSig()
			}

			if err != nil {
				p.logger.Err(err).Msg("samlAuthHandler: BuildAuthRequestDocument error")
				panic(err)
			}

			el := authRequest.SelectElement("AuthnRequest")
			attr := el.SelectAttr("ProtocolBinding")
			attr.Value = saml2.BindingHttpRedirect

			authUrl, err := p.samlSP.BuildAuthURLRedirect(request.RelayState, authRequest)
			if err != nil {
				p.logger.Err(err).Msg("samlAuthHandler: parse IdentityProviderSSOURL error")
				panic(err)
			}

			c.Redirect(http.StatusPermanentRedirect, authUrl)
		}

		if p.config.CanopsisSSOBinding == BindingPOST {
			b, err := p.samlSP.BuildAuthBodyPost(request.RelayState)
			if err != nil {
				p.logger.Err(err).Msg("samlAuthHandler: BuildAuthBodyPost error")
				panic(err)
			}

			c.Data(http.StatusOK, gin.MIMEHTML, b)
		}
	}
}

func (p *provider) SamlAcsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		samlResponse, exists := c.GetPostForm("SAMLResponse")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(errors.New("SAMLResponse doesn't exist")))
			return
		}

		relayState, exists := c.GetPostForm("RelayState")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(errors.New("RelayState doesn't exist")))
			return
		}

		relayUrl, err := url.Parse(relayState)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(errors.New("RelayState is not a valid url")))
			return
		}

		assertionInfo, err := p.samlSP.RetrieveAssertionInfo(samlResponse)
		if err != nil {
			p.logger.Err(err).Msg("assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.InvalidTime {
			p.logger.Err(errors.New("invalid time")).Msg("assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.NotInAudience {
			p.logger.Err(errors.New("not in audience")).Msg("assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		user, err := p.userProvider.FindByExternalSource(c, assertionInfo.NameID, security.SourceSaml)
		if err != nil {
			p.logger.Err(err).Msg("samlAcsHandler: userProvider FindByExternalSource error")
			panic(err)
		}

		if user == nil {
			var ok bool
			user, ok = p.createUser(c, relayUrl, assertionInfo)
			if !ok {
				return
			}
		}

		maintenanceConf, err := p.maintenanceAdapter.GetConfig(c)
		if err != nil {
			panic(err)
		}

		if maintenanceConf.Enabled {
			ok, err := p.enforcer.Enforce(user.ID, apisecurity.PermMaintenance, model.PermissionCan)
			if err != nil {
				panic(err)
			}

			if !ok {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, common.CanopsisUnderMaintenanceResponse)
				return
			}
		}

		err = p.enforcer.LoadPolicy()
		if err != nil {
			panic(fmt.Errorf("reload enforcer error: %w", err))
		}

		var accessToken string
		if p.config.ExpirationInterval != "" || assertionInfo.SessionNotOnOrAfter == nil {
			accessToken, err = p.tokenService.Create(c, *user, security.AuthMethodSaml)
		} else {
			accessToken, err = p.tokenService.CreateWithExpiration(c, *user, security.AuthMethodSaml, *assertionInfo.SessionNotOnOrAfter)
		}
		if err != nil {
			panic(err)
		}

		query := relayUrl.Query()
		query.Set("access_token", accessToken)
		relayUrl.RawQuery = query.Encode()

		c.Redirect(http.StatusPermanentRedirect, relayUrl.String())
	}
}

func (p *provider) getAssocAttribute(attrs saml2.Values, canopsisName, defaultValue string) string {
	return cmp.Or(attrs.Get(p.config.IdpAttributesMap[canopsisName]), defaultValue)
}

func (p *provider) getAssocArrayAttribute(attrs saml2.Values, canopsisName string, defaultValue []string) []string {
	v := attrs.GetAll(p.config.IdpAttributesMap[canopsisName])
	if len(v) != 0 {
		return v
	}

	return defaultValue
}

func (p *provider) SamlSloHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if p.samlSP.IdentityProviderSLOURL == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		samlRequest, exists := c.GetQuery("SAMLRequest")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(errors.New("SAMLRequest doesn't exist")))
			return
		}

		relayState, exists := c.GetQuery("RelayState")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(errors.New("RelayState doesn't exist")))
			return
		}

		request, err := p.samlSP.ValidateEncodedLogoutRequestPOST(samlRequest)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		user, err := p.userProvider.FindByExternalSource(c, request.NameID.Value, security.SourceSaml)
		if err != nil {
			p.logger.Err(err).Msg("samlSloHandler: userProvider FindByExternalSource error")
			panic(err)
		}

		if user == nil {
			responseUrl, err := p.buildLogoutResponseUrl(saml2.StatusCodeUnknownPrincipal, request.ID, relayState)
			if err != nil {
				p.logger.Err(err).Msg("samlSloHandler: buildLogoutResponseUrl error")
				panic(err)
			}

			c.Redirect(http.StatusPermanentRedirect, responseUrl.String())
			return
		}

		err = p.sessionStore.ExpireSessions(c, user.ID, security.AuthMethodSaml)
		if err != nil {
			responseUrl, err := p.buildLogoutResponseUrl(saml2.StatusCodeUnknownPrincipal, request.ID, relayState)
			if err != nil {
				p.logger.Err(err).Msg("samlSloHandler: buildLogoutResponseUrl error")
				panic(err)
			}

			c.Redirect(http.StatusPermanentRedirect, responseUrl.String())
			return
		}

		err = p.tokenService.DeleteBy(c, user.ID, security.AuthMethodSaml)
		if err != nil {
			responseUrl, err := p.buildLogoutResponseUrl(saml2.StatusCodeUnknownPrincipal, request.ID, relayState)
			if err != nil {
				p.logger.Err(err).Msg("samlSloHandler: buildLogoutResponseUrl error")
				panic(err)
			}

			c.Redirect(http.StatusPermanentRedirect, responseUrl.String())
			return
		}

		responseUrl, err := p.buildLogoutResponseUrl(saml2.StatusCodeSuccess, request.ID, relayState)
		if err != nil {
			p.logger.Err(err).Msg("samlSloHandler: buildLogoutResponseUrl error")
			panic(err)
		}

		c.Redirect(http.StatusPermanentRedirect, responseUrl.String())
	}
}

func (p *provider) buildLogoutResponseUrl(status, reqID, relayState string) (*url.URL, error) {
	responseDoc, err := p.samlSP.BuildLogoutResponseDocument(status, reqID)
	if err != nil {
		return nil, err
	}

	buffer, err := p.encodeAndCompress(responseDoc)
	if err != nil {
		return nil, err
	}

	responseUrl, err := url.Parse(p.samlSP.IdentityProviderSLOURL)
	if err != nil {
		return nil, err
	}

	query := responseUrl.Query()
	query.Set("SAMLResponse", buffer.String())
	query.Set("RelayState", relayState)

	responseUrl.RawQuery = query.Encode()

	return responseUrl, nil
}

func (p *provider) encodeAndCompress(doc io.WriterTo) (_ *bytes.Buffer, resErr error) {
	buffer := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)

	defer func() {
		err := encoder.Close()
		if err != nil && resErr == nil {
			resErr = err
		}
	}()

	compressor, err := flate.NewWriter(encoder, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = compressor.Close()
		if err != nil && resErr == nil {
			resErr = err
		}
	}()

	_, err = doc.WriteTo(compressor)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (p *provider) createUser(c *gin.Context, relayUrl *url.URL, assertionInfo *saml2.AssertionInfo) (*security.User, bool) {
	if !p.config.AutoUserRegistration {
		p.logger.Err(fmt.Errorf("user with external_id = %s is not found", assertionInfo.NameID)).Msg("autoUserRegistration is disabled")
		p.errorRedirect(c, relayUrl, "This user is not allowed to log into Canopsis")

		return nil, false
	}

	roles, err := p.roleProvider.GetValidRoleIDs(c, p.getAssocArrayAttribute(assertionInfo.Values, "role", []string{}), p.config.DefaultRole)
	if err != nil {
		roleNotFoundError := roleprovider.ProviderError{}
		if errors.As(err, &roleNotFoundError) {
			p.logger.Err(roleNotFoundError).Msg("user registration failed")
			p.errorRedirect(c, relayUrl, roleNotFoundError.Error())

			return nil, false
		}

		panic(err)
	}

	user := &security.User{
		Name:       p.getAssocAttribute(assertionInfo.Values, "name", assertionInfo.NameID),
		Roles:      roles,
		IsEnabled:  true,
		ExternalID: assertionInfo.NameID,
		Source:     security.SourceSaml,
		Firstname:  p.getAssocAttribute(assertionInfo.Values, "firstname", ""),
		Lastname:   p.getAssocAttribute(assertionInfo.Values, "lastname", ""),
		Email:      p.getAssocAttribute(assertionInfo.Values, "email", ""),
	}
	err = p.userProvider.Save(c, user)
	if err != nil {
		p.logger.Err(err).Msg("samlAcsHandler: userProvider Save error")
		panic(fmt.Errorf("cannot save user: %w", err))
	}

	return user, true
}

func (p *provider) errorRedirect(c *gin.Context, relayUrl *url.URL, errorMessage string) {
	query := relayUrl.Query()
	query.Set("errorMessage", errorMessage)
	relayUrl.RawQuery = query.Encode()

	c.Redirect(http.StatusPermanentRedirect, relayUrl.String())
}

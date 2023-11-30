package saml

import (
	"bytes"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	libsession "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"github.com/beevik/etree"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
	saml2 "github.com/russellhaering/gosaml2"
	samltypes "github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const MetadataReqTimeout = time.Second * 15

const (
	BindingRedirect = "redirect"
	BindingPOST     = "post"

	DefaultUserRole = "Norights"
)

type ServiceProvider interface {
	SamlMetadataHandler() gin.HandlerFunc
	SamlAuthHandler() gin.HandlerFunc
	SamlSessionAcsHandler() gin.HandlerFunc
	SamlAcsHandler() gin.HandlerFunc
	SamlSloHandler() gin.HandlerFunc
}

type serviceProvider struct {
	samlSP             *saml2.SAMLServiceProvider
	userProvider       security.UserProvider
	roleCollection     mongo.DbCollection
	sessionStore       libsession.Store
	enforcer           security.Enforcer
	config             security.Config
	tokenService       apisecurity.TokenService
	maintenanceAdapter config.MaintenanceAdapter
	logger             zerolog.Logger
	defaultRole        string
}

func NewServiceProvider(
	ctx context.Context,
	userProvider security.UserProvider,
	roleCollection mongo.DbCollection,
	sessionStore libsession.Store,
	enforcer security.Enforcer,
	config security.Config,
	tokenService apisecurity.TokenService,
	maintenanceAdapter config.MaintenanceAdapter,
	logger zerolog.Logger,
) (ServiceProvider, error) {
	if config.Security.Saml.IdpMetadataUrl != "" && config.Security.Saml.IdpMetadataXml != "" {
		return nil, fmt.Errorf("should provide only idp metadata url or xml, not both")
	}

	if config.Security.Saml.CanopsisSSOBinding != BindingRedirect && config.Security.Saml.CanopsisSSOBinding != BindingPOST {
		return nil, fmt.Errorf("wrong canopsis_sso_binding value, should be post or redirect")
	}

	if config.Security.Saml.CanopsisACSBinding != BindingRedirect && config.Security.Saml.CanopsisACSBinding != BindingPOST {
		return nil, fmt.Errorf("wrong canopsis_acs_binding value, should be post or redirect")
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

	defaultRole := config.Security.Saml.DefaultRole
	if defaultRole == "" {
		defaultRole = DefaultUserRole
	}

	return &serviceProvider{
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
		roleCollection:     roleCollection,
		sessionStore:       sessionStore,
		enforcer:           enforcer,
		config:             config,
		tokenService:       tokenService,
		maintenanceAdapter: maintenanceAdapter,
		defaultRole:        defaultRole,
		logger:             logger,
	}, nil
}

func (sp *serviceProvider) SamlMetadataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		meta, err := sp.samlSP.MetadataWithSLO(0)
		if err != nil {
			sp.logger.Err(err).Msg("SamlMetadataHandler: MetadataWithSLO error")
			panic(err)
		}

		if len(meta.SPSSODescriptor.SingleLogoutServices) > 0 && sp.config.Security.Saml.CanopsisSSOBinding == BindingRedirect {
			meta.SPSSODescriptor.SingleLogoutServices[0].Binding = saml2.BindingHttpRedirect
		}

		if len(meta.SPSSODescriptor.AssertionConsumerServices) > 0 {
			if sp.config.Security.Saml.CanopsisACSBinding == BindingRedirect {
				meta.SPSSODescriptor.AssertionConsumerServices[0].Binding = saml2.BindingHttpRedirect
			}

			if sp.config.Security.Saml.ACSIndex != nil {
				meta.SPSSODescriptor.AssertionConsumerServices[0].Index = *sp.config.Security.Saml.ACSIndex
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

func (sp *serviceProvider) SamlAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if sp.samlSP.IdentityProviderSSOURL == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		request := samlLoginRequest{}

		if err := c.ShouldBindQuery(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		if sp.config.Security.Saml.CanopsisSSOBinding == BindingRedirect {
			var authRequest *etree.Document
			var err error

			if sp.config.Security.Saml.SignAuthRequest {
				authRequest, err = sp.samlSP.BuildAuthRequestDocument()
			} else {
				authRequest, err = sp.samlSP.BuildAuthRequestDocumentNoSig()
			}

			if err != nil {
				sp.logger.Err(err).Msg("SamlAuthHandler: BuildAuthRequestDocument error")
				panic(err)
			}

			el := authRequest.SelectElement("AuthnRequest")
			attr := el.SelectAttr("ProtocolBinding")
			attr.Value = saml2.BindingHttpRedirect

			authUrl, err := sp.samlSP.BuildAuthURLRedirect(request.RelayState, authRequest)
			if err != nil {
				sp.logger.Err(err).Msg("SamlAuthHandler: parse IdentityProviderSSOURL error")
				panic(err)
			}

			c.Redirect(http.StatusPermanentRedirect, authUrl)
		}

		if sp.config.Security.Saml.CanopsisSSOBinding == BindingPOST {
			b, err := sp.samlSP.BuildAuthBodyPost(request.RelayState)
			if err != nil {
				sp.logger.Err(err).Msg("SamlAuthHandler: BuildAuthBodyPost error")
				panic(err)
			}

			c.Data(http.StatusOK, gin.MIMEHTML, b)
		}
	}
}

func (sp *serviceProvider) SamlSessionAcsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		samlResponse, exists := c.GetPostForm("SAMLResponse")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("SAMLResponse doesn't exist")))
			return
		}

		relayState, exists := c.GetPostForm("RelayState")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("RelayState doesn't exist")))
			return
		}

		relayUrl, err := url.Parse(relayState)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("RelayState is not a valid url")))
			return
		}

		assertionInfo, err := sp.samlSP.RetrieveAssertionInfo(samlResponse)
		if err != nil {
			sp.logger.Err(err).Msg("Assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.InvalidTime {
			sp.logger.Err(fmt.Errorf("invalid time")).Msg("Assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.NotInAudience {
			sp.logger.Err(fmt.Errorf("not in audience")).Msg("Assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		user, err := sp.userProvider.FindByExternalSource(c, assertionInfo.NameID, security.SourceSaml)
		if err != nil {
			sp.logger.Err(err).Msg("SamlAcsHandler: userProvider FindByExternalSource error")
			panic(err)
		}

		if user == nil {
			var ok bool
			user, ok = sp.createUser(c, relayUrl, assertionInfo)
			if !ok {
				return
			}
		}

		session := sp.getSession(c)
		session.Values["user"] = user.ID
		session.Values["provider"] = security.AuthMethodSaml

		if assertionInfo.SessionNotOnOrAfter != nil {
			session.Options.MaxAge = int(time.Until(*assertionInfo.SessionNotOnOrAfter).Seconds())
		}

		err = session.Save(c.Request, c.Writer)
		if err != nil {
			sp.logger.Err(err).Msg("SamlAcsHandler: save session error")
			panic(err)
		}

		err = sp.enforcer.LoadPolicy()
		if err != nil {
			panic(fmt.Errorf("reload enforcer error: %w", err))
		}

		c.Redirect(http.StatusPermanentRedirect, relayUrl.String())
	}
}

func (sp *serviceProvider) SamlAcsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		samlResponse, exists := c.GetPostForm("SAMLResponse")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("SAMLResponse doesn't exist")))
			return
		}

		relayState, exists := c.GetPostForm("RelayState")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("RelayState doesn't exist")))
			return
		}

		relayUrl, err := url.Parse(relayState)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("RelayState is not a valid url")))
			return
		}

		assertionInfo, err := sp.samlSP.RetrieveAssertionInfo(samlResponse)
		if err != nil {
			sp.logger.Err(err).Msg("Assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.InvalidTime {
			sp.logger.Err(fmt.Errorf("invalid time")).Msg("Assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.NotInAudience {
			sp.logger.Err(fmt.Errorf("not in audience")).Msg("Assertion is not valid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		user, err := sp.userProvider.FindByExternalSource(c, assertionInfo.NameID, security.SourceSaml)
		if err != nil {
			sp.logger.Err(err).Msg("SamlAcsHandler: userProvider FindByExternalSource error")
			panic(err)
		}

		if user == nil {
			var ok bool
			user, ok = sp.createUser(c, relayUrl, assertionInfo)
			if !ok {
				return
			}
		}

		maintenanceConf, err := sp.maintenanceAdapter.GetConfig(c)
		if err != nil {
			panic(err)
		}

		if maintenanceConf.Enabled {
			ok, err := sp.enforcer.Enforce(user.ID, apisecurity.PermMaintenance, model.PermissionCan)
			if err != nil {
				panic(err)
			}

			if !ok {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, common.CanopsisUnderMaintenanceResponse)
				return
			}
		}

		err = sp.enforcer.LoadPolicy()
		if err != nil {
			panic(fmt.Errorf("reload enforcer error: %w", err))
		}

		var accessToken string
		if assertionInfo.SessionNotOnOrAfter == nil {
			accessToken, err = sp.tokenService.Create(c, *user, security.AuthMethodSaml)
		} else {
			accessToken, err = sp.tokenService.CreateWithExpiration(c, *user, security.AuthMethodSaml, *assertionInfo.SessionNotOnOrAfter)
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

func (sp *serviceProvider) getAssocAttribute(attrs saml2.Values, canopsisName, defaultValue string) string {
	v := defaultValue

	idpAssoc, ok := sp.config.Security.Saml.IdpAttributesMap[canopsisName]
	if ok {
		v = attrs.Get(idpAssoc)
	}

	return v
}

func (sp *serviceProvider) SamlSloHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if sp.samlSP.IdentityProviderSLOURL == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		samlRequest, exists := c.GetQuery("SAMLRequest")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("SAMLRequest doesn't exist")))
			return
		}

		relayState, exists := c.GetQuery("RelayState")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("RelayState doesn't exist")))
			return
		}

		request, err := sp.samlSP.ValidateEncodedLogoutRequestPOST(samlRequest)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		user, err := sp.userProvider.FindByExternalSource(c, request.NameID.Value, security.SourceSaml)
		if err != nil {
			sp.logger.Err(err).Msg("SamlSloHandler: userProvider FindByExternalSource error")
			panic(err)
		}

		if user == nil {
			responseUrl, err := sp.buildLogoutResponseUrl(saml2.StatusCodeUnknownPrincipal, request.ID, relayState)
			if err != nil {
				sp.logger.Err(err).Msg("SamlSloHandler: buildLogoutResponseUrl error")
				panic(err)
			}

			c.Redirect(http.StatusPermanentRedirect, responseUrl.String())
			return
		}

		err = sp.sessionStore.ExpireSessions(c, user.ID, security.AuthMethodSaml)
		if err != nil {
			responseUrl, err := sp.buildLogoutResponseUrl(saml2.StatusCodeUnknownPrincipal, request.ID, relayState)
			if err != nil {
				sp.logger.Err(err).Msg("SamlSloHandler: buildLogoutResponseUrl error")
				panic(err)
			}

			c.Redirect(http.StatusPermanentRedirect, responseUrl.String())
			return
		}

		err = sp.tokenService.DeleteBy(c, user.ID, security.AuthMethodSaml)
		if err != nil {
			responseUrl, err := sp.buildLogoutResponseUrl(saml2.StatusCodeUnknownPrincipal, request.ID, relayState)
			if err != nil {
				sp.logger.Err(err).Msg("SamlSloHandler: buildLogoutResponseUrl error")
				panic(err)
			}

			c.Redirect(http.StatusPermanentRedirect, responseUrl.String())
			return
		}

		responseUrl, err := sp.buildLogoutResponseUrl(saml2.StatusCodeSuccess, request.ID, relayState)
		if err != nil {
			sp.logger.Err(err).Msg("SamlSloHandler: buildLogoutResponseUrl error")
			panic(err)
		}

		c.Redirect(http.StatusPermanentRedirect, responseUrl.String())
	}
}

func (sp *serviceProvider) buildLogoutResponseUrl(status, reqID, relayState string) (*url.URL, error) {
	responseDoc, err := sp.samlSP.BuildLogoutResponseDocument(status, reqID)
	if err != nil {
		return nil, err
	}

	buffer, err := sp.encodeAndCompress(responseDoc)
	if err != nil {
		return nil, err
	}

	responseUrl, err := url.Parse(sp.samlSP.IdentityProviderSLOURL)
	if err != nil {
		return nil, err
	}

	query := responseUrl.Query()
	query.Set("SAMLResponse", buffer.String())
	query.Set("RelayState", relayState)

	responseUrl.RawQuery = query.Encode()

	return responseUrl, nil
}

func (sp *serviceProvider) encodeAndCompress(doc io.WriterTo) (_ *bytes.Buffer, resErr error) {
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

func (sp *serviceProvider) getSession(c *gin.Context) *sessions.Session {
	session, err := sp.sessionStore.Get(c.Request, security.SessionKey)
	if err != nil {
		panic(err)
	}

	return session
}

func (sp *serviceProvider) createUser(c *gin.Context, relayUrl *url.URL, assertionInfo *saml2.AssertionInfo) (*security.User, bool) {
	if !sp.config.Security.Saml.AutoUserRegistration {
		sp.logger.Err(fmt.Errorf("user with external_id = %s is not found", assertionInfo.NameID)).Msg("AutoUserRegistration is disabled")
		sp.errorRedirect(c, relayUrl, "This user is not allowed to log into Canopsis")

		return nil, false
	}

	role := sp.getAssocAttribute(assertionInfo.Values, "role", sp.defaultRole)

	err := sp.roleCollection.FindOne(c, bson.M{"name": role}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			errMessage := fmt.Errorf("role %s doesn't exist", role)
			sp.logger.Err(errMessage).Msg("User registration failed")
			sp.errorRedirect(c, relayUrl, errMessage.Error())
			return nil, false
		}

		panic(err)
	}

	user := &security.User{
		Name:       sp.getAssocAttribute(assertionInfo.Values, "name", assertionInfo.NameID),
		Roles:      []string{role},
		IsEnabled:  true,
		ExternalID: assertionInfo.NameID,
		Source:     security.SourceSaml,
		Firstname:  sp.getAssocAttribute(assertionInfo.Values, "firstname", ""),
		Lastname:   sp.getAssocAttribute(assertionInfo.Values, "lastname", ""),
		Email:      sp.getAssocAttribute(assertionInfo.Values, "email", ""),
	}
	err = sp.userProvider.Save(c, user)
	if err != nil {
		sp.logger.Err(err).Msg("SamlAcsHandler: userProvider Save error")
		panic(fmt.Errorf("cannot save user: %w", err))
	}

	return user, true
}

func (sp *serviceProvider) errorRedirect(c *gin.Context, relayUrl *url.URL, errorMessage string) {
	query := relayUrl.Query()
	query.Set("errorMessage", errorMessage)
	relayUrl.RawQuery = query.Encode()

	c.Redirect(http.StatusPermanentRedirect, relayUrl.String())
}

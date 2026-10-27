module git.canopsis.net/canopsis/canopsis-community/community/go-engines-community

// The go line declares a required minimum version of Go to use with this module. (https://golang.org/ref/mod#go-mod-file-go)
// The .env:GOLANG_VERSION be greater than or equal to the version below.
go 1.27

toolchain go1.27.3

// Note: External libs under GPL, AGPL or other viral licenses are not allowed here.
// Canopsis Pro contains Canopsis Community, and Canopsis Pro can't contain viral
// licenses, since it's a proprietary product.
//
// Please always maintain this file with the rules described here:
// https://git.canopsis.net/canopsis/canopsis-pro/-/issues/590

require (
	github.com/apognu/gocal v0.9.1
	github.com/beevik/etree v1.8.1
	github.com/brianvoe/gofakeit/v7 v7.18.0
	github.com/bsm/redislock v0.10.0
	github.com/casbin/casbin/v3 v3.11.0
	github.com/chenyahui/gin-cache v1.10.0
	github.com/coreos/go-oidc/v3 v3.21.0
	github.com/dlclark/regexp2 v1.12.0
	github.com/dop251/goja v0.0.0-20250309171923-bcd7cc6bf64c
	github.com/dustinkirkland/golang-petname v0.0.0-20261016130456-fdd95a3d0543
	github.com/gin-gonic/gin v1.12.0
	github.com/go-chi/chi/v5 v5.5.0
	github.com/go-ldap/ldap/v3 v3.4.15
	github.com/go-playground/locales v0.14.2
	github.com/go-playground/universal-translator v0.18.2
	github.com/go-playground/validator/v10 v10.30.5
	github.com/goccy/go-yaml v1.19.2
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golang-migrate/migrate/v4 v4.20.1
	github.com/gorilla/schema v1.4.1
	github.com/gorilla/securecookie v1.1.2
	github.com/gorilla/sessions v1.4.0
	github.com/gorilla/websocket v1.5.3
	github.com/jackc/pgx/v5 v5.11.0
	github.com/json-iterator/go v1.1.12
	github.com/kylelemons/godebug v1.1.0
	github.com/mailru/easyjson v0.9.2
	github.com/michaelklishin/rabbit-hole/v3 v3.6.0
	github.com/pelletier/go-toml/v2 v2.4.3
	github.com/prometheus/client_golang v1.25.0
	github.com/prometheus/procfs v0.22.0
	github.com/rabbitmq/amqp091-go v1.16.0
	github.com/redis/go-redis/v9 v9.23.0
	github.com/rs/zerolog v1.35.1
	github.com/russellhaering/gosaml2 v0.12.0 // since v0.10.0 assertions are not decrypted when SkipSignatureValidation is set, see decryptSAMLResponse in lib/api/auth/providers/saml.
	github.com/russellhaering/goxmldsig v1.6.1
	github.com/smartystreets/goconvey v1.8.1
	github.com/teambition/rrule-go v1.8.2
	github.com/valyala/fastjson v1.6.10
	go.mongodb.org/mongo-driver/v2 v2.9.2
	go.uber.org/mock v0.6.0
	go.yaml.in/yaml/v3 v3.0.5
	golang.org/x/crypto v0.58.0
	golang.org/x/oauth2 v0.37.0
	golang.org/x/sync v0.24.0
	golang.org/x/text v0.43.0
)

require (
	github.com/Azure/go-ntlmssp v0.1.1 // indirect
	github.com/ChannelMeter/iso8601duration v0.0.0-20150204201828-8da3af7a2a61 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.10.3 // indirect
	github.com/bytedance/gopkg v0.1.4 // indirect
	github.com/bytedance/sonic v1.15.4 // indirect
	github.com/bytedance/sonic/loader v0.5.3 // indirect
	github.com/casbin/govaluate v1.10.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.7 // indirect
	github.com/coreos/go-systemd/v22 v22.7.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.16 // indirect
	github.com/gin-contrib/sse v1.1.2 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.8 // indirect
	github.com/go-jose/go-jose/v4 v4.1.6 // indirect
	github.com/go-openapi/jsonpointer v1.0.2 // indirect
	github.com/go-openapi/jsonreference v1.0.3 // indirect
	github.com/go-openapi/spec v1.0.1 // indirect
	github.com/go-openapi/swag/conv v0.29.2 // indirect
	github.com/go-openapi/swag/jsonutils v0.29.2 // indirect
	github.com/go-openapi/swag/loading v0.29.2 // indirect
	github.com/go-openapi/swag/pools v0.29.2 // indirect
	github.com/go-openapi/swag/stringutils v0.29.2 // indirect
	github.com/go-openapi/swag/typeutils v0.29.2 // indirect
	github.com/go-openapi/swag/yamlutils v0.29.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sourcemap/sourcemap v2.1.4+incompatible // indirect
	github.com/goccy/go-json v0.11.2 // indirect
	github.com/google/pprof v0.0.0-20261016150608-1facf8f91b50 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.21.0 // indirect
	github.com/jackc/pgerrcode v0.0.0-20250907135507-afb5586c32a6 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.3 // indirect
	github.com/jellydator/ttlcache/v2 v2.11.1 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.20.3 // indirect
	github.com/klauspost/cpuid/v2 v2.4.1 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/common v0.72.1 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.64.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	github.com/swaggo/swag v1.16.6 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.2 // indirect
	github.com/urfave/cli/v2 v2.27.7 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xrash/smetrics v0.0.0-20261013231157-39c47d93efba // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	golang.org/x/arch v0.32.0 // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/tools v0.52.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

// No effect on the real canopsis-community repo, but necessary when it's part of the canopsis-pro monorepo
replace git.canopsis.net/canopsis/canopsis-community/community/go-engines-community => ./

tool (
	github.com/mailru/easyjson/easyjson
	github.com/swaggo/swag/cmd/swag
	go.uber.org/mock/mockgen
)

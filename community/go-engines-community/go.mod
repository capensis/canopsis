module git.canopsis.net/canopsis/canopsis-community/community/go-engines-community

// The go line declares a required minimum version of Go to use with this module. (https://golang.org/ref/mod#go-mod-file-go)
// The .env:GOLANG_VERSION be greater than or equal to the version below.
go 1.21

// Note: External libs under GPL, AGPL or other viral licenses are not allowed here.
// Canopsis Pro contains Canopsis Community, and Canopsis Pro can't contain viral
// licenses, since it's a proprietary product.
//
// Please always maintain this file with the rules described here:
// https://git.canopsis.net/canopsis/canopsis-pro/-/issues/590

require (
	// No GPL or AGPL libs allowed below!

	github.com/ajg/form v1.5.1
	github.com/alecthomas/participle v0.7.1
	github.com/apognu/gocal v0.9.0
	github.com/beevik/etree v1.2.0
	github.com/brianvoe/gofakeit/v6 v6.23.2
	github.com/bsm/redislock v0.9.4
	github.com/casbin/casbin/v2 v2.77.2
	github.com/chenyahui/gin-cache v1.8.1
	github.com/dlclark/regexp2 v1.10.0
	github.com/dop251/goja v0.0.0-20230828202809-3dbe69dd2b8e
	github.com/gin-gonic/gin v1.9.1
	github.com/go-ldap/ldap/v3 v3.4.5
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.15.3
	github.com/goccy/go-yaml v1.11.0
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/golang-migrate/migrate/v4 v4.16.2
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.1
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.5.0
	github.com/jackc/pgx/v5 v5.4.3
	github.com/json-iterator/go v1.1.12
	github.com/kylelemons/godebug v1.1.0
	github.com/mailru/easyjson v0.7.7
	github.com/pelletier/go-toml/v2 v2.1.0
	github.com/rabbitmq/amqp091-go v1.8.1
	github.com/redis/go-redis/v9 v9.1.0
	github.com/rs/zerolog v1.30.0
	github.com/russellhaering/gosaml2 v0.9.1
	github.com/russellhaering/goxmldsig v1.4.0
	github.com/smartystreets/goconvey v1.7.2
	github.com/teambition/rrule-go v1.8.2
	github.com/valyala/fastjson v1.6.4
	go.mongodb.org/mongo-driver v1.12.1
	golang.org/x/sync v0.3.0
	golang.org/x/text v0.14.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/coreos/go-oidc/v3 v3.6.0
	github.com/jellydator/ttlcache/v2 v2.11.1
	golang.org/x/crypto v0.17.0
	golang.org/x/oauth2 v0.13.0
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/ChannelMeter/iso8601duration v0.0.0-20150204201828-8da3af7a2a61 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/pprof v0.0.0-20230207041349-798e818bf904 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.0 // indirect
	github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.1 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jonboulle/clockwork v0.3.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/onsi/gomega v1.20.2 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.uber.org/atomic v1.10.0 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

// No effect on the real canopsis-community repo, but necessary when it's part of the canopsis-pro monorepo
replace git.canopsis.net/canopsis/canopsis-community/community/go-engines-community => ./

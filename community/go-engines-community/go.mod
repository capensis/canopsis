module git.canopsis.net/canopsis/canopsis-community/community/go-engines-community

// The go line declares a required minimum version of Go to use with this module. (https://golang.org/ref/mod#go-mod-file-go)
// The .env:GOLANG_VERSION be greater than or equal to the version below.
go 1.22

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
	github.com/beevik/etree v1.3.0
	github.com/brianvoe/gofakeit/v6 v6.28.0
	github.com/bsm/redislock v0.9.4
	github.com/casbin/casbin/v2 v2.87.1
	github.com/chenyahui/gin-cache v1.9.0
	github.com/dlclark/regexp2 v1.11.0
	github.com/dop251/goja v0.0.0-20240220182346-e401ed450204
	github.com/gin-gonic/gin v1.9.1
	github.com/go-ldap/ldap/v3 v3.4.7
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.19.0
	github.com/goccy/go-yaml v1.11.3
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/golang-migrate/migrate/v4 v4.17.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/securecookie v1.1.2
	github.com/gorilla/sessions v1.2.2
	github.com/gorilla/websocket v1.5.1
	github.com/jackc/pgx/v5 v5.5.5
	github.com/json-iterator/go v1.1.12
	github.com/kylelemons/godebug v1.1.0
	github.com/mailru/easyjson v0.7.7
	github.com/pelletier/go-toml/v2 v2.2.0
	github.com/rabbitmq/amqp091-go v1.9.0
	github.com/redis/go-redis/v9 v9.5.1
	github.com/rs/zerolog v1.32.0
	github.com/russellhaering/gosaml2 v0.9.1
	github.com/russellhaering/goxmldsig v1.4.0
	github.com/smartystreets/goconvey v1.8.1
	github.com/teambition/rrule-go v1.8.2
	github.com/valyala/fastjson v1.6.4
	go.mongodb.org/mongo-driver v1.14.0
	golang.org/x/sync v0.7.0
	golang.org/x/text v0.14.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/coreos/go-oidc/v3 v3.10.0
	github.com/jellydator/ttlcache/v2 v2.11.1
	github.com/prometheus/procfs v0.13.0
	golang.org/x/crypto v0.22.0
	golang.org/x/oauth2 v0.19.0
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/ChannelMeter/iso8601duration v0.0.0-20150204201828-8da3af7a2a61 // indirect
	github.com/bytedance/sonic v1.11.3 // indirect
	github.com/casbin/govaluate v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustinkirkland/golang-petname v0.0.0-20240428194347-eebcea082ee0 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sourcemap/sourcemap v2.1.4+incompatible // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/pprof v0.0.0-20240409012703-83162a5b38cd // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/onsi/gomega v1.20.2 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/smarty/assertions v1.15.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/arch v0.7.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

// No effect on the real canopsis-community repo, but necessary when it's part of the canopsis-pro monorepo
replace git.canopsis.net/canopsis/canopsis-community/community/go-engines-community => ./

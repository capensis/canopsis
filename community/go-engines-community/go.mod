module git.canopsis.net/canopsis/canopsis-community/community/go-engines-community

// Needs to be synced with Makefile.var:GOLANG_IMAGE_TAG and go.mod from Pro
// Only two integers are allowed here (https://golang.org/ref/mod#go-mod-file-go)
go 1.16

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
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/beevik/etree v1.1.0
	github.com/bsm/redislock v0.7.1
	github.com/casbin/casbin/v2 v2.36.1
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/cucumber/godog v0.12.0
	github.com/dlclark/regexp2 v1.4.0
	github.com/gin-gonic/gin v1.7.4
	github.com/go-chi/chi v4.0.2+incompatible // indirect
	github.com/go-ldap/ldap/v3 v3.4.1
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/go-redis/redis/v8 v8.11.3
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/json-iterator/go v1.1.11
	github.com/kr/pty v1.1.5 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pelletier/go-toml v1.9.3
	github.com/rs/zerolog v1.23.0
	github.com/russellhaering/gosaml2 v0.6.0
	github.com/russellhaering/goxmldsig v1.1.0
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.3.1
	github.com/swaggo/swag v1.6.7
	github.com/teambition/rrule-go v1.7.1
	github.com/tidwall/pretty v1.1.0 // indirect
	github.com/urfave/cli/v2 v2.3.0 // indirect
	github.com/valyala/fastjson v1.6.3
	go.mongodb.org/mongo-driver v1.7.1
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
	golang.org/x/text v0.3.7
	golang.org/x/tools v0.1.5 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// No effect on the real canopsis-community repo, but necessary when it's part of the canopsis-pro monorepo
replace git.canopsis.net/canopsis/canopsis-community/community/go-engines-community => ./

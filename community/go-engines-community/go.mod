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
	github.com/casbin/casbin/v2 v2.31.4
	github.com/cucumber/godog v0.11.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/dlclark/regexp2 v1.4.0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-ldap/ldap/v3 v3.3.0
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.6.1
	github.com/go-redis/redis/v8 v8.10.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.2.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.11
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pelletier/go-toml v1.9.3
	github.com/rs/zerolog v1.23.0
	github.com/russellhaering/gosaml2 v0.6.0
	github.com/russellhaering/goxmldsig v1.1.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/streadway/amqp v1.0.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	github.com/teambition/rrule-go v1.7.0
	github.com/tidwall/pretty v1.1.0 // indirect
	github.com/valyala/fastjson v1.6.3
	github.com/vmihailenco/msgpack/v4 v4.3.12
	go.mongodb.org/mongo-driver v1.5.3
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/text v0.3.6
	gopkg.in/yaml.v2 v2.4.0
)

// No effect on the real canopsis-community repo, but necessary when it's part of the canopsis-pro monorepo
replace git.canopsis.net/canopsis/canopsis-community/community/go-engines-community => ./

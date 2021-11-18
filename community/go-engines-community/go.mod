module git.canopsis.net/canopsis/go-engines

// If there's a hard requirement, always sync this and Makefile.var:GOLANG_IMAGE_TAG
// Only two integers are allowed here (https://golang.org/ref/mod#go-mod-file-go)
go 1.15

// Note: GPL, AGPL (and other viral) libs are not allowed here, because of CAT.
// Please always maintain this file with the rules described here:
// https://git.canopsis.net/canopsis/canopsis/-/issues/2697
// TODO Following dependencies updates require code changes :
//      github.com/go-redis/redis/v7 cannot be updated to github.com/go-redis/redis/v8
//          because method signatures have been changed
//      github.com/bsm/redislock v0.5.0 cannot be updated
//          because newer versions use github.com/go-redis/redis/v8
//      go.mongodb.org/mongo-driver v1.3.7 cannot be updated to v1.4
//          because error wrapping has been added and errors are checking in
//          UnmarshalBSONValue methods as err.(CustomError) instead of errors.As(err, &CustomError{})
require (
	github.com/BurntSushi/toml v0.3.1 // UNMAINTAINED
	github.com/alecthomas/participle v0.7.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/bsm/redislock v0.5.0
	github.com/casbin/casbin/v2 v2.19.4
	github.com/cucumber/godog v0.10.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/dlclark/regexp2 v1.4.0
	github.com/gin-gonic/gin v1.6.3
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v7 v7.4.0
	github.com/golang/mock v1.4.4
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/google/uuid v1.1.2
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/influxdata/influxdb v1.8.3
	github.com/json-iterator/go v1.1.10
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/streadway/amqp v1.0.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	github.com/teambition/rrule-go v1.6.2
	github.com/tidwall/gjson v1.6.4
	github.com/vmihailenco/msgpack/v4 v4.3.7
	go.mongodb.org/mongo-driver v1.3.7
	golang.org/x/net v0.0.0-20201216054612-986b41b23924
	golang.org/x/text v0.3.4
	gopkg.in/yaml.v2 v2.4.0
)

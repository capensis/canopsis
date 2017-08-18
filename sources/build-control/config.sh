# - Variable UPDATE_ETC_ACTION
# - Possible values : keep replace
#
# If defined, the user won't be asked what to do on config file update

UPDATE_ETC_ACTION="ask"

# - Variable MONGODB_ADMIN_PASSWD
# If defined, the user won't be asked to type a password for the MongoDB admin user

MONGODB_ADMIN_PASSWD="admin"

# - Variable MONGODB_USER
# If defined, the user won't be asked to enter a username for MongoDB

MONGODB_USER="cpsmongo"

# - Variable MONGODB_USER_PASSWD
# If defined, the user won't be asked to type a password for the MongoDB canopsis user

MONGODB_USER_PASSWD="canopsis"

# - Variable INFLUXDB_ADMIN_PASSWD
# If defined, the user won't be asked to type a password for the InfluxDB admin user

INFLUXDB_ADMIN_PASSWD="admin"

# - Variable INFLUXDB_USER
# If defined, the user won't be asked to enter a username for InfluxDB

INFLUXDB_USER="cpsinflux"

# - Variable INFLUXDB_USER_PASSWD
# If defined, the user won't be asked to type a password for the InfluxDB canopsis user

INFLUXDB_USER_PASSWD="canopsis"

# - Variable RABBITMQ_ADMIN_PASSWD
# If defined, the user won't be asked to type a password for the RabbitMQ admin user

RABBITMQ_ADMIN_PASSWD="admin"

# - Variable RABBITMQ_USER
# If defined, the user won't be asked to enter a username for RabbitMQ

RABBITMQ_USER="cpsrabbit"

# - Variable RABBITMQ_USER_PASSWD
# If defined, the user won't be asked to type a password for the RabbitMQ canopsis user

RABBITMQ_USER_PASSWD="canopsis"

# - Variable SSL_CA_PATH
# If defined and path exists, will copy the CA to Canopsis instead of generate a new one
# SSL_CA_PATH=""

# - Variable SSL_CAKEY_PATH
# If defined and path exists, will copy the CA key to Canopsis instead of generate a new one
# SSL_CAKEY_PATH=""

# - Variable SSL_KEY_PASS
# Password for CA certificate key if we are generating one. Will prompt the user if empty
SSL_KEY_PASS="cpsnode"

# - Variable SSL_KEY_BITS
# Key size in bits
SSL_KEY_BITS="1024"

# - Variable SSL_CHECK_SECONDS
# Range used to check for certificate expiration from now
SSL_CHECK_SECONDS="86400"

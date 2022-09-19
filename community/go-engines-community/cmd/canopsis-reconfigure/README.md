# Reconfigure.

The *canopsis-reconfigure* command is used to initialize the RabbitMQ queues and exchanges to using a configuration and
initialize the Mongo database.

- [Run](#run)
    - [Arguments](#arguments)
    - [Environment vars](#environment-vars)
- [Configuration](#configuration)
- [Migrations](#migrations)
- [Fixtures](#fixtures)
  - [Creating fixtures](#creating-fixtures)
    - [Functions](#functions)
    - [References](#references)
    - [Templates](#templates)

## Run

#### Arguments

Each command argument has default value. Change them to fit your needs.

- `conf` - The configuration file used to initialize Canopsis.
- `migrate-mongo` - If true, it will execute Mongo migration scripts.
- `mongo-migration-directory` - The directory with Mongo migration scripts.
- `mongo-fixture-directory` - The directory with Mongo fixtures.
- `migrate-postgres` - If true, it will execute Postgres migration scripts.
- `postgres-migration-directory` - The directory with Postgres migration scripts.
- `postgres-migration-mode` - Should be up or down.
- `postgres-migration-steps` - Number of migration steps, will execute all migrations if empty or 0.

#### Environment vars

Before run define following env vars :

- `CPS_MONGO_URL` - Mongo DB connection url.
- `CPS_AMQP_URL` - RMQ connection url.
- `CPS_MAX_RETRY` - Max attempts to initialize app.
- `CPS_MAX_DELAY` - Delay between attempt.
- `CPS_WAIT_FIRST_ATTEMPT` - Timeout before first initialize attempt.
- `CPS_POSTGRES_URL` - Postgres DB connection url.

## Configuration

Use [canopsis-community.toml](./canopsis-community.toml) as configuration template for community version
and [canopsis-pro.toml](./canopsis-pro.toml) for pro version.

Rerun *canopsis-reconfigure* after configuration modification. Rerun API or engines if modified parameters aren't reloaded by them automatically.  

## Migrations

Migrations are used for versioning database schema. Migrations provide easy and safe way to deploy database schema changes.
Migration scripts contain Javascript code which is executed in Mongo console.  

Cmd *canopsis-reconfigure* automatically migrates new migrations from [database/migrations](../../database/migrations) directory.

Use [cmd/mongo-migrations](../mongo-migrations) to manage migrations. 

## Fixtures

Fixtures are used to create predefined data for new production environment and to create fake data for functional testing.
Fixtures are YAML files which content is inserted to database collections. 

If Canopsis database is empty *canopsis-reconfigure* automatically applies fixtures from [database/fixtures](../../database/fixtures) directory.

#### Creating fixtures

Many documents of different collections can be defined in one file.

```yaml
users:
  user1:
    name: John
  user2:
    name: Bob

roles:
  role_admin:
    name: admin
  role_support:
    name: support
```

Here `users` and `roles` are collection names.

##### Functions

Different functions can be used to generate field values.

```yaml
users:
  user1:
    _id: <UUID()>
    name: <Name()>
    password: <Password(test)>
    created: <NowUnix()>
    birthday: <DateUnix()>
```

Following functions are available

- `<Password(pass)>` returns hashed password.
- `<NowUnix()>` returns Unix timestamp for current time.
- `<DateUnix()>` returns Unix timestamp for random time.
- `<Index>` returns index of current document in all documents which are defined under collection name.
- `<Current().field>` returns value of `field` of current document which is defined above.

Check [gofakeit](https://github.com/brianvoe/gofakeit) for more functions. Arguments of function can be only literals.

##### References

Identifier of previously defined document in the same file can be used as value of field.

```yaml
roles:
  role_admin:
    _id: <UUID()>
    name: admin

users:
  user1:
    _id: <UUID()>
    name: John
    role: "@role_admin"  
```

##### Templates

Use template to not repeat fields in documents.

```yaml
template:
  - &default_user {
    _id: <UUID()>,
    created: <NowUnix()>,
    role: "@role_admin"
  }

users:
  user1:
    <<: *default_user
    name: John
    password: <Password(test)>

  user2:
    <<: *default_user
    name: Bob
    password: <Password(test)>
    role: "@role_support"
```

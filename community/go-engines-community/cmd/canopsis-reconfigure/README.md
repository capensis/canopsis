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
- `migration-directory` - The directory with migration scripts.
- `fixture-directory` - The directory with fixtures.

#### Environment vars

Before run define following env vars :

- `CPS_MONGO_URL` - Mongo DB connection url.
- `CPS_AMQP_URL` - RMQ connection url.
- `CPS_MAX_RETRY` - Max attempts to initialize app.
- `CPS_MAX_DELAY` - Delay between attempt.
- `CPS_WAIT_FIRST_ATTEMPT` - Timeout before first initialize attempt.

## Configuration

Use [canopsis-core.toml.example](./canopsis-core.toml.example) as configuration template for community version
and [canopsis-cat.toml.example](./canopsis-cat.toml.example) for pro version.

Rerun *canopsis-reconfigure* after configuration modification. Rerun API or engines if modified parameters aren't reloaded by them automatically.  

## Migrations

Migrations are used for versioning database schema. Migrations provide easy and safe way to deploy database schema changes.
Migration scripts contain Javascript code which is executed in Mongo console.  

Cmd *canopsis-reconfigure* automatically migrates new migrations from [database/migrations](../../database/migrations) directory.

Use [cmd/mongo-migrations](../mongo-migrations) to manage migrations. 

## Fixtures

Fixtures are used to create predefined data for new production environment and to create fake data for functional testing.
Fixtures are YAML files which content is inserted to database collections. 

If Canopsis database is empty *canopsis-reconfigure* automatically applies fixtures from [database/fixtures/prod](../../database/fixtures/prod) directory.

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
    name: admin

users:
  user1:
    name: John
    role: "@role_admin"  
```

##### Templates

Use template to not repeat fields in documents.

```yaml
users:
  default_user (template):
    _id: <UUID()>
    created: <NowUnix()>
    role: "@role_admin"

  user1 (extend default_user):
    name: John
    password: <Password(test)>

  user2 (extend default_user):
    name: Bob
    password: <Password(test)>
    role: "@role_support"
```

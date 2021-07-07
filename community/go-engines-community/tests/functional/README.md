# Functional tests

This folder contains functional tests for Go API and Go engines.

Function tests are implemented using [Godog](https://github.com/cucumber/godog) BDD framework. It uses Gherkin formatted scenarios in the format of Given, When, Then.

- [Run](#run)
    - [Environment vars](#environment-vars)
- [Database](#database)
- [API](#api)
- [Engines](#engines)
    - [Events](#events)
    - [Periodical process](#periodical-process)
    - [Engine-webhook](#engine-webhook)
- [Best practices](#best-practices)

## Run

Execute Go test command to run functional tests. Use required `paths` argument to run selected tests. 

All api tests from [api](./features/api) folder can be executed by  

```bash
go test --paths feauture/api
```  

All engines tests from [engines](./features/engines) folder can be executed by  

```bash
go test --paths feauture/engines
``` 

#### Environment vars

Before run test define following env vars :

- `API_URL` - Go API url
- `CPS_MONGO_URL` - Mongo DB connection url
- `CPS_AMQP_URL` - RMQ connection url
- `CPS_REDIS_URL` - Redis connection url

**Be sure to run application before go test. Go test doesn't run anything except tests.**

## Database

Fixtures are used to load data into a database that can then be used int tests.
Fixtures are stored in [fixtures](../../fixtures) folder. Each file in folder is named like corresponding mongo collection.
A file contains json with list of data.

Fixtures are reloaded before and after test suite.

Env var `CPS_MONGO_URL` is used to connect to Mongo DB.

**Be sure to run separate database instance. All data will be overwritten during tests.**  

## Redis

Redis cache is flushed before and after test suite.

Env var `CPS_REDIS_URL` is used to connect to Redis.

## API

Run API before execute tests.

API tests are in [features/api](./features/api) folder.

Env var `API_URL` is used to execute API requests.

## Engines

Run Engines before execute tests.

### Events

Testing app receives `FIFO Ack` events to detect the end of event processing.
Prepare testing environment so testing app can receive these events :

- Bind `FIFO ack queue` to `amq.direct` exchange

Init command can be used. Update init config as following and run init command. 

```toml
[[RabbitMQ.queues]]
name = "FIFO_ack"
durable = true
autoDelete = false
exclusive = false
noWait = false
  [RabbitMQ.queues.bind]
  key = "FIFO_ack"
  exchange = "amq.direct"
  noWait = false
```

- Run `action` and `che` engines with argument `-fifoAckExchange amq.direct` 

Arguments `-ewe amq.direct` and `-ewk FIFO_ack` provides queue info to test app. They already have default values.    

Env var `CPS_AMQP_URL` is used to connect to RMQ.  

Argument `--eventlogs events.log` logs all events which were catched by `When I wait the end of event processing`
and `When I wait the end of N events processing` for each test scenario.

Argument `--checkUncaughtEvents` enables catching event after each scenario. Functional tests stop by panic if event is caught. Please note, it adds extra timeout for each scenario.

### Periodical process

Run all engines with shorter periodical wait time. Use argument `-periodicalWaitTime 2s`.

Argument `-pwt 2200ms` provides interval which test app should wait before the end of next periodical process.

### Engine-webhook

API must be accessible by Engine-webhook to execute webhook tests. 

If all engines and API run by IDE use `API_URL=localhost:8082` env var.

If all engine and API are run by docker add
```
api 127.0.0.1
```
to `/etc/hosts` and use `API_URL=api:8082` env var.

## Best practices

To avoid conflicts between test scenarios :

- Use unique names for rules, alarms, entities, etc. which are used in each test scenario.
- Catch all events which scenario causes.
- If some rules are used to generate events (idle rules, etc.) make sure create
  them with parameters with which events aren't generated during tests. 

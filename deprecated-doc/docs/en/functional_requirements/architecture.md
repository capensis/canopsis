# Canopsis Architecture

This document describes the architecture requirements of Canopsis.

## Contents

### Description

Because Canopsis is built on top of several technologies, it has
specific requirements.

To know those requirements we have to know its needs :

- Canopsis sources will gather informations to send them as message to Canopsis
- it will need a messaging queue system in order to transport messages
- in order to store incoming data, it will need a database platform
- to expose data, it will need a web server to process requests

### Data Source

A data source MUST be independent and MUST NOT rely on Canopsis
toolchain to produce messages.

### Messaging Queue System

Each Canopsis engine listen on a queue and chain the message to the 
next configured queues, each one of them associated to an engine.

The technical solution MUST support :

- message queuing
- load-balanced consumers
- consume to all consumers
- consumers subscription to a set of messages
- authentication support

### Data storage

The data storage is the main bridge between the engines and the front-end.

The technical solution MUST support :

- API with JSON support
- indexing
- in memory cache
- authentication support
- file storage support

### Data exposure

Data will be exposed to the front-end via an API which MUST support :

- JSON custom protocol
- authentication (with credentials or token)
- external authentication (with CAS, LDAP)

## Schema

![image](../../img/requirement/architecture.png)

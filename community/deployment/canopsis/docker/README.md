# Canopsis Community Docker Compose environment

## Starting the environment

In order to start the stack, you need to set env variable `CPS_EDITION` to
`community`:

```bash
export CPS_EDITION=community
docker compose up -d
```

or:

```bash
CPS_EDITION=community docker compose up -d
```

More information on the [official documentation (french)][doc].

[doc]: https://doc.canopsis.net/guide-administration/installation/installation-conteneurs/

# Upgrading to Canopsis 2.7

## Redis

Canopsis now requires a Redis server.

You now need to export the following environment variable:
```bash
CPS_REDIS_URL=redis://REDIS_SERVER_ADDRESS:6379/0
```

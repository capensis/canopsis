## Deprecated AXE flags
The `featureHideResources` and `postProcessorsDirectory` flags are marked as  deprecated, since those features were removed from `engine-axe`.

## New FIFO flag
The `enableMetaAlarmProcessing` flag was added to the `engine-fifo`. It's `true` by default. Since `engine-correlation` was moved to `cat` version it should be possible to turn off meta-alarm processing from `engine-fifo`. If you use `core` version, set `enableMetaAlarmProcessing` to `false`

## Multi-axe concurrency

The one thing instances compete for is periodical process. Only one instance should do periodical process when it ticks. It is resolved by `redlock`. If an instance acquire the lock it does the periodical process or skips it otherwise. `"Obtain redis lock: unexpected error"` printed to log in case of other error.

## Start multi-axe

To start with engine-axe multiple instances docker-compose can be run with command line arguments:  
`--scale axe=N` where N is number of engine-axe instances.  

## Testing Axe Multi-instance
It is possible to test with `random-feeder`. Also it's possible to use a new flag:

`alarms N` -- where N should be from 1-100. The N means a percentage of alarms, which changes their state from OK to CRITICAL or vice versa.
BEGIN;

ALTER TABLE sli_duration SET (timescaledb.compress = true);
SELECT add_compression_policy('sli_duration', INTERVAL '7 days');

COMMIT;

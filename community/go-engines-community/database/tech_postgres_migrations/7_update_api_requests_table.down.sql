BEGIN;

ALTER TABLE api_requests
    SET (timescaledb.compress = false);
ALTER TABLE api_requests ALTER COLUMN url TYPE VARCHAR(255);
ALTER TABLE api_requests
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'url');

COMMIT;

BEGIN;

SELECT remove_compression_policy('sli_duration');

CREATE TEMP TABLE sli_duration_compressed_chunks AS
SELECT chunk_schema || '.' || chunk_name AS full_chunk_name
FROM chunk_compression_stats('sli_duration')
WHERE compression_status = 'Compressed';

DO $$
    DECLARE
        chunk TEXT;
    BEGIN
        FOR chunk IN SELECT full_chunk_name FROM sli_duration_compressed_chunks
            LOOP
                PERFORM decompress_chunk(chunk);
            END LOOP;
    END;
$$;

ALTER TABLE sli_duration SET (timescaledb.compress = false);

COMMIT;

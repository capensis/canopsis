BEGIN;

CREATE TABLE IF NOT EXISTS llm_token_usage
(
    time               TIMESTAMP    NOT NULL,
    model              VARCHAR(255) NOT NULL,
    input_tokens       INT          NOT NULL DEFAULT 0,
    output_tokens      INT          NOT NULL DEFAULT 0,
    thinking_tokens    INT          NOT NULL DEFAULT 0,
    cache_write_tokens INT          NOT NULL DEFAULT 0,
    cache_read_tokens  INT          NOT NULL DEFAULT 0,
    retries            INT          NOT NULL DEFAULT 0
);
SELECT create_hypertable('llm_token_usage', 'time', if_not_exists => TRUE);

CREATE MATERIALIZED VIEW IF NOT EXISTS llm_token_usage_hourly
            (
             time,
             model,
             input_tokens,
             output_tokens,
             thinking_tokens,
             cache_write_tokens,
             cache_read_tokens,
             retries,
             count
                )
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time),
       model,
       sum(input_tokens),
       sum(output_tokens),
       sum(thinking_tokens),
       sum(cache_write_tokens),
       sum(cache_read_tokens),
       sum(retries),
       count(*)
FROM llm_token_usage
GROUP BY time_bucket('1 hour', time), model
    WITH NO DATA;

SELECT add_continuous_aggregate_policy(
           'llm_token_usage_hourly',
           start_offset => INTERVAL '3 hours',
           end_offset => INTERVAL '1 hour',
           schedule_interval => INTERVAL '1 hour',
           if_not_exists => TRUE
       );
SELECT add_retention_policy('llm_token_usage', drop_after => INTERVAL '24 hours', if_not_exists => TRUE);

ALTER MATERIALIZED VIEW llm_token_usage_hourly SET (timescaledb.compress = true);
SELECT add_compression_policy('llm_token_usage_hourly', compress_after => '1 day'::interval, if_not_exists => TRUE);
SELECT add_retention_policy('llm_token_usage_hourly', drop_after => INTERVAL '14 days', if_not_exists => TRUE);

END;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'metric_type') THEN
CREATE TYPE metric_type AS ENUM ('counter', 'gauge');
END IF;
END$$;

CREATE TABLE IF NOT EXISTS metrics.metrics(
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    metric_type metric_type NOT NULL,
    delta BIGINT default NULL,
    value DOUBLE PRECISION default NULL
);

---- create above / drop below ----
DROP TABLE IF EXISTS metrics.metrics;

DROP TYPE IF EXISTS metrics.metric_type;
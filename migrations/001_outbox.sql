CREATE TABLE IF NOT EXISTS outbox (
    id         UUID        PRIMARY KEY,
    type       TEXT        NOT NULL,
    content    JSONB       NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
);

-- índice parcial: o worker só varre linhas pendentes (updated_at IS NULL),
-- então o índice fica pequeno e o poll rápido mesmo com a tabela grande.
CREATE INDEX IF NOT EXISTS idx_outbox_pending
    ON outbox (created_at)
    WHERE updated_at IS NULL;

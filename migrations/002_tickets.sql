CREATE TABLE IF NOT EXISTS tickets (
    id       TEXT          PRIMARY KEY,
    event_id TEXT          NOT NULL REFERENCES events(id),
    price    NUMERIC(10,2) NOT NULL,
    user_id  TEXT          NOT NULL,
    status   TEXT          NOT NULL,
    seat     TEXT          NOT NULL
);

-- busca rápida de todos os tickets de um evento (relação 1->N).
CREATE INDEX IF NOT EXISTS idx_tickets_event_id ON tickets (event_id);

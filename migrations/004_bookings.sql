CREATE TABLE IF NOT EXISTS bookings (
    id         TEXT        PRIMARY KEY,
    ticket_id  TEXT        NOT NULL UNIQUE REFERENCES tickets(id),
    seat       TEXT        NOT NULL,
    user_id    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ticket_id UNIQUE: um ticket é reservado por exatamente um usuário. Também dá
-- idempotência ao consumer (reprocessar a mesma msg não duplica o booking).
-- Índice para "meus bookings" (buscar por usuário).
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings (user_id);

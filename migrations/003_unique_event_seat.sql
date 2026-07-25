-- Remove duplicatas pré-existentes (mantém 1 linha por event_id+seat),
-- senão o índice único abaixo falharia ao ser criado.
DELETE FROM tickets a
USING tickets b
WHERE a.ctid < b.ctid
  AND a.event_id = b.event_id
  AND a.seat = b.seat;

-- Garante unicidade de assento POR evento. Índice único (em vez de ADD
-- CONSTRAINT) para poder usar IF NOT EXISTS = idempotente. Violação gera 23505.
CREATE UNIQUE INDEX IF NOT EXISTS uq_event_seat ON tickets (event_id, seat);

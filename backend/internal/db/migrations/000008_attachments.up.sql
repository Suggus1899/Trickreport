-- Migration: 008_attachments
-- Creates: ticket_attachments

CREATE TABLE ticket_attachments (
    id           UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id    UUID          NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id      UUID          NOT NULL REFERENCES users(id),
    filename     VARCHAR(255)  NOT NULL,
    content_type VARCHAR(100)  NOT NULL,
    file_size    BIGINT        NOT NULL,
    file_data    BYTEA         NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX ticket_attachments_ticket_idx ON ticket_attachments(ticket_id);

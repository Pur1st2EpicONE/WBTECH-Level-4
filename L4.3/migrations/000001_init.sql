-- +goose Up
CREATE TABLE IF NOT EXISTS archived_events (
    event_id UUID PRIMARY KEY,
    user_id INT NOT NULL,
    event_date DATE NOT NULL,
    text TEXT,
    archived_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_archived_user_date ON archived_events(user_id, event_date);

-- +goose Down
DROP TABLE IF EXISTS archived_events CASCADE;

CREATE TABLE active_sessions (
    session_id CHAR(86) PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX user_id_id ON active_sessions (user_id)

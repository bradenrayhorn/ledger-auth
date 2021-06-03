CREATE TABLE active_sessions (
    session_id char(86) PRIMARY KEY,
    user_id char(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    INDEX user_id_idx (user_id)
);

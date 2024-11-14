-- db/migrations/messages_index_nnn.sql

DROP TABLE IF EXISTS messages_index;

CREATE TABLE IF NOT EXISTS messages_index (
    msg_id SERIAL PRIMARY KEY,
    msg_platform VARCHAR(200) DEFAULT NULL,
    msg_external_id VARCHAR(200) DEFAULT NULL,
    msg_created_at TIMESTAMP DEFAULT NULL,
    msg_language VARCHAR(200) DEFAULT NULL,
    msg_url VARCHAR(200) DEFAULT NULL,
    msg_content TEXT DEFAULT NULL,
    msg_external_account_id VARCHAR(200) DEFAULT NULL
);

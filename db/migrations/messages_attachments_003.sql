-- db/migrations/messages_attachments_nnn.sql

DROP TABLE IF EXISTS messages_attachments;

CREATE TABLE IF NOT EXISTS messages_attachments (
    attachment_id SERIAL PRIMARY KEY,
    attachment_msg_id INT DEFAULT NULL,
    attachment_external_id VARCHAR(200) DEFAULT NULL,
    attachment_url VARCHAR(200) DEFAULT NULL,
    attachment_type VARCHAR(200) DEFAULT NULL,
    attachment_meta_description TEXT DEFAULT NULL
);

CREATE TABLE messages(
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    uuid varchar(36) NOT NULL,
    from_user varchar(36) NOT NULL,
    to_user varchar(36),
    content text,
    message_type smallint,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
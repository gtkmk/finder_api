CREATE TABLE user_event (
    id VARCHAR(191) PRIMARY KEY,
    user_id VARCHAR(191) NOT NULL,
    event ENUM('login', 'new_user', 'edit_user'),
    old_value MEDIUMTEXT,
    new_value MEDIUMTEXT,
    ip VARCHAR(90),
    device VARCHAR(191),

    created_at DATETIME(3) NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES user (id)
);
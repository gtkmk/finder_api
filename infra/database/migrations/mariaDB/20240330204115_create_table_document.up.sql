CREATE TABLE IF NOT EXISTS document (
    id VARCHAR(191) PRIMARY KEY,
    type ENUM('media', 'profile_picture', 'profile_banner_picture') NOT NULL,
    path VARCHAR(250) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    post_id VARCHAR(191),
    owner_id VARCHAR(191) NOT NULL,
    deleted_reason VARCHAR(100),
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    CONSTRAINT fk_user_document FOREIGN KEY (owner_id) REFERENCES user (id),
    CONSTRAINT fk_post_document FOREIGN KEY (post_id) REFERENCES post (id)
);

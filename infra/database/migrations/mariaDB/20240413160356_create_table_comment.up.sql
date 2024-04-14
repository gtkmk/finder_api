CREATE TABLE IF NOT EXISTS comment (
    id VARCHAR(191) PRIMARY KEY,
    post_id VARCHAR(191) NOT NULL,
    user_id VARCHAR(191) NOT NULL,
    text TEXT NOT NULL,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);
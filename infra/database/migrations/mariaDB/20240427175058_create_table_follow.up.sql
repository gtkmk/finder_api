CREATE TABLE IF NOT EXISTS follow (
    id VARCHAR(191) PRIMARY KEY,
    follower_id VARCHAR(191) NOT NULL,
    followed_id VARCHAR(191) NOT NULL,
    created_at DATETIME(3) NULL,
    FOREIGN KEY (follower_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE INDEX idx_follower_id ON follow(follower_id);
CREATE INDEX idx_followed_id ON follow(followed_id);

CREATE TABLE IF NOT EXISTS interaction_likes (
    id VARCHAR(191) PRIMARY KEY,
    like_type ENUM('post', 'comment') NOT NULL,
    post_id VARCHAR(191),
    comment_id VARCHAR(191),
    user_id VARCHAR(191) NOT NULL,
    created_at DATETIME(3) NULL,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES `comment`(id) ON DELETE CASCADE
);
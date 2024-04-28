CREATE TABLE IF NOT EXISTS follow (
    id INT AUTO_INCREMENT PRIMARY KEY,
    follower_id INT NOT NULL,
    followed_id INT NOT NULL,
    created_at DATETIME(3) NULL,
    FOREIGN KEY (follower_id) REFERENCES user(id),
    FOREIGN KEY (followed_id) REFERENCES user(id)
);

CREATE INDEX idx_follower_id ON follow(follower_id);
CREATE INDEX idx_followed_id ON follow(followed_id);

CREATE TABLE IF NOT EXISTS post (
    id VARCHAR(191) NOT NULL,
    text TEXT DEFAULT NULL,
    location VARCHAR(255) DEFAULT NULL,
    reward TINYINT(1) DEFAULT NULL,
    lost_found ENUM('lost', 'found') DEFAULT NULL,
    privacy ENUM('public', 'private', 'friends_only') DEFAULT 'public',
    shares_count INT(11) DEFAULT 0,
    category ENUM('default', 'paid', 'add') DEFAULT 'default',
    animal_type ENUM('cachorro', 'gato', 'ave', 'outro') DEFAULT NULL,
    animal_size ENUM('pequeno', 'médio', 'grande') DEFAULT NULL,
    user_id VARCHAR(191) DEFAULT NULL,
    deleted_reason VARCHAR(100) DEFAULT NULL,
    created_at DATETIME(3) DEFAULT NULL,
    updated_at DATETIME(3) DEFAULT NULL,
    deleted_at DATETIME(3) DEFAULT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES user(id)
);

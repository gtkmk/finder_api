CREATE TABLE IF NOT EXISTS user (
    id VARCHAR(191) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    cpf varchar(14) NOT NULL,
    cellphone_number varchar(15),
    status ENUM('expired', 'pending', 'logged') NOT NULL,
    is_active TINYINT(1) DEFAULT 1 NOT NULL,
    password_reset TINYINT(1) DEFAULT 0,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    CONSTRAINT email UNIQUE (email)
);

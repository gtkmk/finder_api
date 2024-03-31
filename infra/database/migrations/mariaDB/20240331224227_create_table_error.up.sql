CREATE TABLE error (
   id VARCHAR(191) PRIMARY KEY,
   message TEXT NOT NULL,
   stack LONGTEXT NULL,
   created_at DATETIME(3) NULL
);

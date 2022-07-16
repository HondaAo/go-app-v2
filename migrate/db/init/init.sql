USE `videos-app`;

CREATE USER 'root'@'%' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;

CREATE TABLE IF NOT EXISTS `user` (
    `user_id`   VARCHAR(100) PRIMARY KEY,
    `first_name`   VARCHAR(32) NOT NULL,
    `last_name`   VARCHAR(32) NOT NULL,
    `email`   VARCHAR(64) UNIQUE NOT NULL,
    `password`   VARCHAR(250) NOT NULL,
    `role`   VARCHAR(10) NOT NULL DEFAULT 'user',
    `country`   VARCHAR(30),
    `created_at`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
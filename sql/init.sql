-- db for user
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `user_id` VARCHAR(37) NOT NULL,
  `user_name` VARCHAR(30) NOT NULL,
  `password` VARCHAR(90) NOT NULL,
  `prefect`  VARCHAR(30) NOT NULL,
  `gender`   VARCHAR(30) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_name` (`user_name`)
);

-- db for logic delete info
DROP TABLE IF EXISTS `user_deletes`;
CREATE TABLE `user_deletes` (
  `flag` SMALLINT DEFAULT 0,
  `user_id` VARCHAR(37) NOT NULL,
  PRIMARY KEY (`user_id`),
);

-- db for private and open chat content
DROP TABLE IF EXISTS `chats`;
CREATE TABLE `chats` (
  `chat_id` VARCHAR(37) NOT NULL,
  `post` VARCHAR(100) NOT NULL,
  `postuser_id`  VARCHAR(37) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`chat_id`),
);

-- db for chat owner
DROP TABLE IF EXISTS `chat_ownerships`;
CREATE TABLE `chat_ownerships` (
  `chat_id` VARCHAR(37) NOT NULL,
  `user_id` VARCHAR(37) NOT NULL,
  PRIMARY KEY (`chat_id`, `user_id`),
);

-- db for session info
DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
  `session_id` VARCHAR(37) NOT NULL,
  `user_id` VARCHAR(37) NOT NULL,
  PRIMARY KEY (`session_id`)
);
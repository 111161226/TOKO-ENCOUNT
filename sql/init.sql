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
  PRIMARY KEY (`user_id`)
);

-- db for private and open chat content
DROP TABLE IF EXISTS `chats`;
CREATE TABLE `chats` (
  `chat_id` VARCHAR(37) NOT NULL,
  `room_id` VARCHAR(37) NOT NULL, -- 全体チャットは0
  `destination_user_id` VARCHAR(37) NOT NULL, -- 全体チャットは0
  `post` VARCHAR(100) NOT NULL,
  `post_user_id`  VARCHAR(37) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`chat_id`)
);

-- db for room data
DROP TABLE IF EXISTS `room_datas`;
CREATE TABLE `room_datas` (
  `room_id` VARCHAR(37) NOT NULL,
  `user_id` VARCHAR(37) NOT NULL,
  `latest_access` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `not_read` BIGINT(100) NOT NULL DEFAULT 0,
  PRIMARY KEY (`room_id`, `user_id`)
);

-- db for session info
DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
  `session_id` VARCHAR(37) NOT NULL,
  `user_id` VARCHAR(37) NOT NULL,
  PRIMARY KEY (`session_id`)
);

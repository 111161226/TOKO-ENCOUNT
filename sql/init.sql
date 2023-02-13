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
  `user_name` VARCHAR(37) NOT NULL,
  PRIMARY KEY (`user_name`)
);

-- db for private and open chat content
DROP TABLE IF EXISTS `chats`;
CREATE TABLE `chats` (
  `chat_id` VARCHAR(37) NOT NULL,
  `room_id` VARCHAR(37) NOT NULL, -- 全体チャットは0
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

-- db for room name
DROP TABLE IF EXISTS `room_names`;
CREATE TABLE `room_names` (
  `room_id` VARCHAR(37) NOT NULL,
  `room_name` VARCHAR(30) NOT NULL,
  PRIMARY KEY (`room_id`)
);

-- db for session info
DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
  `session_id` VARCHAR(37) NOT NULL,
  `user_id` VARCHAR(37) NOT NULL,
  PRIMARY KEY (`session_id`)
);

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` (`user_id`, `user_name`, `password`, `prefect`, `gender`)
VALUES
  ('system', 'system', 'system', 'system', 'system');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;


LOCK TABLES `chats` WRITE;
/*!40000 ALTER TABLE `chats` DISABLE KEYS */;
INSERT INTO `chats` (`chat_id`, `room_id`, `post`, `post_user_id`)
VALUES
  ('0', '0', 'This is general chat.', 'system');
/*!40000 ALTER TABLE `chats` ENABLE KEYS */;
UNLOCK TABLES;

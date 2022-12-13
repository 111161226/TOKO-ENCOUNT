DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `user_id` VARCHAR(37) NOT NULL,
  `user_name` VARCHAR(30) NOT NULL,
  `password` VARCHAR(90) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_name` (`user_name`)
);

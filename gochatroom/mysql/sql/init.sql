CREATE DATABASE IF NOT EXISTS db_chatroom;

CREATE TABLE `users` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(20) NOT NULL,
    `nickname` varchar(20) NOT NULL,
    `avatar_id` bigint(20) DEFAULT NULL,
    `password_digest` longtext NOT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_users_username` (`username`)
) ENGINE = InnoDB AUTO_INCREMENT = 4 DEFAULT CHARSET = utf8mb4;

CREATE TABLE `messages` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) DEFAULT NULL,
    `to_user_id` bigint(20) DEFAULT NULL,
    `room_id` bigint(20) DEFAULT NULL,
    `content` longtext DEFAULT NULL,
    `image_url` longtext DEFAULT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_messages_user_id` (`user_id`),
    KEY `idx_messages_to_user_id` (`to_user_id`),
    KEY `idx_messages_room_id` (`room_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 661 DEFAULT CHARSET = utf8mb4;
DROP TABLE IF EXISTS `oauth2_client_scopes`;
CREATE TABLE `oauth2_client_scopes` (
  `id` varchar(255) NOT NULL,
  `key_code` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `oauth2_client_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `oauth2_clients`;
CREATE TABLE `oauth2_clients` (
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `secret_key` varchar(255) NOT NULL,
  `redirect_uri` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `oauth2_client_auth_codes`;
CREATE TABLE `oauth2_client_auth_codes` (
  `id` varchar(255) NOT NULL,
  `code` varchar(255) NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `oauth2_client_id` varchar(255) NOT NULL,
  `expired_at` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `oauth2_client_auth_code_scopes`;
CREATE TABLE `oauth2_client_auth_code_scopes` (
  `oauth2_client_auth_code_id` varchar(255) NOT NULL,
  `oauth2_client_scope_id` varchar(255) NOT NULL,
  PRIMARY KEY (`oauth2_client_auth_code_id`,`oauth2_client_scope_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

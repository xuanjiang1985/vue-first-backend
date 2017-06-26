CREATE TABLE `users` (
     `id` INT(11) unsigned NOT NULL AUTO_INCREMENT,
     `name` varchar(30) DEFAULT '匿名用户',
     `phone` varchar(12) DEFAULT NULL,
     `password` varchar(60) DEFAULT NULL,
     `header` varchar(255) DEFAULT '/public/images/header.jpg',
     `sex` tinyint(4) DEFAULT 0,
     `admin` tinyint(4) DEFAULT 0,
     `created_at` int(11) DEFAULT 0,
     `updated_at` int(11) DEFAULT 0,
     PRIMARY KEY (`id`),
     UNIQUE KEY `users_email_unique` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


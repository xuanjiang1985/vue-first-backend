CREATE TABLE `pictures` (
     `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	 `num` tinyint(4) DEFAULT 1,
	 `description` varchar(255) DEFAULT '',
     `first` varchar(255) DEFAULT '',
     `photos` varchar(1200) DEFAULT '',
     `created_at` int(11) DEFAULT 0,
     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
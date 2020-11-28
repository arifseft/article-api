CREATE DATABASE  IF NOT EXISTS `db_articles`;
USE `db_articles`;

CREATE TABLE IF NOT EXISTS `articles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
  `body` longtext COLLATE utf8_unicode_ci NOT NULL,
  `author` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO `articles` VALUES (1,'Title One','Body Article One','Brad Pitt','2017-05-18 13:50:19'),(2,'Title Two','Body Article Two','Johnny Depp','2017-05-18 13:50:19'),(3,'Title Three','Body Article Three','Tom Cruise','2017-05-18 13:50:19');

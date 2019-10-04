# /*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
# /*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
# /*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
# /*!40101 SET NAMES utf8 */;
# /*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
# /*!40103 SET TIME_ZONE='+00:00' */;
# /*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
# /*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
# /*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
# /*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `MUXIKSTACK` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `MUXIKSTACK`;

--
-- Table structure for table `tb_users`
--

DROP TABLE IF EXISTS `tb_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tb_users` (
  `sid` bigint(20) unsigned NOT NULL,
  `username` varchar(255) NOT NULL,
  `avatar` varchar(255) NOT NULL,
  `loginWay` int(8) unsigned NOT NULL,
  `loginCode` varchar(255) NOT NULL,
  PRIMARY KEY (`sid`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Insert mock data
--

CREATE TABLE "course_comment" (
  `course_id`         VARCHAR(50) NOT NULL,
  `course_name`       VARCHAR(20) NOT NULL,
  `star`              INT         NOT NULL,
  `attendance_check`  INT         NOT NULL,
  `exam_check`        INT         NOT NULL,
  `content`           text        NOT NULL,
  `is_annoymous`      TINYINT(1)  NOT NULL,
  `like_num`          INT         NOT NULL,
  `comment_num`       INT         NOT NULL,
  `tags`              VARCHAR(255) NOT NULL,
  `time`              VARCHAR(50) NOT NULL,
  `is_valid`          TINYINT(1)  NOT NULL,
)

CREATE TABLE "comment" (
  `time`              VARCHAR(20) NOT NULL,
  `content`           text        NOT NULL,
  `like_num`          INT         NOT NULL,
  `user_id`           INT         NOT NULL,
  `parent_id`         INT         NOT NULL,
  `is_root`           TINYINT(1)  NOT NULL,
  `comment_target_id` INT         NOT NULL,
  `subcomment_num`    INT         NOT NULL,
)

CREATE TABLE "course_comment_like" (
  `course_comment_id` INT NOT NULL,
  `user_id`           INT NOT NULL,
)

CREATE TABLE "comment_like" (
  `comment_id` INT NOT NULL,
  `user_id`    INT NOT NULL,
)

CREATE TABLE "class_table" (
  `user_id` INT  NOT NULL,
  `courses` text NOT NULL,
)

CREATE TABLE "tag" (
  `tag_name` VARCHAR(20) NOT NULL,
)


# LOCK TABLES `tb_users` WRITE;
# /*!40000 ALTER TABLE `tb_users` DISABLE KEYS */;
# INSERT INTO `tb_users` VALUES (0,'admin','$2a$10$veGcArz47VGj7l9xN7g2iuT9TF21jLI1YGXarGzvARNdnt4inC9PG','2018-05-27 16:25:33','2018-05-27 16:25:33',NULL);
# /*!40000 ALTER TABLE `tb_users` ENABLE KEYS */;
# UNLOCK TABLES;
# /*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;
#
# /*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
# /*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
# /*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
# /*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
# /*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
# /*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
# /*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-05-28  0:25:41

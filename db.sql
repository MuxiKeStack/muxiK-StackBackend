DROP DATABASE IF EXISTS `MUXIKSTACK`;

CREATE DATABASE `MUXIKSTACK`;

USE `MUXIKSTACK`;

CREATE TABLE `user` (
  `id`         INT              unsigned  NOT NULL AUTO_INCREMENT,
  `sid`        VARCHAR(10)      NOT NULL COMMENT   "学生学号",
  `username`   VARCHAR(25)      ,
  `avatar`     VARCHAR(255)     ,
  `is_blocked` TINYINT          NOT NULL DEFAULT 0,

  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_evaluation` (
  `id`                    INT unsigned NOT NULL AUTO_INCREMENT,
  `course_name`           VARCHAR(50)  NOT NULL,
  `rate`                  INT          NOT NULL DEFAULT 0,
  `attendance_check_type` INT          NOT NULL DEFAULT 0 COMMENT "考勤方式，经常点名/偶尔点名/签到点名，标识为 0/1/2",
  `exam_check_type`       INT          NOT NULL DEFAULT 0 COMMENT "考核方式，无考核/闭卷考试/开卷考试/论文考核，标识为 0/1/2/3",
  `content`               TEXT         NOT NULL           COMMENT "评课内容",
  `is_anonymous`          TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否匿名评课",
  `like_num`              INT          NOT NULL DEFAULT 0 COMMENT "点赞数",
  `comment_num`           INT          NOT NULL DEFAULT 0 COMMENT "一级评论数",
  `tags`                  VARCHAR(255) NOT NULL           COMMENT "标签id列表，逗号分隔",
  `time`                  VARCHAR(20)  NOT NULL           COMMENT "评课时间，时间戳",
  `is_valid`              TINYINT(1)            DEFAULT 1 COMMENT "是否有效，未被折叠",

  `course_id`             VARCHAR(50)  NOT NULL,
  `user_id`               INT          NOT NULl,

  PRIMARY KEY (`id`),
  KEY `course_id` (`course_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

-- CREATE TABLE `parent_comment` (
--   `id`                INT unsigned NOT NULL AUTO_INCREMENT,
--   `time`              VARCHAR(20)  NOT NULL           COMMENT "评课时间，时间戳",
--   `content`           TEXT         NOT NULl           COMMENT "评论内容",
--   `like_num`          INT          NOT NULL DEFAULT 0 COMMENT "点赞数",
--   `sub_comment_num`   INT          NOT NULL DEFAULT 0 COMMENT "子评论数",
--   `is_anonymous`      TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否匿名",
--   `is_valid`          TINYINT(1)            DEFAULT 1 COMMENT "是否有效，未被折叠",
--
--   `user_id`           INT          NOT NULL,
--   `evaluation_id`     INT          NOT NULL COMMENT "评课id",
--
--   PRIMARY KEY (`id`),
--   KEY `user_id` (`user_id`),
--   KEY `evaluation_id` (`evaluation_id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;
--
-- CREATE TABLE `sub_comment` (
--   `id`                INT unsigned NOT NULL AUTO_INCREMENT,
--   `time`              VARCHAR(20)  NOT NULL           COMMENT "评课时间，时间戳",
--   `content`           TEXT         NOT NULl           COMMENT "评论内容",
--   `like_num`          INT          NOT NULL DEFAULT 0 COMMENT "点赞数",
--   `is_anonymous`      TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否匿名",
--   `is_valid`          TINYINT(1)            DEFAULT 1 COMMENT "是否有效，未被折叠",
--
--   `user_id`           INT          NOT NULL,
--   `parent_id`         INT          NOT NULL,
--   `comment_target_id` INT          NOT NULL COMMENT "评论对象id",
--
--   PRIMARY KEY (`id`),
--   KEY `user_id` (`user_id`),
--   KEY `parent_id` (`parent_id`),
--   KEY `comment_target_id` (`comment_target_id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `comment` (
  `id`                INT unsigned NOT NULL AUTO_INCREMENT,
  `time`              VARCHAR(20)  NOT NULL           COMMENT "评课时间，时间戳",
  `content`           TEXT         NOT NULl           COMMENT "评论内容",
  `like_num`          INT          NOT NULL DEFAULT 0 COMMENT "点赞数",
  `sub_comment_num`   INT          NOT NULL DEFAULT 0 COMMENT "子评论数",
  `is_root`           TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否一级评论",
  `is_anonymous`      TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否匿名",
  `is_valid`          TINYINT(1)            DEFAULT 1 COMMENT "是否有效，未被折叠",

  `user_id`           INT          NOT NULL,
  `parent_id`         INT          NOT NULL,
  `comment_target_id` INT          NOT NULL COMMENT "评论对象id",

  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `parent_id` (`parent_id`),
  KEY `comment_target_id` (`comment_target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_evaluation_like` (
  `id`            INT unsigned NOT NULL AUTO_INCREMENT,
  `evaluation_id` INT          NOT NULL COMMENT "评课id",
  `user_id`       INT          NOT NULL,

  PRIMARY KEY (`id`),
  KEY `evaluation_id` (`evaluation_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `comment_like` (
  `id`         INT unsigned NOT NULL AUTO_INCREMENT,
  `comment_id` INT          NOT NULL COMMENT "评论id",
  `user_id`    INT          NOT NULL,

  PRIMARY KEY (`id`),
  KEY `comment_id` (`comment_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `class_table` (
  `id`      INT unsigned NOT NULL AUTO_INCREMENT,
  `user_id` INT          NOT NULL,
  `name`    VARCHAR(20)  NOT NULL DEFAULT "课表",
  `classes` TEXT         NOT NULL COMMENT "课堂 hash 列表，逗号分隔",

  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `tag` (
  `id`       INT unsigned NOT NULL AUTO_INCREMENT,
  `tag_name` VARCHAR(20)  NOT NULL,

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `report` (
  `id`            INT unsigned NOT NULL AUTO_INCREMENT,
  `evaluation_id` INT          NOT NULL,
  `user_id`       INT          NOT NULL,
  `pass`          TINYINT(1)   NOT NULL DEFAULT 1 COMMENT "举报审核是否通过",

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_list` (
  `id`        INT unsigned NOT NULL AUTO_INCREMENT,
  `user_id`   INT          NOT NULL,
  `course_id` VARCHAR(50)  NOT NULL,

  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `history_course` (
  `id`      INT unsigned NOT NULL AUTO_INCREMENT,
  `hash`    VARCHAR(50)  NOT NULL COMMENT "课程id + 教师名 hash 生成的唯一标识",
  `name`    VARCHAR(50)  NOT NULL,
  `teacher` VARCHAR(20)  NOT NULL,
  `type`    INT          NOT NULL COMMENT "课程类型（专业课，公共课）",
    `rate`           FLOAT        NOT NULL DEFAULT 0 COMMENT "课程评价星级",
  `stars_num`      INT          NOT NULL DEFAULT 0 COMMENT "参与评课人数",
  `credit`         FLOAT        NOT NULL DEFAULT 0 COMMENT "学分",

  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `using_course` (
  `id`             INT unsigned NOT NULL AUTO_INCREMENT,
  `hash`           VARCHAR(50)  NOT NULL           COMMENT "课程id + 教师名 hash 生成的唯一标识",
  `name`           VARCHAR(50)  NOT NULL,
  `teacher`        VARCHAR(20)  NOT NULL,
  `course_id`      INT          NOT NULL           COMMENT "课程号",
  `class_id`       INT          NOT NULL           COMMENT "教学班编号",
  `type`           INT          NOT NULL           COMMENT "通识必修，通识选修，通识核心，专业必修，专业选修分别为 0/1/2/3/4",
  `credit_type`    INT          NOT NULL           COMMENT "学分类别，文科理科艺术之类的，加索引（筛选条件）",
  `time1`          VARCHAR(20)  NOT NULL,
  `place1`         VARCHAR(20)  NOT NULL,
  `time2`          VARCHAR(20)  NOT NULL,
  `place2`         VARCHAR(20)  NOT NULL,
  `time3`          VARCHAR(20)  NOT NULL,
  `place3`         VARCHAR(20)  NOT NULL,
  `weeks`          VARCHAR(20)  NOT NULL,
  `region`         INT          NOT NULL COMMENT "上课地区，南湖，东区，西区。加索引（筛选条件）",

  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

DROP DATABASE IF EXISTS `muxikstack`;

CREATE DATABASE `muxikstack`;

USE `muxikstack`;

CREATE TABLE `user` (
  `id`         INT              unsigned  NOT NULL AUTO_INCREMENT,
  `sid`        VARCHAR(10)      NOT NULL COMMENT   "学生学号",
  `username`   VARCHAR(25)      ,
  `avatar`     VARCHAR(255)     ,
  `is_blocked` TINYINT          NOT NULL DEFAULT 0,

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `message` (
  `id`          INT UNSIGNED NOT NULL auto_increment,
  `pub_user_id` INT UNSIGNED NOT NULL DEFAULT 0,
  `sub_user_id` INT UNSIGNED NOT NULL DEFAULT 0,
  `is_like`     TINYINT(1),
  `is_read`     TINYINT(1)   NOT NULL  DEFAULT 0,
  `reply`       VARCHAR(255),
  `time`        VARCHAR(20)  NOT NULL,
  `course_info` VARCHAR(255) NOT NULL,

  PRIMARY KEY (`id`),
  KEY sub_user_id (`sub_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_evaluation` (
  `id`                    INT unsigned NOT NULL AUTO_INCREMENT,
  `course_name`           VARCHAR(50)  NOT NULL,
  `rate`                  FLOAT        NOT NULL DEFAULT 0,
  `attendance_check_type` INT          NOT NULL DEFAULT 0 COMMENT "考勤方式，经常点名/偶尔点名/签到点名，标识为 0/1/2",
  `exam_check_type`       INT          NOT NULL DEFAULT 0 COMMENT "考核方式，无考核/闭卷考试/开卷考试/论文考核，标识为 0/1/2/3",
  `content`               TEXT                            COMMENT "评课内容",
  `is_anonymous`          TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否匿名评课",
  `like_num`              INT          NOT NULL DEFAULT 0 COMMENT "点赞数",
  `comment_num`           INT          NOT NULL DEFAULT 0 COMMENT "一级评论数",
  `tags`                  VARCHAR(255)                    COMMENT "标签id列表，逗号分隔",
  `time`                  DATETIME     NOT NULL           COMMENT "评课时间",
  `is_valid`              TINYINT(1)            DEFAULT 1 COMMENT "是否有效，未被折叠",
  `deleted_at`            TIMESTAMP    NULL     DEFAULT NULL,

  `course_id`             VARCHAR(50)  NOT NULL,
  `user_id`               INT          NOT NULl,

  PRIMARY KEY (`id`),
  KEY `course_id` (`course_id`),
  KEY `user_id` (`user_id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `parent_comment` (
  `id`              VARCHAR(40) NOT NULL           COMMENT "uuid",
  `time`            DATETIME    NOT NULL           COMMENT "评课时间",
  `content`         TEXT                           COMMENT "评论内容",
  `sub_comment_num` INT         NOT NULL DEFAULT 0 COMMENT "子评论数",
  `is_anonymous`    TINYINT(1)  NOT NULL DEFAULT 0 COMMENT "是否匿名",
  `is_valid`        TINYINT(1)  NOT NULL DEFAULT 1 COMMENT "是否有效，未被折叠",
  `deleted_at`      TIMESTAMP   NULL     DEFAULT NULL,

  `user_id`         INT         NOT NULL,
  `evaluation_id`   INT         NOT NULL COMMENT "评课id",

  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `evaluation_id` (`evaluation_id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `sub_comment` (
  `id`           VARCHAR(40) NOT NULL           COMMENT "uuid",
  `time`         DATETIME    NOT NULL           COMMENT "评课时间",
  `content`      TEXT                           COMMENT "评论内容",
  `is_anonymous` TINYINT(1)  NOT NULL DEFAULT 0 COMMENT "是否匿名",
  `is_valid`     TINYINT(1)  NOT NULL DEFAULT 1 COMMENT "是否有效，未被折叠",
  `deleted_at`   TIMESTAMP   NULL     DEFAULT NULL,

  `parent_id`      VARCHAR(40) NOT NULL,
  `user_id`        INT         NOT NULL,
  `target_user_id` INT         NOT NULL COMMENT "评论的目标用户id",

  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `parent_id` (`parent_id`),
  KEY `target_user_id` (`target_user_id`),
  KEY `deleted_at` (`deleted_at`)
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
  `comment_id` VARCHAR(40)  NOT NULL COMMENT "评论id",
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

CREATE TABLE `tags` (
  `id`   INT unsigned NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(20)  NOT NULL,

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_tag` (
  `id`        INT unsigned NOT NULL AUTO_INCREMENT,
  `tag_id`    INT          NOT NULL,
  `course_id` VARCHAR(50)  NOT NULL,
  `num`       INT          NOT NULL,

  PRIMARY KEY (`id`),
  KEY `tag_id` (`tag_id`),
  KEY `course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `report` (
  `id`            INT unsigned NOT NULL AUTO_INCREMENT,
  `evaluation_id` INT          NOT NULL,
  `user_id`       INT          NOT NULL,
  `pass`          TINYINT(1)   NOT NULL DEFAULT 1 COMMENT "举报审核是否通过",

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_list` (
  `id`             INT unsigned NOT NULL AUTO_INCREMENT,
  `user_id`        INT          NOT NULL,
  `course_hash_id` VARCHAR(50)  NOT NULL,

  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_course_hash_id` (`course_hash_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `history_course` (
  `id`        INT unsigned NOT NULL AUTO_INCREMENT,
  `hash`      VARCHAR(50)  NOT NULL COMMENT "课程id + 教师名 hash 生成的唯一标识",
  `name`      VARCHAR(50)  NOT NULL,
  `teacher`   VARCHAR(20)  NOT NULL,
  `type`      INT          NOT NULL COMMENT "课程类型（专业课，公共课）",
  `rate`      FLOAT        NOT NULL DEFAULT 0 COMMENT "课程评价星级",
  `stars_num` INT          NOT NULL DEFAULT 0 COMMENT "参与评课人数",
  `credit`    FLOAT        NOT NULL DEFAULT 0 COMMENT "学分",

  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`),
  FULLTEXT KEY (`name`, `teacher`) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `using_course` (
  `id`             INT unsigned NOT NULL AUTO_INCREMENT,
  `hash`           VARCHAR(50)  NOT NULL           COMMENT "课程id + 教师名 hash 生成的唯一标识",
  `name`           VARCHAR(50)  NOT NULL DEFAULT "",
  `teacher`        VARCHAR(20)  NOT NULL DEFAULT "",
  `course_id`      VARCHAR(8)   NOT NULL           COMMENT "课程号",
  `class_id`       INT          NOT NULL           COMMENT "教学班编号",
  `type`           INT          NOT NULL           COMMENT "通识必修，通识选修，通识核心，专业必修，专业选修分别为 0/1/2/3/4",
  `time1`          VARCHAR(20)  NOT NULL DEFAULT "",
  `place1`         VARCHAR(20)  NOT NULL DEFAULT "",
  `time2`          VARCHAR(20)  NOT NULL DEFAULT "",
  `place2`         VARCHAR(20)  NOT NULL DEFAULT "",
  `time3`          VARCHAR(20)  NOT NULL DEFAULT "",
  `place3`         VARCHAR(20)  NOT NULL DEFAULT "",
  `weeks1`         VARCHAR(20)  NOT NULL DEFAULT "",
  `weeks2`         VARCHAR(20)  NOT NULL DEFAULT "",
  `weeks3`         VARCHAR(20)  NOT NULL DEFAULT "",
  `region`         INT          NOT NULL COMMENT "上课地区，1-南湖，2-东区，3-西区。加索引（筛选条件）",

  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`),
  FULLTEXT KEY (`name`, `course_id`, `teacher`) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

INSERT INTO `tags` (name) VALUES ("简单易学"), ("干货满满"), ("生动有趣"), ("作业量少"), ("老师温柔"), ("云课堂资料全");


-- mock data

INSERT INTO `history_course` (hash, name, teacher, type) VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', 0);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', 10, 3, '1-2#1', '7205', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('213f89eyguiguhy', '数据库原理', '喻莹', '98767654', 20, 3, '3-4#2', '9201', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('298yr3y9f8euibf', '数字逻辑', '赵甫哲', '34789865', 30, 3, '1-2#3', '9402', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('38yry8fe7guiwb3', '计算机组成原理', '李沛', '23345678', 40, 3, '10-11#2', '6221', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('298y38efub2ef32', '计算机网络', '王林平', '09898767', 50, 3, '5-6#5', '8205', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('28yy89dqube12d8', '面向对象程序设计', '胡珀', '34569876', 60, 3, '7-8#1', 'JKSYS-3', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('723fguib98y2e1h', 'Python程序设计', '胡珀', '23456987', 70, 3, '1-2#5', '9403', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('2e154de56gyubdq', '高级语言程序设计', '沈显军', '56982345', 80, 3, '3-4#5', '1122', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('0s9uighvg121efe', 'Java程序设计', '张连发', '09865423', 90, 3, '9-10#4', '9501', '2-17#0', 2);


-- mock data

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2016213456', '你管我叫啥呢', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2017213213', 'fucking...', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2018243789', 'Wow, FPXnb!', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2019211678', '大爱法学院', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2016221983', '赵凌云大傻逼', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2017908932', 'i华大牛逼', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2018923872', '当代恶臭网民', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2018214830', '随便呗', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2019526782', '孙笑川', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2017211712', '中华人民共和国湖北省武汉市洪山区当代恶臭网民孙笑川', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2017909876', 'GITHUB', '0');
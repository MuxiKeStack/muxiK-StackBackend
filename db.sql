-- DROP DATABASE IF EXISTS `muxikstack`;

CREATE DATABASE IF NOT EXISTS `muxikstack`;

USE `muxikstack`;

CREATE TABLE `user` (
  `id`         INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sid`        VARCHAR(10)  NOT NULL UNIQUE COMMENT "学生学号",
  `username`   VARCHAR(25)  ,
  `avatar`     VARCHAR(255) ,
  `is_blocked` TINYINT(1)   NOT NULL DEFAULT 0,
  `licence`    TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "成绩查看许可",

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `message` (
  `id`          INT UNSIGNED NOT NULL auto_increment,
  `pub_user_id` INT UNSIGNED NOT NULL DEFAULT 0,
  `sub_user_id` INT UNSIGNED NOT NULL DEFAULT 0,
  `kind`        TINYINT(1) UNSIGNED   NOT NULL DEFAULT 0  COMMENT "消息提醒的种类，0是点赞，1是评论，2是举报",
  `is_read`     TINYINT(1)   NOT NULL DEFAULT 0,
  `reply`       VARCHAR(255),
  `time`        VARCHAR(20)  NOT NULL,
  `course_id`   VARCHAR(255),
  `course_name` VARCHAR(255),
  `teacher`     VARCHAR(255),
  `evaluation_id` INT UNSIGNED,
  `content`     VARCHAR(255),
  `sid`         VARCHAR(255),
  `parent_comment_id`     VARCHAR(255),

  PRIMARY KEY (`id`),
  KEY sub_user_id (`sub_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_evaluation` (
  `id`                    INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `course_name`           VARCHAR(50)  NOT NULL,
  `rate`                  FLOAT        NOT NULL DEFAULT 0,
  `attendance_check_type` INT          NOT NULL DEFAULT 0 COMMENT "考勤方式，经常点名/偶尔点名/签到点名，标识为 1/2/3",
  `exam_check_type`       INT          NOT NULL DEFAULT 0 COMMENT "考核方式，无考核/闭卷考试/开卷考试/论文考核，标识为 1/2/3/4",
  `content`               TEXT                            COMMENT "评课内容",
  `is_anonymous`          TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否匿名评课",
  `like_num`              INT          NOT NULL DEFAULT 0 COMMENT "点赞数",
  `comment_num`           INT          NOT NULL DEFAULT 0 COMMENT "一级评论数",
  `tags`                  VARCHAR(255)                    COMMENT "标签id列表，逗号分隔",
  `time`                  DATETIME     NOT NULL           COMMENT "评课时间",
  `is_valid`              TINYINT(1)            DEFAULT 1 COMMENT "是否有效，未被折叠",
  `deleted_at`            TIMESTAMP    NULL     DEFAULT NULL,

  `course_id`             VARCHAR(50)  NOT NULL,
  `user_id`               INT UNSIGNED NOT NULL,

  PRIMARY KEY (`id`),
  KEY `course_id` (`course_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `parent_comment` (
  `id`              VARCHAR(40) NOT NULL           COMMENT "uuid",
  `time`            DATETIME    NOT NULL           COMMENT "评课时间",
  `content`         TEXT                           COMMENT "评论内容",
  `sub_comment_num` INT         NOT NULL DEFAULT 0 COMMENT "子评论数",
  `is_anonymous`    TINYINT(1)  NOT NULL DEFAULT 0 COMMENT "是否匿名",
  `is_valid`        TINYINT(1)  NOT NULL DEFAULT 1 COMMENT "是否有效，未被折叠",
  `deleted_at`      TIMESTAMP   NULL     DEFAULT NULL,

  `user_id`         INT UNSIGNED NOT NULL,
  `evaluation_id`   INT UNSIGNED NOT NULL COMMENT "评课id",

  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `evaluation_id` (`evaluation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `sub_comment` (
  `id`           VARCHAR(40) NOT NULL           COMMENT "uuid",
  `time`         DATETIME    NOT NULL           COMMENT "评课时间",
  `content`      TEXT                           COMMENT "评论内容",
  `is_anonymous` TINYINT(1)  NOT NULL DEFAULT 0 COMMENT "是否匿名",
  `is_valid`     TINYINT(1)  NOT NULL DEFAULT 1 COMMENT "是否有效，未被折叠",
  `deleted_at`   TIMESTAMP   NULL     DEFAULT NULL,

  `parent_id`      VARCHAR(40) NOT NULL,
  `user_id`        INT UNSIGNED  NOT NULL,
  `target_user_id` INT UNSIGNED  NOT NULL COMMENT "评论的目标用户id",

  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `parent_id` (`parent_id`),
  KEY `target_user_id` (`target_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_evaluation_like` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `evaluation_id` INT UNSIGNED NOT NULL COMMENT "评课id",
  `user_id`       INT UNSIGNED NOT NULL,

  PRIMARY KEY (`id`),
  KEY `evaluation_id` (`evaluation_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `comment_like` (
  `id`         INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `comment_id` VARCHAR(40)  NOT NULL COMMENT "评论id",
  `user_id`    INT UNSIGNED NOT NULL,

  PRIMARY KEY (`id`),
  KEY `comment_id` (`comment_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `class_table` (
  `id`      INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` INT UNSIGNED NOT NULL,
  `name`    VARCHAR(20)  NOT NULL DEFAULT "新课表",
  `classes` TEXT         NOT NULL COMMENT "课堂 hash id 列表，逗号分隔",

  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `tags` (
  `id`   INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(20)  NOT NULL,

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_tag` (
  `id`        INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tag_id`    INT UNSIGNED NOT NULL,
  `course_id` VARCHAR(50)  NOT NULL,
  `num`       INT          NOT NULL,

  PRIMARY KEY (`id`),
  KEY `tag_id` (`tag_id`),
  KEY `course_id` (`course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `report` (
  `id`            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `evaluation_id` INT UNSIGNED NOT NULL,
  `user_id`       INT UNSIGNED NOT NULL,
  `pass`          TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "举报审核是否通过",
  `reason`        VARCHAR(200) NOT NULL,

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `course_list` (
  `id`             INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`        INT UNSIGNED NOT NULL,
  `course_hash_id` VARCHAR(50)  NOT NULL,

  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_course_hash_id` (`course_hash_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `history_course` (
  `id`        INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `hash`      VARCHAR(50)  NOT NULL COMMENT "课程id + 教师名 hash 生成的唯一标识",
  `name`      VARCHAR(50)  NOT NULL,
  `teacher`   VARCHAR(50)  NOT NULL,
  `course_id`      VARCHAR(9)   NOT NULL            COMMENT "课程号",
  `type`      INT          NOT NULL COMMENT "课程类型（根据学校提供的特定位进行判定）0-通必,1-专必,2-专选,3-通选,5-通核",
  `rate`      FLOAT        NOT NULL DEFAULT 0 COMMENT "课程评价星级",
  `stars_num` INT          NOT NULL DEFAULT 0 COMMENT "参与评课人数",
  `total_grade` FLOAT      NOT NULL DEFAULT 0 COMMENT "总成绩",
  `usual_grade` FLOAT      NOT NULL DEFAULT 0 COMMENT "平时成绩",
  `grade_sample_size` INT  NOT NULL DEFAULT 0 COMMENT "成绩样本人数",
  `grade_section_1`   INT  NOT NULL DEFAULT 0 COMMENT "成绩区间1,85以上",
  `grade_section_2`   INT  NOT NULL DEFAULT 0 COMMENT "成绩区间2,70-85",
  `grade_section_3`   INT  NOT NULL DEFAULT 0 COMMENT "成绩区间3,70以下",

  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`),
  FULLTEXT KEY (`name`, `teacher`) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `using_course` (
  `id`             INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `hash`           VARCHAR(50)  NOT NULL            COMMENT "课程id + 教师名 hash 生成的唯一标识",
  `name`           VARCHAR(50)  NOT NULL DEFAULT "",
  `academy`        VARCHAR(25)  NOT NULL DEFAULT "" COMMENT "开课学院",
  `teacher`        VARCHAR(200) NOT NULL DEFAULT "",
  `credit`         FLOAT        NOT NULL DEFAULT 0  COMMENT "学分",
  `course_id`      VARCHAR(9)   NOT NULL            COMMENT "课程号",
  `class_id`       VARCHAR(50)  NOT NULL            COMMENT "教学班编号",
  `type`           INT          NOT NULL            COMMENT "通识必修，通识选修，通识核心，专业必修，专业选修分别为 0/1/2/3/4",
  `time1`          VARCHAR(30)  NOT NULL DEFAULT "",
  `place1`         VARCHAR(30)  NOT NULL DEFAULT "",
  `time2`          VARCHAR(30)  NOT NULL DEFAULT "",
  `place2`         VARCHAR(30)  NOT NULL DEFAULT "",
  `time3`          VARCHAR(30)  NOT NULL DEFAULT "",
  `place3`         VARCHAR(30)  NOT NULL DEFAULT "",
  `weeks1`         VARCHAR(30)  NOT NULL DEFAULT "",
  `weeks2`         VARCHAR(30)  NOT NULL DEFAULT "",
  `weeks3`         VARCHAR(30)  NOT NULL DEFAULT "",
  `region`         INT          NOT NULL COMMENT "上课地区，1-南湖，2-东区，3-西区。加索引（筛选条件）",

  PRIMARY KEY (`id`),
  KEY `hash` (`hash`),
  FULLTEXT KEY (`name`, `course_id`, `teacher`) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `grade` (
  `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `course_name` VARCHAR(50)  NOT NULL DEFAULT "",
  `total_score` FLOAT        NOT NULL DEFAULT 0 COMMENT "总成绩",
  `usual_score` FLOAT        NOT NULL DEFAULT 0 COMMENT "平时成绩",
  `final_score` FLOAT        NOT NULL DEFAULT 0 COMMENT "期末成绩",
  `has_added`   TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否已加入统计样本中",

  `user_id`        INT UNSIGNED NOT NULL,
  `course_hash_id` VARCHAR(50)  NOT NULL,

  PRIMARY KEY (`id`),
  KEY `idx_user_hash` (`user_id`, `course_hash_id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

-- CREATE TABLE `self_course` (
--   `id`      INT UNSIGNED      NOT NULL AUTO_INCREMENT,
--   `user_id` INT UNSIGNED      NOT NULL,
--   `course_hash_id` VARCHAR(50) NOT NULL,
--   `year_term`      CHAR(5)   NOT NULL DEFAULT '' COMMENT "学年和学期，'20181'->2018学年第一学期",
-- --   `num`     SMALLINT UNSIGNED NOT NULL COMMENT "课程数",
-- --   `courses` TEXT              NOT NULL COMMENT "课程 hash id 列表，逗号分隔",
--
--   PRIMARY KEY (`id`),
--   KEY `idx_user_id` (`user_id`),
--   KEY `idx_year_term` (`year_term`),
--   KEY `idx_hash_id` (`course_hash_id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

INSERT INTO `tags` (name) VALUES ("课程简单易学"), ("课程干货满满"), ("老师严谨负责"),
("老师温柔随和"), ("老师风趣幽默"), ("平时作业少"), ("期末划重点"), ("云课堂资料全");

-- 系统用户 以及 匿名用户信息 匿名用户信息都为空

INSERT INTO `user` 
  (`id`, `sid`, `username`)
VALUES 
  (1, "0", "系统提醒"),
  (2, "1", null);


-- mock data
/*
INSERT INTO `history_course` (hash, name, teacher, type) VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', 0);
INSERT INTO `history_course` (hash, name, teacher, type) VALUES ('213f89eyguiguhy', '数据库原理', '喻莹', 0);
INSERT INTO `history_course` (hash, name, teacher, type) VALUES ('2e154de56gyubdq', '高级语言程序设计', '沈显军', 1);
INSERT INTO `history_course` (hash, name, teacher, type) VALUES ('0s9uighvg121efe', 'Java程序设计', '张连发', 1);
INSERT INTO `history_course` (hash, name, teacher, type) VALUES ('28yy89dqube12d8', '面向对象程序设计', '胡珀', 2);
INSERT INTO `history_course` (hash, name, teacher, type) VALUES ('723fguib98y2e1h', 'Python程序设计', '胡珀', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "1", 3, '1-2#1', '7205', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "2", 3, '3-4#1', '5-6#3', '7207', '7105', '2-17#0', '2-15#2', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2,time3, place1, place2, place3, weeks1, weeks2, weeks3, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "3", 3, '9-10#2', '5-6#4', '1-2#5', '7207', '7105', '7201', '2-17#0', '2-17#0', '2-17#2', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "4", 3, '3-4#1', '5-6#3', '7207', '7105', '2-17#0', '2-15#2', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "5", 3, '3-4#1', '5-6#1', '7207', '7105', '2-17#1', '2-15#2', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "6", 3, '3-4#1', '5-6#3', '7207', '7105', '2-17#0', '2-15#1', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "7", 3, '3-4#1', '5-6#3', '7207', '7105', '2-17#1', '2-15#2', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "12", 3, '3-4#1', '5-6#3', '7207', '7105', '2-17#0', '2-15#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "42", 3, '3-4#1', '5-6#3', '7207', '7105', '2-17#1', '2-15#2', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2,time3, place1, place2, place3, weeks1, weeks2, weeks3, region)
VALUES ('112d34testsvggase', '高等数学A', '宋冰玉', '45677654', "55", 3, '1-2#4', '5-6#4', '9-10#4', '7207', '7105', '7201', '2-17#0', '2-17#0', '2-15#2', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, time3, place1, place2, place3 ,weeks1, weeks2, weeks3, region)
VALUES ('213f89eyguiguhy', '数据库原理', '喻莹', '98767654', "1", 3, '3-4#2', '7-8#4', '1-2#5', '9201', '8204', '7309', '2-17#0', '5-20#1','2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('213f89eyguiguhy', '数据库原理', '喻莹', '98767654', "2", 3, '1-2#3', '5-6#1', '9201', '8412','2-17#0', '5-18#1',2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('213f89eyguiguhy', '数据库原理', '喻莹', '98767654', "3", 3, '7-8#4', '9201', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('298yr3y9f8euibf', '数字逻辑', '赵甫哲', '34789865', "1", 3, '1-2#3', '9402', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('38yry8fe7guiwb3', '计算机组成原理', '李沛', '23345678', "1", 3, '10-11#2', '6221', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('298y38efub2ef32', '计算机网络', '王林平', '09898767', "1", 3, '5-6#5', '8205', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('28yy89dqube12d8', '面向对象程序设计', '胡珀', '34569876', "1", 3, '7-8#1', 'JKSYS-3', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('723fguib98y2e1h', 'Python程序设计', '胡珀', '23456987', "1", 3, '1-2#5', '9403', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('2e154de56gyubdq', '高级语言程序设计', '沈显军', '56982345', "1", 3, '3-4#5', '1122', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, place1, weeks1, region)
VALUES ('0s9uighvg121efe', 'Java程序设计', '张连发', '09865423', "1", 3, '9-10#4', '9501', '2-17#0', 2);

INSERT INTO `using_course` (hash, name, teacher, course_id, class_id, type, time1, time2, place1, place2, weeks1, weeks2, region)
VALUES ('asdf1232314123dasdf', '概率统计A', '涂佳娟', '123712322', "12", 2, '3-4#1', '5-6#3', '7207', '7105', '2-17#0', '2-15#2', 2);

INSERT INTO `course_list` (user_id, course_hash_id) values (1, '0s9uighvg121efe'), (1, '112d34testsvggase');
INSERT INTO `course_evaluation_like` (user_id, evaluation_id) values (8, 1), (8, 2);

-- mock data

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2018214830', '随便呗', '0');

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
VALUES ('2019526782', '孙笑川', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2017211712', '中华人民共和国湖北省武汉市洪山区当代恶臭网民孙笑川', '0');

INSERT INTO `user` (`sid`, `username`, `is_blocked`)
VALUES ('2017909876', 'GITHUB', '0');

-- mock data

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('1', '98767654', '2');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('2', '98767654', '5');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('3', '98767654', '10');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('4', '98767654', '7');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('5', '98767654', '1');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('6', '98767654', '4');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('1', '45677654', '20');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('2', '45677654', '5');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('3', '45677654', '1');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('4', '45677654', '2');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('5', '45677654', '10');

INSERT INTO `course_tag` (`tag_id`, `course_id`, `num`)
VALUES ('6', '45677654', '6');
*/

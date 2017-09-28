CREATE TABLE `token` (
  `id`           BIGINT(20)          NOT NULL AUTO_INCREMENT,
  `token_type`   TINYINT(3) UNSIGNED NOT NULL DEFAULT '0',
  `token_symbol` VARCHAR(255)        NOT NULL DEFAULT '',
  `token_name`   VARCHAR(255)        NOT NULL DEFAULT '',
  `token_intro`  VARCHAR(255)        NOT NULL DEFAULT '',
  `token_fee`    DOUBLE              NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `token_type` (`token_type`),
  UNIQUE KEY `token_symbol` (`token_symbol`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 3
  DEFAULT CHARSET = utf8;

CREATE TABLE `token_record` (
  `id`             BIGINT(20)          NOT NULL AUTO_INCREMENT,
  `record_time`    DATETIME            NOT NULL,
  `record_type`    TINYINT(3) UNSIGNED NOT NULL DEFAULT '0',
  `token_id`       BIGINT(20)          NOT NULL,
  `transaction_id` VARCHAR(255)        NOT NULL DEFAULT '',
  `user_id`        BIGINT(20)          NOT NULL,
  `record_status`  BIGINT(20)          NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `transaction_id` (`transaction_id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 16
  DEFAULT CHARSET = utf8;
CREATE TABLE `user` (
  `id`            BIGINT(20)   NOT NULL AUTO_INCREMENT,
  `user_name`     VARCHAR(255) NOT NULL DEFAULT '',
  `user_password` VARCHAR(255) NOT NULL DEFAULT '',
  `mail_address`  VARCHAR(255) NOT NULL DEFAULT '',
  `mail_code`     VARCHAR(255) NOT NULL DEFAULT '',
  `phone_no`      VARCHAR(11)  NOT NULL DEFAULT '',
  `country_code`  VARCHAR(255) NOT NULL DEFAULT '',
  `create_time`   DATETIME     NOT NULL,
  `user_portrait` VARCHAR(255) NOT NULL DEFAULT '',
  `user_intro`    VARCHAR(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_name` (`user_name`),
  UNIQUE KEY `phone_no` (`phone_no`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 3
  DEFAULT CHARSET = utf8;
CREATE TABLE `user_token` (
  `id`            BIGINT(20)   NOT NULL AUTO_INCREMENT,
  `user_id`       BIGINT(20)   NOT NULL,
  `token_id`      BIGINT(20)   NOT NULL,
  `token_address` VARCHAR(255) NOT NULL DEFAULT '',
  `private_key`   VARCHAR(255) NOT NULL DEFAULT '',
  `token_balance` DOUBLE       NOT NULL DEFAULT '0',
  `token_extra`   VARCHAR(255) NOT NULL DEFAULT '',
  `lock_balance`  DOUBLE       NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8;
CREATE TABLE `withdrawal` (
  `id`       BIGINT(20)   NOT NULL AUTO_INCREMENT,
  `user_id`  BIGINT(20)   NOT NULL,
  `address`  VARCHAR(255) NOT NULL DEFAULT '0',
  `tag`      VARCHAR(255) NOT NULL DEFAULT '',
  `token_id` BIGINT(20)   NOT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8;

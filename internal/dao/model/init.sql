CREATE TABLE `traces`
(
    `id`          bigint UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '主键',
    `uuid`        char(36)                       NOT NULL COMMENT 'uuid',
    `type`        tinyint                        NOT NULL COMMENT '类型(0:注册 1:激活 2:登录)',
    `ip`          char(15)                       NOT NULL COMMENT 'ip',
    `create_time` timestamp                      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `uuid-type` (`uuid`, `type`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `users`
(
    `id`          bigint UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '主键',
    `uuid`        char(36)                       NOT NULL COMMENT 'uuid',
    `email`       varchar(320)                   NOT NULL COMMENT '邮箱',
    `password`    char(32)                       NOT NULL COMMENT '密码',
    `salt`        char(4)                        NOT NULL COMMENT 'salt',
    `status`      tinyint                        NOT NULL COMMENT '状态(0:未验证 1:正常)',
    `verify_code` char(20)                       NOT NULL COMMENT '验证码(激活)',
    `create_time` timestamp                      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_time` timestamp                      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `email-password-status-uuid` (`email`, `password`, `status`, `uuid`),
    KEY `uuid-verify_code` (`uuid`, `verify_code`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `apps`
(
    `id`          bigint UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '主键',
    `name`        varchar(256)                   NOT NULL COMMENT '应用名称',
    `app_id`      char(36)                       NOT NULL COMMENT 'appId',
    `status`      tinyint                        NOT NULL COMMENT '状态(0:锁定 1:正常)',
    `redirect`    varchar(256)                   NOT NULL COMMENT '跳转链接',
    `create_time` timestamp                      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    Primary Key (`id`),
    KEY `app_id-status` (`app_id`, `status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

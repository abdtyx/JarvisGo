CREATE TABLE `blacklists` (
    `pk` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `id` INT UNSIGNED NOT NULL,
    `type` VARCHAR(20) NOT NULL,
    `comment` VARCHAR(20) NOT NULL,
    `created_at` DATETIME(3) NOT NULL COMMENT '创建日期',
    `updated_at` DATETIME(3) NOT NULL COMMENT '更新日期',
    `deleted_at` DATETIME(3) NULL COMMENT '删除日期',
    PRIMARY KEY (`pk`),
    INDEX `identifier`(`id`) USING BTREE,
    INDEX `type`(`type`) USING BTREE
);

CREATE TABLE `jeminders` (
    `id` INT UNSIGNED NOT NULL,
    `created_at` DATETIME(3) NOT NULL COMMENT '创建日期',
    `updated_at` DATETIME(3) NOT NULL COMMENT '更新日期',
    `deleted_at` DATETIME(3) NULL COMMENT '删除日期',
    PRIMARY KEY (`id`)
);

CREATE TABLE `blacklist` (
    `pk` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `id` INT UNSIGNED NOT NULL,
    `type` VARCHAR(20) NOT NULL,
    PRIMARY KEY (`pk`),
    INDEX `identifier`(`id`) USING BTREE,
    INDEX `type`(`type`) USING BTREE
);

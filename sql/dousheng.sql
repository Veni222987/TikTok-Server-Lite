/*
 Navicat Premium Data Transfer

 Source Server         : dousheng
 Source Server Type    : MySQL
 Source Server Version : 80034
 Source Host           : 47.113.149.158:3306
 Source Schema         : dousheng

 Target Server Type    : MySQL
 Target Server Version : 80034
 File Encoding         : 65001

 Date: 21/08/2023 16:37:07
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for account
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account`  (
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`username`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` bigint NULL DEFAULT NULL,
  `video_id` bigint NULL DEFAULT NULL,
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_date` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for follow
-- ----------------------------
DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id_a` bigint NOT NULL COMMENT '关注者',
  `user_id_b` bigint NOT NULL COMMENT '被关注者',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 33 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for like
-- ----------------------------
DROP TABLE IF EXISTS `like`;
CREATE TABLE `like`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` bigint NULL DEFAULT NULL,
  `video_id` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 39 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint NOT NULL COMMENT '用户id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名称',
  `follow_count` int NULL DEFAULT NULL COMMENT '关注总数',
  `follower_count` int NULL DEFAULT NULL COMMENT '粉丝总数',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'https://img0.baidu.com/it/u=3798218377,2567059835&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=501' COMMENT '用户头像',
  `background_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'https://up.enterdesk.com/edpic/10/12/ee/1012eed2819b19b55872a6f6493c00e9.jpg' COMMENT '用户个人页顶部大图',
  `signature` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '这个人很懒，没有留下任何信息。' COMMENT '个人简介',
  `total_favorited` int NULL DEFAULT NULL COMMENT '获赞数量',
  `work_count` int NULL DEFAULT 0 COMMENT '作品数',
  `favorite_count` int NULL DEFAULT NULL COMMENT '喜欢数',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`  (
  `id` bigint NOT NULL,
  `author_id` bigint NULL DEFAULT NULL COMMENT '作者id',
  `play_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '视频url',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '封面url',
  `favorite_count` int NULL DEFAULT NULL COMMENT '点赞数量',
  `comment_count` int NULL DEFAULT NULL COMMENT '评论数量',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '视频标题',
  `time` bigint NULL DEFAULT NULL COMMENT '视频发布时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;

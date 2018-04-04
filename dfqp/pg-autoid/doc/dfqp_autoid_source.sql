-- phpMyAdmin SQL Dump
-- version 4.0.10.11
-- http://www.phpmyadmin.net
--
-- 主机: 127.0.0.1:3388
-- 生成日期: 2018-03-06 10:48:58
-- 服务器版本: 5.6.16-log
-- PHP 版本: 5.5.26

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- 数据库: `dfqp_common`
--

-- --------------------------------------------------------

--
-- 表的结构 `dfqp_autoid_source`
--

CREATE TABLE IF NOT EXISTS `dfqp_autoid_source` (
  `btag` varchar(128) NOT NULL DEFAULT '' COMMENT 'id生成器类型',
  `max_id` bigint(20) NOT NULL DEFAULT '1' COMMENT '当前最大id',
  `step` int(11) NOT NULL COMMENT '每次更新时的步伐值',
  `des` varchar(256) DEFAULT NULL COMMENT '类型描述',
  `update_time` int(11) NOT NULL COMMENT '最近一次更新时间',
  PRIMARY KEY (`btag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `dfqp_autoid_source`
--

INSERT INTO `dfqp_autoid_source` (`btag`, `max_id`, `step`, `des`, `update_time`) VALUES
('user', 160, 10, '用户自增表', 1520304494);

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

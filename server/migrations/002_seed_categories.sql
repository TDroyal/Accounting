-- 002_seed_categories.sql
-- 系统预置分类树（user_id=0），用户首次注册时复制一份到该用户名下。
-- 对应 docs/01-需求报告 §4.1.3。
-- 注意：MySQL 不允许 INSERT/UPDATE 目标表直接 SELECT 自身（报错 1093），
-- 这里先插入一级分类，再用会话变量记录其 id，最后插入二级分类。

-- 一级分类
INSERT INTO `categories` (`user_id`, `parent_id`, `name`, `type`, `sort`, `status`, `icon`) VALUES
(0, 0, '交通', 0, 1, 1, 'traffic'),
(0, 0, '话费', 0, 2, 1, 'phone'),
(0, 0, '转账', 1, 3, 1, 'transfer'),
(0, 0, '外卖', 0, 4, 1, 'takeout'),
(0, 0, '吃饭', 0, 5, 1, 'meal'),
(0, 0, '房租', 0, 6, 1, 'rent'),
(0, 0, '购物', 0, 7, 1, 'shopping'),
(0, 0, '娱乐', 0, 8, 1, 'entertainment'),
(0, 0, '其他', 0, 9, 1, 'other');

-- 用会话变量记录各一级分类 id（用户变量赋值不触发 1093 限制）
SET @traffic   = (SELECT id FROM `categories` WHERE name='交通'   AND parent_id=0 AND user_id=0);
SET @phone     = (SELECT id FROM `categories` WHERE name='话费'   AND parent_id=0 AND user_id=0);
SET @transfer  = (SELECT id FROM `categories` WHERE name='转账'   AND parent_id=0 AND user_id=0);
SET @takeout   = (SELECT id FROM `categories` WHERE name='外卖'   AND parent_id=0 AND user_id=0);
SET @meal      = (SELECT id FROM `categories` WHERE name='吃饭'   AND parent_id=0 AND user_id=0);
SET @rent      = (SELECT id FROM `categories` WHERE name='房租'   AND parent_id=0 AND user_id=0);

-- 二级分类
INSERT INTO `categories` (`user_id`, `parent_id`, `name`, `type`, `sort`, `status`) VALUES
(0, @traffic,  '地铁',   0, 1, 1),
(0, @traffic,  '公交',   0, 2, 1),
(0, @traffic,  '打车',   0, 3, 1),
(0, @traffic,  '加油',   0, 4, 1),
(0, @traffic,  '停车',   0, 5, 1),
(0, @phone,    '手机',   0, 1, 1),
(0, @phone,    '宽带',   0, 2, 1),
(0, @transfer, '还款',   1, 1, 1),
(0, @transfer, '红包',   1, 2, 1),
(0, @transfer, 'AA',     1, 3, 1),
(0, @takeout,  '美团',   0, 1, 1),
(0, @takeout,  '饿了么', 0, 2, 1),
(0, @meal,     '早餐',   0, 1, 1),
(0, @meal,     '午餐',   0, 2, 1),
(0, @meal,     '晚餐',   0, 3, 1),
(0, @meal,     '聚餐',   0, 4, 1),
(0, @rent,     '租金',   0, 1, 1),
(0, @rent,     '水电',   0, 2, 1),
(0, @rent,     '物业',   0, 3, 1);

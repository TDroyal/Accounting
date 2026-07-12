-- 002_seed_categories.sql
-- 系统预置分类树（user_id=0），用户首次注册时复制一份到该用户名下。
-- 对应 docs/01-需求报告 §4.1.3。

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

-- 二级分类（parent_id 引用上面一级的自增 id）
INSERT INTO `categories` (`user_id`, `parent_id`, `name`, `type`, `sort`, `status`) VALUES
(0, (SELECT id FROM categories WHERE name='交通' AND parent_id=0), '地铁', 0, 1, 1),
(0, (SELECT id FROM categories WHERE name='交通' AND parent_id=0), '公交', 0, 2, 1),
(0, (SELECT id FROM categories WHERE name='交通' AND parent_id=0), '打车', 0, 3, 1),
(0, (SELECT id FROM categories WHERE name='交通' AND parent_id=0), '加油', 0, 4, 1),
(0, (SELECT id FROM categories WHERE name='交通' AND parent_id=0), '停车', 0, 5, 1),
(0, (SELECT id FROM categories WHERE name='话费' AND parent_id=0), '手机', 0, 1, 1),
(0, (SELECT id FROM categories WHERE name='话费' AND parent_id=0), '宽带', 0, 2, 1),
(0, (SELECT id FROM categories WHERE name='转账' AND parent_id=0), '还款', 1, 1, 1),
(0, (SELECT id FROM categories WHERE name='转账' AND parent_id=0), '红包', 1, 2, 1),
(0, (SELECT id FROM categories WHERE name='转账' AND parent_id=0), 'AA', 1, 3, 1),
(0, (SELECT id FROM categories WHERE name='外卖' AND parent_id=0), '美团', 0, 1, 1),
(0, (SELECT id FROM categories WHERE name='外卖' AND parent_id=0), '饿了么', 0, 2, 1),
(0, (SELECT id FROM categories WHERE name='吃饭' AND parent_id=0), '早餐', 0, 1, 1),
(0, (SELECT id FROM categories WHERE name='吃饭' AND parent_id=0), '午餐', 0, 2, 1),
(0, (SELECT id FROM categories WHERE name='吃饭' AND parent_id=0), '晚餐', 0, 3, 1),
(0, (SELECT id FROM categories WHERE name='吃饭' AND parent_id=0), '聚餐', 0, 4, 1),
(0, (SELECT id FROM categories WHERE name='房租' AND parent_id=0), '租金', 0, 1, 1),
(0, (SELECT id FROM categories WHERE name='房租' AND parent_id=0), '水电', 0, 2, 1),
(0, (SELECT id FROM categories WHERE name='房租' AND parent_id=0), '物业', 0, 3, 1);

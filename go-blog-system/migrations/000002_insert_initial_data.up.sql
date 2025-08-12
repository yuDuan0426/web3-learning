-- 插入用户数据
INSERT INTO users (username, email, password, bio) VALUES
('admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '系统管理员'),
('user1', 'user1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '普通用户1'),
('user2', 'user2@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '普通用户2');

-- 插入标签数据
INSERT INTO tags (name) VALUES
('技术'),
('生活'),
('旅行'),
('美食'),
('编程');

-- 插入文章数据
INSERT INTO posts (title, content, user_id, status) VALUES
('Go语言入门指南', '这是一篇关于Go语言入门的文章，介绍了Go语言的基础语法和特性...', 1, 'published'),
('我的旅行日记', '上周末我去了杭州西湖，那里的景色真美...', 2, 'published'),
('如何做一道美味的红烧肉', '准备食材：五花肉500g，酱油，糖，葱姜蒜...
烹饪步骤：
1. 将五花肉切成块
2. 冷水下锅，焯水去血水
3. 锅中放油，放入葱姜蒜爆香
4. 放入肉块翻炒
5. 加入酱油、糖等调料
6. 加水没过肉块
7. 大火烧开后转小火慢炖1小时
8. 收汁即可出锅', 3, 'published'),
('Go并发编程实践', '本文将介绍Go语言中的goroutine和channel，以及如何利用它们进行并发编程...', 1, 'draft');

-- 插入文章标签关联数据
INSERT INTO post_tags (post_id, tag_id) VALUES
(1, 1), (1, 5),
(2, 2), (2, 3),
(3, 2), (3, 4),
(4, 1), (4, 5);

-- 插入评论数据
INSERT INTO comments (content, user_id, post_id) VALUES
('非常好的入门文章，谢谢分享！', 2, 1),
('我也想去西湖，请问有什么推荐的景点吗？', 3, 2),
('我按照这个方法做了，味道很赞！', 1, 3),
('期待更多Go相关的文章', 3, 1);

-- 插入回复评论
INSERT INTO comments (content, user_id, post_id, parent_id) VALUES
('谢谢支持，后续会更新更多内容', 1, 1, 1),
('可以去断桥、雷峰塔、三潭印月等地方', 2, 2, 2);
-- 清空所有表中的数据（按照依赖关系的反序）
DELETE FROM comments;
DELETE FROM post_tags;
DELETE FROM posts;
DELETE FROM tags;
DELETE FROM users;
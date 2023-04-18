# 设计取舍

1. 因为 MySQL 的语法也有一定复杂性，这里无需实现所有的 MySQL 语句，覆盖现有的大部分语句即可。
1. 为了复用已有代码，可以用本库来生成 SQL 语句，再用 sqlx 做查询和数据绑定等操作。
1. 在不改动现有 struct 的前提下，需要再创建一个 struct 与 table 进行绑定。由于有别名等存在，需要允许用户修改，这里不建议使用 `go generate` 生成，可以实现一个工具来转换创建表结构的 sql 文件到 go 文件。

# 需求

1. 支持 MySQL 8.x，无需向下兼容
1. 表名、字段名和别名等在和关键字相同时自动转义（考虑到 MySQL 关键字太多，改成全部转义）
1. 支持 `JOIN`
1. 支持 `COUNT()`、`GROUP_CONCAT()` 等函数
1. 支持 `GROUP BY`
1. 支持 `ORDER BY x ASC/DESC`
1. `WHERE` 条件支持 `IN (?)`
1. `WHERE` 条件支持 `OR`、`AND` 和括号
1. 支持 `FOR UPDATE/SHARE`，无需后接 `OF table`、`NOWAIT` 和 `SKIP LOCKED` 等
1. 支持这种不带别名、`HAVING`、`ANY` 等修饰的子查询：`SELECT * FROM permission WHERE role_id IN (SELECT user_role FROM security_domain_user WHERE ...)`
1. `INSERT INTO table (a, buf, c) VALUES (?, ?, ?)` 能自动匹配 `?` 数量
1. 支持 `INSERT IGNORE INTO`
1. 支持 `INSERT INTO ... SELECT ...`
1. 支持 `UPDATE ... SET a=a+?`
1. 表名可作为占位符，例如 `SELECT 1 FROM %s WHERE id=? FOR UPDATE`
1. 支持用 `NamedExec` 来批量插入，例如 `INSERT INTO table (a, buf, c) VALUES (:a, :buf, :c)`
1. 支持同时查询结果和 count
1. 缓存和预编译 SQL

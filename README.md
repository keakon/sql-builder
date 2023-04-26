# 设计取舍

1. 因为 MySQL 的语法也有一定复杂性，这里无需实现所有的 MySQL 语句，覆盖现有的大部分语句即可。
1. 不做过多检查，允许用户构造错误的 SQL。
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
1. 支持同时查询结果和 count（暂未实现）
1. 缓存和预编译 SQL（需搭配 sqlx，实测对于简单的语句，预编译后提升大概 3%，作用不大）

# 部分不支持的特性

1. `UNION` 语句
2. 复杂的子查询

# 使用

## 定义表结构

```go
type UserTable struct {
	Table `db:"user"` // 表名
	ID    Column `db:"id"`
	Name  Column // 未写 tag 时，将取小写形式
	age   Column // 未导出字段忽略
}
```

## 创建表对象

```go
u := New[UserTable]("")   // 无别名
u2 := New[UserTable]("u2") // 别名为 u2
u3 := New[UserTable]("u3")
```

## 查询

* 查询单个表
	```go
	Select(u).From(u).String() // SELECT * FROM `user`
	// 以下示例均省略 .String()
	u.Select() // SELECT * FROM `user`
	Select(u.ID, u.Name).From(u) // SELECT `id`, `name` FROM `user`
	u.Select(u.ID, u.Name) // 同上
	```
* 限制返回数
	```go
	u.Select().Limit(10)            // SELECT * FROM `user` LIMIT 10
	u.Select().Limit(10).Offset(20) // SELECT * FROM `user` LIMIT 20, 10
	```
* 排序
	```go
	u.Select().OrderBy(u.ID.Desc()).OrderBy(u.Name.Asc()) // SELECT * FROM `user` ORDER BY `id` DESC, `name`
	u.Select().OrderBy(u.ID.Asc(), u.Name.Desc())         // SELECT * FROM `user` ORDER BY `id`, `name` DESC"
	```
* 分组
	```go
	u.Select().GroupBy(u.Name).GroupBy(u.ID) // SELECT * FROM `user` GROUP BY `name`, `id`
	u.Select().GroupBy(u.Name, u.ID)
	```
* 加锁
	```go
	u.Select().LockForShare()  // SELECT * FROM `user` FOR SHARE
	u.Select().LockForUpdate() // SELECT * FROM `user` FOR UPDATE
	```
* 查询条件
	```go
	u.Select().Where(u.ID.Eq(Expr("1"))) // SELECT * FROM `user` WHERE `id` = 1
	u.Select().Where(u.ID.Eq(PH))        // SELECT * FROM `user` WHERE `id` = ?
	u.Select().Where(u.ID.Eq(nil))       // SELECT * FROM `user` WHERE `id` IS NULL
	u.Select().Where(u.ID.In(PH))        // SELECT * FROM `user` WHERE `id` IN (?)
	u.Select().Where(And(u.ID.Eq(Expr("1")), u1.Name.Eq(PH), Not(Or(u.ID.Ne(Expr("2")), u.ID.Gt(Expr("3"))))))    // SELECT * FROM `user` WHERE `id` = 1 AND `name` = ? AND (NOT (`id` != 2 OR `id` > 3))
	u.Select().Where(u.ID.Eq(Expr("1")).And(u1.Name.Eq(PH)).And(u.ID.Ne(Expr("2")).Or(u.ID.Gt(Expr("3")))).Not()) // 同上
	```
	可用比较表达式有 `Eq`、`Ne`、`Gt`、`Ge`、`Lt`、`Le`、`In` 和 `NotIn`，逻辑表达式有 `And`、`Or` 和 `Not`。`PH` 是占位符的缩写。
* Join
	```go
	Select(u).FromJoin(u.InnerJoin(u2, u.ID.Eq(u2.ID))) // SELECT `user`.* FROM `user` JOIN `user` AS `u2` ON `u`.`id` = `u2`.`id`
	Select(u, u2.ID.As("other_id")).FromJoin(u.LeftJoin(u2, u.Name.Eq(u2.Name)).OuterJoin(u3, u2.ID.Eq(u3.ID))) // SELECT `u`.*, `u2`.`id` AS `other_id` FROM `user` LEFT JOIN `user` AS `u2` ON `u`.`name` = `u2`.`name` OUTER JON `user` AS `u3` ON `u2`.`id` = `u3`.`id`
	```
	当有 join 时，会自动使用别名，并引入表名。注意不能用 `From`，而要用 `FromJoin`。可用 join 方式有 `InnerJoin`、`LeftJoin`、`RightJoin` 和 `OuterJoin`。
* 函数
	```go
	Select(Func("SUM", u.ID).As("sum"), Func("COUNT", Expr("1"))).From(u)      // SELECT SUM(`id`) AS `sum`, COUNT(1) FROM `user`
	Select(Func("GROUP_CONCAT", Concat(u.Name, OrderBys{u.ID.Asc()}))).From(u) // SELECT GROUP_CONCAT(`name` ORDER BY `id`) FROM `user`
	```
	`Concat` 用于将多个表达式连接起来，因为 `GROUP_CONCAT` 这个 MySQL 函数比较特殊，它不用 ", " 来分隔表达式。`OrderBys` 这个结构可以输出 ` ORDER BY ...`。
* 子查询
	```go
	Select(Expr("*")).From(u).Where(u.ID.In(Select(Func("DISTINCT", u.ID)).From(u))) // SELECT * FROM `user` WHERE `id` IN (SELECT DISTINCT(`id`) FROM `user`)
	```

## 插入
* 插入默认值
	```go
	Insert(u)          // INSERT INTO `user` () VALUES ()
	u.Insert()         // 同上
	Insert(u).Ignore() // INSERT IGNORE INTO `user` () VALUES ()
	```
* 设置插入值
	```go
	Insert(u).Columns(u.ID, u.Name)                           // INSERT INTO `user` (`id`, `name`) VALUES (?, ?)
	u.Insert(u.ID, u.Name)                                    // 同上
	Insert(u).Columns(u.ID, u.Name).Values(nil, Expr(`"1"`))  // INSERT INTO `user` (`id`, `name`) VALUES (NULL, "1")
	Insert(u).Columns(u.ID, u.Name).NamedValues(u.ID, u.Name) // INSERT INTO `user` (`id`, `name`) VALUES (:id, :name)
	Insert(u).Columns(u.ID, u.Name).NamedValues()             // 同上，可自动使用 Columns 来作为 NamedValues
	```
	当不提供 `Values` 时，会自动根据 `Columns` 的数量生成相应数量的占位符。
* 冲突时更新
	```go
	Insert(u).Columns(u.ID, u.Name).OnDuplicateKeyUpdate(u.ID.Assign(u.ID.Plus(PH))) // INSERT INTO `user` (`id`, `name`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `id`=`id`+?
	```
	可用的操作有 `Plus`、`Minus`、`Multiply`、`Div` 和 `Mod`。
* INSERT INTO ... SELECT ...
	```go
	Insert(u).Columns(u.ID, u.Name).Select(u, Expr("1"), u.Name) // INSERT INTO `user` (`id`, `name`) SELECT 1, `name` FROM `user`
	```

## 更新
* 更新全表
	```go
	Update(u).Set(u.Name.Assign(Expr(`"1"`))) // UPDATE `user` SET `name`="1"
	u.Update(u.ID.Assign(Expr("1")))          // UPDATE `user` SET `id`=1
	```
* 条件更新
	```go
	Update(u).Set(u.ID.Assign(u.ID.Plus(Expr("1")))).Where(u.ID.Gt(PH)) // UPDATE `user` SET `name`=`id`+1 WHERE `id` > ?
	```
* 限制更新条数
	```go
	Update(u).Set(u.ID.Assign(Expr("1"))).OrderBy(u.Name.Asc(), u.ID.Desc()).Limit(10) // UPDATE `user` SET `id`=1 ORDER BY `name`, `id` DESC LIMIT 10
	```

## 删除
* 删除全表
	```go
	Delete(u)  // DELETE `user`
	u.Delete() // 同上
	```
* 条件删除全表
	```go
	Delete(u).Where(u.ID.Gt(PH)) // DELETE `user` WHERE `id` > ?
	```
* 限制删除条数
	```go
	Delete(u).OrderBy(u.Name.Asc(), u.ID.Desc()).Limit(10) // DELETE `user` ORDER BY `name`, `id` DESC LIMIT 10
	```

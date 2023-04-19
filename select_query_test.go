package sb

import "testing"

type UserTable struct {
	Table `db:"user"`
	ID    Column `json:"id" db:"id"`
	Name  Column `db:"name"`
}

type DeptTable struct {
	Table `db:"dept"`
	ID    Column `db:"id"`
	Name  Column `db:"name"`
}

type DeptUserTable struct {
	Table  `db:"dept_user"`
	DeptID Column
	UserID Column
}

func TestSelectQuery(t *testing.T) {
	u1 := New[UserTable]("u1")
	u2 := New[UserTable]("u2")
	d1 := New[DeptTable]("d1")
	// d2 := New[DeptTable]("d2")
	du := New[DeptUserTable]("du")

	tests := []struct {
		query    *SelectQuery
		expected string
	}{
		{
			query:    Select(u1).From(u1),
			expected: "SELECT * FROM `user`",
		},
		{
			query:    Select(u1.ID, u1.Name).From(u1),
			expected: "SELECT `id`, `name` FROM `user`",
		},
		{
			query:    Select(u1).From(u1).Limit(10),
			expected: "SELECT * FROM `user` LIMIT 10",
		},
		{
			query:    Select(u1).From(u1).Limit(10).Offset(20),
			expected: "SELECT * FROM `user` LIMIT 20, 10",
		},
		{
			query:    Select(u1).From(u1).Offset(20),
			expected: "SELECT * FROM `user` LIMIT 20, 0",
		},
		{
			query:    Select(u1).From(u1).OrderBy(u1.ID.Desc()).OrderBy(u1.Name.Asc()),
			expected: "SELECT * FROM `user` ORDER BY `id` DESC, `name`",
		},
		{
			query:    Select(u1).From(u1).OrderBy(u1.ID.Asc(), u1.Name.Desc()),
			expected: "SELECT * FROM `user` ORDER BY `id`, `name` DESC",
		},
		{
			query:    Select(u1).From(u1).GroupBy(u1.Name).GroupBy(u1.ID),
			expected: "SELECT * FROM `user` GROUP BY `name`, `id`",
		},
		{
			query:    Select(d1).From(d1).GroupBy(d1.Name, d1.ID),
			expected: "SELECT * FROM `dept` GROUP BY `name`, `id`",
		},
		{
			query:    Select(u1).From(u1).LockForShare(),
			expected: "SELECT * FROM `user` FOR SHARE",
		},
		{
			query:    Select(u1).From(u1).LockForUpdate(),
			expected: "SELECT * FROM `user` FOR UPDATE",
		},
		{
			query:    Select(u1).From(u1).Where(u1.ID.Eq(Expr("1"))),
			expected: "SELECT * FROM `user` WHERE `id` = 1",
		},
		{
			query:    Select(u1).From(u1).Where(And(u1.ID.Eq(Expr("1")), Or(u1.ID.Ne(Expr("2")), u1.ID.Gt(Expr("3"))))),
			expected: "SELECT * FROM `user` WHERE `id` = 1 AND (`id` != 2 OR `id` > 3)",
		},
		{
			query:    Select(u1).From(u1).Where(u1.ID.Eq(Placeholder)),
			expected: "SELECT * FROM `user` WHERE `id` = ?",
		},
		{
			query:    Select(u1).From(u1).Where(u1.ID.Eq(nil)),
			expected: "SELECT * FROM `user` WHERE `id` IS NULL",
		},
		{
			query:    Select(u1).From(u1).Where(u1.ID.In(Placeholder)),
			expected: "SELECT * FROM `user` WHERE `id` IN (?)",
		},
		{
			query:    Select(u1).FromJoin(u1.InnerJoin(du, u1.ID.Eq(du.UserID))),
			expected: "SELECT `u1`.* FROM `user` AS `u1` JOIN `dept_user` AS `du` ON `u1`.`id` = `du`.`userid`",
		},
		{
			query:    Select(u1, u2.ID.As("other_id")).FromJoin(u1.InnerJoin(u2, u1.Name.Eq(u2.Name))).Where(u2.ID.Gt(Expr("1"))),
			expected: "SELECT `u1`.*, `u2`.`id` AS `other_id` FROM `user` AS `u1` JOIN `user` AS `u2` ON `u1`.`name` = `u2`.`name` WHERE `other_id` > 1",
		},
		{
			query:    Select(u1).FromJoin(u1.InnerJoin(du, u1.ID.Eq(du.UserID)).LeftJoin(d1, d1.ID.Eq(du.DeptID))),
			expected: "SELECT `u1`.* FROM `user` AS `u1` JOIN `dept_user` AS `du` ON `u1`.`id` = `du`.`userid` LEFT JOIN `dept` AS `d1` ON `d1`.`id` = `du`.`deptid`",
		},
		{
			query:    Select(Func("SUM", u1.ID).As("sum"), Func("COUNT", Expr("1"))).From(u1),
			expected: "SELECT SUM(`id`) AS `sum`, COUNT(1) FROM `user`",
		},
		{
			query:    Select(Func("GROUP_CONCAT", Concat(u1.Name, OrderBys{u1.ID.Asc()}))).From(u1),
			expected: "SELECT GROUP_CONCAT(`name` ORDER BY `id`) FROM `user`",
		},
		{
			query:    Select(Expr("*")).From(u1).Where(u1.ID.In(Select(Func("DISTINCT", du.UserID)).From(du))),
			expected: "SELECT * FROM `user` WHERE `id` IN (SELECT DISTINCT(`userid`) FROM `dept_user`)",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.query.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}

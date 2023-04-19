package sb

import "testing"

func TestInsertQuery(t *testing.T) {
	u := New[UserTable]("u")

	tests := []struct {
		query    *InsertQuery
		expected string
	}{
		{
			query:    Insert(u),
			expected: "INSERT INTO `user` () VALUES ()",
		},
		{
			query:    Insert(u).Ignore(),
			expected: "INSERT IGNORE INTO `user` () VALUES ()",
		},
		{
			query:    Insert(u).Columns(u.ID, u.Name),
			expected: "INSERT INTO `user` (`id`, `name`) VALUES (?, ?)",
		},
		{
			query:    Insert(u).Columns(u.ID, u.Name).Values(),
			expected: "INSERT INTO `user` (`id`, `name`) VALUES (?, ?)",
		},
		{
			query:    Insert(u).Columns(u.ID, u.Name).Values(Placeholder),
			expected: "INSERT INTO `user` (`id`, `name`) VALUES (?)", // 这是一个错误的 SQL，只用于测试能否正确生成
		},
		{
			query:    Insert(u).Columns(u.ID, u.Name).Values(nil, Expr(`"1"`)),
			expected: "INSERT INTO `user` (`id`, `name`) VALUES (NULL, \"1\")",
		},
		{
			query:    Insert(u).Columns(u.ID, u.Name).NamedValues(u.ID, u.Name),
			expected: "INSERT INTO `user` (`id`, `name`) VALUES (:id, :name)",
		},
		{
			query:    Insert(u).Columns(u.ID, u.Name).NamedValues(),
			expected: "INSERT INTO `user` (`id`, `name`) VALUES (:id, :name)",
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

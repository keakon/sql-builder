package sb

import "testing"

func TestDeleteQuery(t *testing.T) {
	u := New[UserTable]("u")

	tests := []struct {
		query    *DeleteQuery
		expected string
	}{
		{
			query:    Delete(u),
			expected: "DELETE `user`",
		},
		{
			query:    Delete(u).Where(u.ID.Gt(Placeholder)),
			expected: "DELETE `user` WHERE `id` > ?",
		},
		{
			query:    Delete(u).OrderBy(u.Name.Asc(), u.ID.Desc()).Limit(10),
			expected: "DELETE `user` ORDER BY `name`, `id` DESC LIMIT 10",
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

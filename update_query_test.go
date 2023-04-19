package sb

import "testing"

func TestUpdateQuery(t *testing.T) {
	u := New[UserTable]("u")

	tests := []struct {
		query    *UpdateQuery
		expected string
	}{
		{
			query:    Update(u).Set(u.Name.Assign(Expr(`"1"`))),
			expected: "UPDATE `user` SET `name`=\"1\"",
		},
		{
			query:    Update(u).Set(u.Name.Assign(u.ID.Plus(Expr("1")))).Where(u.ID.Gt(Placeholder)),
			expected: "UPDATE `user` SET `name`=`id`+1 WHERE `id` > ?",
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
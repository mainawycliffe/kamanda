package auth

import "testing"

func TestFindUserCriteria_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    FindUserCriteria
		want bool
	}{
		{"Test if Search by Email Address", ByUserEmailCriteria, true},
		{"Test if Search By Phone", ByUserPhoneCriteria, true},
		{"Test if Search By UID", ByUserUIDCriteria, true},
		{"Test for Invalid Search Criteria", 9, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("FindUserCriteria.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

package utils

import "testing"

func TestDocsFrontMatter(t *testing.T) {
	type args struct {
		filename string
		summary  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Auth Command Test",
			args: args{
				filename: "docs/content/commands/auth/kamanda_auth_addUser.md",
				summary:  "Add a new Firebase Email/Password user (Accepts Custom Claims)",
			},
			want: `---
title: "kamanda auth addUser"
slug: kamanda_auth_addUser
url: /commands/kamanda_auth_adduser/
summary: "Add a new Firebase Email/Password user (Accepts Custom Claims)"
---
`,
		},
		{
			name: "Logout Command Test",
			args: args{
				filename: "docs/content/commands/kamanda_logout.md",
				summary:  "Logout kamanda from Firebase",
			},
			want: `---
title: "kamanda logout"
slug: kamanda_logout
url: /commands/kamanda_logout/
summary: "Logout kamanda from Firebase"
---
`,
		},
		{
			name: "Version Command Test",
			args: args{
				filename: "docs/content/commands/kamanda_version.md",
				summary:  "Version will output the current build information",
			},
			want: `---
title: "kamanda version"
slug: kamanda_version
url: /commands/kamanda_version/
summary: "Version will output the current build information"
---
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DocsFrontMatter(tt.args.filename, tt.args.summary); got != tt.want {
				t.Errorf("DocsFrontMatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

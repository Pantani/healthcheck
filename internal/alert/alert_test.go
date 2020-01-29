package alert

import "testing"

func Test_getDescription(t *testing.T) {
	type args struct {
		namespace string
		name      string
		path      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test full values", args{
			namespace: "ethereum",
			name:      "block_height",
			path:      "api",
		}, "ethereum test block_height api"},
		{"test empty namespace", args{
			namespace: "",
			name:      "block_height",
			path:      "api",
		}, "test block_height api"},
		{"test empty name", args{
			namespace: "ethereum",
			name:      "",
			path:      "api",
		}, "ethereum test api"},
		{"test empty path", args{
			namespace: "ethereum",
			name:      "block_height",
			path:      "",
		}, "ethereum test block_height"},
		{"test empty", args{
			namespace: "",
			name:      "",
			path:      "",
		}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDescription(tt.args.namespace, tt.args.name, tt.args.path); got != tt.want {
				t.Errorf("getDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

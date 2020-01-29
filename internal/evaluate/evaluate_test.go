package evaluate

import "testing"

func TestEvaluate(t *testing.T) {
	type args struct {
		exp       string
		lastValue interface{}
		newValue  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"success true case",
			args{
				exp:       "lastValue > newValue",
				lastValue: 10,
				newValue:  3,
			},
			true, false,
		},
		{
			"success false case",
			args{
				exp:       "lastValue > newValue",
				lastValue: 3,
				newValue:  33,
			},
			false, false,
		},
		{
			"string success true case",
			args{
				exp:      "len(newValue) > 0",
				newValue: "test",
			},
			true, false,
		},
		{
			"string success false case",
			args{
				exp:      "len(newValue) == 0",
				newValue: "test",
			},
			false, false,
		},
		{"failed case 1",
			args{
				exp:       "lastValue > newValue",
				lastValue: nil,
				newValue:  33,
			},
			false, true,
		},
		{
			"failed case 2",
			args{
				exp:       "lastValue > newValue",
				lastValue: 33,
				newValue:  nil,
			},
			false, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.args.exp, tt.args.lastValue, tt.args.newValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

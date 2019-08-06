package parser

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *Command
		wantErr bool
	}{
		{
			name:    "Empty line",
			args:    args{"\n"},
			wantErr: false,
			want:    nil,
		},
		{
			name:    "Wrong Prefix",
			args:    args{"This is not an mqtt message"},
			wantErr: false,
			want:    nil,
		},
		{
			name:    "No Command",
			args:    args{"mqtt"},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "Unknown Command",
			args:    args{"mqtt bananarama"},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "Subscribe",
			args:    args{"mqtt subscribe topic"},
			wantErr: false,
			want: &Command{
				Command: "subscribe",
				Topic:   "topic",
			},
		},
		{
			name:    "Subscribe not enough parameter",
			args:    args{"mqtt subscribe"},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "Publish",
			args:    args{"mqtt publish topic 1.231"},
			wantErr: false,
			want: &Command{
				Command: "publish",
				Topic:   "topic",
				Value:   "1.231",
			},
		},
		{
			name:    "Publish not enough parameter",
			args:    args{"mqtt publish"},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

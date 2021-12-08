package talktime

import (
	"reflect"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Transcript
		wantErr bool
	}{
		{
			name:    "not found",
			args:    args{filename: "notfound.txt"},
			wantErr: true,
			want:    nil,
		},
		{
			name:    "found",
			args:    args{filename: "sample.vtt"},
			wantErr: false,
			want: &Transcript{
				Summaries: map[string]Summary{
					"Johnny": Summary{Duration: time.Millisecond * 20850},
					"Bob":    Summary{Duration: time.Millisecond * 3575},
				},
				Duration: time.Millisecond * 28320,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("run %s, Open() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got == nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() = %+v, want %+v", got, tt.want)
			}

		})
	}
}

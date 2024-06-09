package botcode

import (
	"reflect"
	"testing"
	"time"
)

func Test_monthChange(t *testing.T) {
	type args struct {
		month  time.Month
		change int
	}
	tests := []struct {
		name string
		args args
		want time.Month
	}{
		// TODO: Add test cases.
		{"December -12", args{time.December, -12}, time.December},
		{"January -1", args{time.January, -1}, time.December},
		{"February -1", args{time.February, -1}, time.January},
		{"December +13", args{time.December, 13}, time.January},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := monthChange(tt.args.month, tt.args.change); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("monthChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

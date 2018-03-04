package parser

import (
	"reflect"
	"testing"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
)

func Test_acestreamParser_Parse(t *testing.T) {
	type args struct {
		comment string
	}
	tests := []struct {
		name string
		a    *acestreamParser
		args args
		want []*soccerstreams.Stream
	}{
		{
			args: args{
				comment: "acestream://cbb102edb320d1cd826af5f9342b10d66efe6fd6  [1920x1080]  [BT Sport]",
			},
			want: []*soccerstreams.Stream{
				{
					Link:    "acestream://cbb102edb320d1cd826af5f9342b10d66efe6fd6",
					Quality: "HD",
					Channel: "BT Sport",
				},
			},
		},
		{
			args: args{
				comment: "acestream://cbb102edb320d1cd826af5f9342b10d66efe6fd6  [520p]  [BT Sport 2]  [English]",
			},
			want: []*soccerstreams.Stream{
				{
					Link:    "acestream://cbb102edb320d1cd826af5f9342b10d66efe6fd6",
					Quality: "520p",
					Channel: "BT Sport 2",
				},
			},
		},
		{
			args: args{
				comment: "acestream://cbb102edb320d1cd826af5f9342b10d66efe6fd6  [520P]  [BT Sport 2]  [Unknown Language]",
			},
			want: []*soccerstreams.Stream{
				{
					Link:    "acestream://cbb102edb320d1cd826af5f9342b10d66efe6fd6",
					Quality: "520p",
					Channel: "Acestream",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &acestreamParser{}
			if got := a.Parse(tt.args.comment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("acestreamParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

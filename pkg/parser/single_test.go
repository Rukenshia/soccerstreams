package parser

import (
	"reflect"
	"testing"

	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
)

func Test_singleStreamParser_parseQuality(t *testing.T) {
	type args struct {
		fragment string
	}
	tests := []struct {
		name  string
		w     *singleStreamParser
		args  args
		want  string
		want1 bool
	}{
		{
			args: args{
				fragment: "",
			},
			want:  "",
			want1: false,
		},
		{
			args: args{
				fragment: "SD",
			},
			want:  "SD",
			want1: true,
		},
		{
			args: args{
				fragment: "HD",
			},
			want:  "HD",
			want1: true,
		},
		{
			args: args{
				fragment: "**HD**",
			},
			want:  "HD",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &singleStreamParser{}
			got, got1 := w.parseQuality(tt.args.fragment)
			if got != tt.want {
				t.Errorf("singleStreamParser.parseQuality() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("singleStreamParser.parseQuality() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_singleStreamParser_parseChannel(t *testing.T) {
	type args struct {
		fragment string
	}
	tests := []struct {
		name  string
		w     *singleStreamParser
		args  args
		want  string
		want1 string
		want2 bool
		want3 bool
	}{
		{
			args: args{
				fragment: "",
			},
			want3: false,
		},
		{
			args: args{
				fragment: "[X](Y)",
			},
			want:  "X",
			want1: "Y",
			want3: true,
		},
		{
			args: args{
				fragment: "[ X ](Y)",
			},
			want:  "X",
			want1: "Y",
			want3: true,
		},
		{
			args: args{
				fragment: "[ X ] (Y)",
			},
			want:  "X",
			want1: "Y",
			want3: true,
		},
		{
			args: args{
				fragment: "[ X ] ( Y)",
			},
			want:  "X",
			want1: "Y",
			want3: true,
		},
		{
			args: args{
				fragment: "[ X ](Y )",
			},
			want:  "X",
			want1: "Y",
			want3: true,
		},
		{
			args: args{
				fragment: "[ X ]( Y )",
			},
			want:  "X",
			want1: "Y",
			want3: true,
		},
		{
			args: args{
				fragment: "[ X | **HD** ]( Y )",
			},
			want:  "X | **HD**",
			want1: "Y",
			want2: true,
			want3: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &singleStreamParser{}
			got, got1, got2, got3 := w.parseChannel(tt.args.fragment)
			if got != tt.want {
				t.Errorf("singleStreamParser.parseChannel() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("singleStreamParser.parseChannel() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("singleStreamParser.parseChannel() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("singleStreamParser.parseChannel() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func Test_singleStreamParser_Parse(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		w    *singleStreamParser
		args args
		want []*soccerstreams.Stream
	}{
		{
			args: args{
				message: `**HD** | [ Everton vs Brighton | BEIN 1 | ARAB | Clicks :1 | MISR: 1 Mbps | Mobile No | Chromecast: Yes ](http://arsenewenger.cf/ch24.html)`,
			},
			want: []*soccerstreams.Stream{
				{
					Channel:        "Everton vs Brighton",
					Link:           "http://arsenewenger.cf/ch24.html",
					IsNSFW:         false,
					MISR:           "1 Mbps",
					Quality:        "HD",
					MobileFriendly: false,
					Clicks:         1,
				},
			},
		}, {
			args: args{
				message: `**HD** | [ Newcastle United vs  Southampton | English | ads 2 | Mobile: Yes ] (http://welovesports.xyz/newcastle-vs-southampton)`,
			},
			want: []*soccerstreams.Stream{
				{
					Channel:        "Newcastle United vs  Southampton",
					Link:           "http://welovesports.xyz/newcastle-vs-southampton",
					IsNSFW:         false,
					MISR:           "",
					Quality:        "HD",
					MobileFriendly: true,
					Clicks:         0,
				},
			},
		},
		{
			args: args{
				message: `520p Stream | EN | [Basel vs Manchester City](http://buffstreamz.com/watch/soccer-2.php) | MISR : 1mbps | Mobile : Yes | Clicks : 3`,
			},
			want: []*soccerstreams.Stream{
				{
					Channel:        "Basel vs Manchester City",
					Link:           "http://buffstreamz.com/watch/soccer-2.php",
					IsNSFW:         false,
					MISR:           "1mbps",
					MobileFriendly: true,
					Clicks:         3,
				},
			},
		},
		{
			args: args{
				message: `520p Stream | EN | [Juventus vs Tottenham Hotspur](http://buffstreamz.com/watch/soccer.php) | MISR : 1mbps | Mobile : Yes | Clicks : 3`,
			},
			want: []*soccerstreams.Stream{
				{
					Channel:        "Juventus vs Tottenham Hotspur",
					Link:           "http://buffstreamz.com/watch/soccer.php",
					IsNSFW:         false,
					MISR:           "1mbps",
					MobileFriendly: true,
					Clicks:         3,
				},
			},
		},
		{
			args: args{
				message: `**HD** | [ENGLISH TSN4 1080p] (http://foundationsports.com/crvcsk/) | MISR : 3mbps | Ad Overlay : 1 | Clicks : 2  | Mobile : Yes.`,
			},
			want: []*soccerstreams.Stream{
				{
					Quality:        "HD",
					Channel:        "ENGLISH TSN4 1080p",
					Link:           "http://foundationsports.com/crvcsk/",
					IsNSFW:         false,
					MISR:           "3mbps",
					MobileFriendly: true,
					Clicks:         2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &singleStreamParser{}
			if got := w.Parse(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("singleStreamParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

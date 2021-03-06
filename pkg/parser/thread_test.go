package parser

import (
	"reflect"
	"testing"
	"time"

	"github.com/Rukenshia/graw/reddit"
	"github.com/Rukenshia/soccerstreams/pkg/soccerstreams"
)

func Test_threadParser_Parse(t *testing.T) {
	now := time.Now()
	gmt, _ := time.LoadLocation("GMT")
	gmt805pm := time.Date(now.Year(), now.Month(), now.Day(), 20, 5, 0, 0, gmt)
	gmt1am := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, gmt)

	type args struct {
		p *reddit.Post
	}
	tests := []struct {
		name string
		t    *threadParser
		args args
		want *soccerstreams.Matchthread
	}{
		{
			args: args{
				p: &reddit.Post{
					Title:         "[1:00 GMT] København vs Atlético Madrid",
					LinkFlairText: "UEFA Europa League",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "København",
				Team2:   "Atlético Madrid",
				Kickoff: &gmt1am,
			},
		},
		{
			args: args{
				p: &reddit.Post{
					Title:         "[20:05 GMT] København vs Atlético Madrid",
					LinkFlairText: "UEFA Europa League",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "København",
				Team2:   "Atlético Madrid",
				Kickoff: &gmt805pm,
			},
		},
		{
			args: args{
				p: &reddit.Post{
					Title:         "[20:05GMT] København vs Atlético Madrid",
					LinkFlairText: "UEFA Europa League",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "København",
				Team2:   "Atlético Madrid",
				Kickoff: &gmt805pm,
			},
		},
		{
			args: args{
				p: &reddit.Post{
					Title:         "[20:05 GMT] Celtic vs. Zenit",
					LinkFlairText: "UEFA Europa League",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "Celtic",
				Team2:   "Zenit",
				Kickoff: &gmt805pm,
			},
		},
		{
			args: args{
				p: &reddit.Post{
					Title:         "[20:05 GMT] Celtic x Zenit",
					LinkFlairText: "UEFA Europa League",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "Celtic",
				Team2:   "Zenit",
				Kickoff: &gmt805pm,
			},
		},
		{
			args: args{
				p: &reddit.Post{
					Title: "[20.05 GMT] Celtic v Zenit",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "Celtic",
				Team2:   "Zenit",
				Kickoff: &gmt805pm,
			},
		},
		{
			args: args{
				p: &reddit.Post{
					Title: "[20:05 GMT] Celtic v Zenit",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "Celtic",
				Team2:   "Zenit",
				Kickoff: &gmt805pm,
			},
		},
		{
			args: args{
				p: &reddit.Post{
					Title: "[20:05 GMT] Celtic - Zenit",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "Celtic",
				Team2:   "Zenit",
				Kickoff: &gmt805pm,
			},
		},
		{
			args: args{
				p: &reddit.Post{
					Title: "[20:05 GMT] Celtic VS Zenit",
				},
			},
			want: &soccerstreams.Matchthread{
				Team1:   "Celtic",
				Team2:   "Zenit",
				Kickoff: &gmt805pm,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &threadParser{}
			if got := tp.Parse(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("threadParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

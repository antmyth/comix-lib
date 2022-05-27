package comicvine

import "testing"

func TestExtractIdFromSiteUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "blackammer reborn",
			args: args{
				url: "https://comicvine.gamespot.com/black-hammer-reborn-part-1-1-volume-5/4000-908362/",
			},
			want: "4000-908362",
		},
		{
			name: "2",
			args: args{
				url: "https://comicvine.gamespot.com/camelot-3000-1-the-past-and-future-king/4000-22636/",
			},
			want: "4000-22636",
		},
		{
			name: "3",
			args: args{
				url: "https://comicvine.gamespot.com/camelot-3000-2-many-are-called/4000-22757/",
			},
			want: "4000-22757",
		},
		{
			name: "4",
			args: args{
				url: "https://comicvine.gamespot.com/camelot-3000-3-knight-quest/4000-22830/",
			},
			want: "4000-22830",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractIdFromSiteUrl(tt.args.url); got != tt.want {
				t.Errorf("ExtractIdFromSiteUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

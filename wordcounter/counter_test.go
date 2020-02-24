package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestCounter(t *testing.T) {
	tests := []struct {
		name string
		args []Comment
		want WordCounts
	}{
		{"case 01", []Comment{Comment{Body: "Working with God counts are things all"},
			Comment{Body: "with God all things are possible word counts"},
			Comment{Body: "with God all things are possible God Testing counts with"},
			Comment{Body: "with God all things are possible Money"},
			Comment{Body: "Testing word counts are God things all with"},
			Comment{Body: "Testing word counts with things possible are"},
			Comment{Body: "God is  all things with Testing word counts are God things all with God"},
			Comment{Body: "Testing is  with God all word"},
			Comment{Body: "Money is all God counts word are things God all with"},
		},
			WordCounts{WordCount{"working", 1},
				WordCount{"money", 2},
				WordCount{"is", 3},
				WordCount{"possible", 4},
				WordCount{"testing", 5},
				WordCount{"word", 6},
				WordCount{"counts", 7},
				WordCount{"are", 8},
				WordCount{"things", 9},
				WordCount{"all", 10},
				WordCount{"with", 11},
				WordCount{"god", 12},
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Counter(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Counter() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestProcessComment(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		limit int
		want WordCounts
		wantErr bool
	}{
		{ "case 01", []byte(`[
	{ "body": "all starting with golang all the way gun"},
	{ "body": "all starting with the golang gun"},
	{ "body": "with golang all the way the gun"},
	{ "body": "with golang all the way with gun"},
	{ "body": "the man with all gun"}]`), 5,
	WordCounts{WordCount{"man", 1},
			WordCount{"starting", 2},
			WordCount{"way", 3},
			WordCount{"golang", 4},
			WordCount{"gun", 5},
		}, false},
		{ "case 02", []byte(`[
	{ "body": "all starting with golang all the way gun]`), 5, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := ioutil.NopCloser(bytes.NewBuffer(tt.args))
			got, err  := ProcessComment(body, tt.limit)
			if err != nil && !tt.wantErr {
				t.Errorf("unable to process comment :%v", err)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessComment() = %s, want %s", got, tt.want)
			}
		})
	}



}
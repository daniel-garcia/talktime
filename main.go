package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/daniel-garcia/talktime/pkg/talktime"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("talktime only accepts a single zoom text transcript as an argument")
	}

	ts, err := talktime.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not process transcript: %s", err)
	}

	fmt.Printf("%-30s: %s\n", "Meeting Duration", ts.Duration)

	var summaries []talktime.Summary

	var talkTime time.Duration
	for _, summary := range ts.Summaries {
		summaries = append(summaries, summary)
		talkTime = talkTime + summary.Duration
	}

	silence := ts.Duration - talkTime
	summaries = append(summaries, talktime.Summary{Name: "!silence!", Duration: silence})

	sort.Sort(talktime.ByDuration(summaries))
	for _, summary := range summaries {
		fmt.Printf("%-30s: %-20s\t %03.1f %%\n", summary.Name, summary.Duration, float64(summary.Duration)/float64(ts.Duration)*100)
	}

}

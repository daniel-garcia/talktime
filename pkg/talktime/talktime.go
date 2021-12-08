package talktime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Transcript struct {
	Duration  time.Duration
	Summaries map[string]Summary
}

type Summary struct {
	Duration time.Duration
	Name     string
}

type ByDuration []Summary

func (s ByDuration) Len() int           { return len(s) }
func (s ByDuration) Less(i, j int) bool { return s[i].Duration < s[j].Duration }
func (s ByDuration) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type ByName []Summary

func (s ByName) Len() int           { return len(s) }
func (s ByName) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s ByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func Open(filename string) (*Transcript, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	header := scanner.Text()
	if header != ".WEBVTT" {
		return nil, fmt.Errorf("bad header, got %s", header)
	}

	people := make(map[string]Summary)

	var totalDuration time.Duration
	for {

		if !scanner.Scan() {
			break
		}
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		if !scanner.Scan() {
			break
		}
		timeline := scanner.Text()

		if !scanner.Scan() {
			break
		}
		wordline := scanner.Text()

		duration, dtime, err := parseTimeline(timeline)
		if err != nil {
			return nil, fmt.Errorf("speaking line %s, bad timeline: %v => %s", line, timeline, err)
		}
		totalDuration = dtime
		person, err := parseWordline(wordline)
		if err != nil {
			return nil, fmt.Errorf("speaking line %s, bad wordline: %v => %s", line, wordline, err)
		}

		summary, found := people[person]
		if !found {
			summary = Summary{Name: person}
			people[person] = summary
		}
		summary.Duration = summary.Duration + duration
		people[person] = summary
	}

	return &Transcript{
		Duration:  totalDuration,
		Summaries: people,
	}, nil
}

func parseWordline(s string) (string, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "unk", nil
	}
	return parts[0], nil
}

const (
	timeLayout = "2006-01-02 15:04:05.000"
)

func parseTimeline(s string) (time.Duration, time.Duration, error) {

	parts := strings.SplitN(s, " ", 3)
	if len(parts) != 3 {
		return time.Second * 0, time.Second * 0, fmt.Errorf("malformed timeline")
	}
	if parts[1] != "-->" {
		return time.Second * 0, time.Second * 0, fmt.Errorf("expected --> in middle of timeline, got %s", parts[1])
	}

	fromStr := "2006-01-02 " + parts[0]
	from, err := time.Parse(timeLayout, fromStr)
	if err != nil {
		return time.Second * 0, time.Second * 0, fmt.Errorf("could not parse from time: %s", err)
	}
	toStr := "2006-01-02 " + parts[2]
	to, err := time.Parse(timeLayout, toStr)
	if err != nil {
		return time.Second * 0, time.Second * 0, fmt.Errorf("could not parse to time: %s", err)
	}

	startTime, _ := time.Parse(timeLayout, "2006-01-02 00:00:00.000")
	return to.Sub(from), to.Sub(startTime), nil
}

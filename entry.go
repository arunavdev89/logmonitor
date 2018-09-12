package main

import (
	"fmt"
	"github.com/satyrius/gonx"
	"strings"
)

const (
	callerIp  = "caller_ip"
	callerId  = "caller_id"
	timestamp = "timestamp"
	method    = "method"
	website   = "website"
	status    = "status"
	size      = "size"
	protocol  = "protocol"
)

var (
	dimensionNames = []string{callerId, website, method, status}
)

//Entry represents a parsed log line entry
type Entry struct {
	CallerIp  string
	CallerId  string
	Timestamp string
	Method    string
	Website   string
	Status    string
}

//represents w3c-formatted HTTP access log
type LogFormat struct {
	Format string
}

func NewW3CLogFormat() *LogFormat {
	return &LogFormat{
		Format: fmt.Sprintf("$%s - $%s [$%s] \"$%s $%s $%s\" $%s $%s",
			callerIp,
			callerId,
			timestamp,
			method,
			website,
			protocol,
			status,
			size),
	}
}

func NewEntry(line string) (entry *Entry, err error) {
	format := NewW3CLogFormat()
	logparser := gonx.NewParser(format.Format)
	fields, err := logparser.ParseString(line)
	entry = &Entry{}
	if err != nil {
		return
	}
	entryCallerIp, _ := fields.Field(callerIp)
	entryCallerId, _ := fields.Field(callerId)
	entryTimestamp, _ := fields.Field(timestamp)
	entryMethod, _ := fields.Field(method)
	entryWebsite, _ := fields.Field(website)
	//take upto 3rd / characters
	entryWebsite = splitUptoNthSlash(entryWebsite, 3)

	entryStatus, _ := fields.Field(status)
	entry = &Entry{
		CallerIp:  entryCallerIp,
		CallerId:  entryCallerId,
		Timestamp: entryTimestamp,
		Method:    entryMethod,
		Website:   entryWebsite,
		Status:    entryStatus,
	}
	return
}

func (e *Entry) Dimension() *MetricDimension {
	return &MetricDimension{
		dimensionNames:  dimensionNames,
		dimensionValues: []string{e.CallerId, e.Website, e.Method, e.Status},
	}
}

func splitUptoNthSlash(s string, n int) string {
	split := strings.SplitN(s, "/", n+2)
	return strings.Join(split[:n+1], "/")
}

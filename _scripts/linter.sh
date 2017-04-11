#!/bin/sh

param=$1

retool do gometalinter --config=.gometalinter.json ./... \
	--linter="aligncheck:aligncheck ${param} {path}:^(?:[^:]+: )?(?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.+)$" \
	--linter="deadcode:deadcode ${param} {path}:^deadcode: (?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.*)$" \
	--linter="dupl:dupl ${param} -plumbing -threshold {duplthreshold} {path}/*.go:^(?P<path>.*?\.go):(?P<line>\d+)-\d+:\s*(?P<message>.*)$" \
	--linter="errcheck:errcheck ${param} -abspath {path}:PATH:LINE:COL:MESSAGE" \
	--linter="gas:gas ${param} -fmt=csv {path}/*.go:^(?P<path>.*?\.go),(?P<line>\d+),(?P<message>[^,]+,[^,]+,[^,]+)" \
	--linter="goconst:goconst ${param} -min-occurrences {min_occurrences} -min-length {min_const_length} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="gocyclo:gocyclo ${param} -over {mincyclo} {path}:^(?P<cyclo>\d+)\s+\S+\s(?P<function>\S+)\s+(?P<path>.*?\.go):(?P<line>\d+):(\d+)$" \
	--linter="gofmt:gofmt ${param} -l -s {path}/*.go:^(?P<path>.*?\.go)$" \
	--linter="goimports:goimports ${param} -l {path}/*.go:^(?P<path>.*?\.go)$" \
	--linter="golint:golint ${param} -min_confidence {min_confidence} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="gosimple:gosimple ${param} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="gotype:gotype ${param} -e {tests=-a} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="ineffassign:ineffassign ${param} -n {path}:PATH:LINE:COL:MESSAGE" \
	--linter="interfacer:interfacer ${param} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="lll:lll ${param} -g -l {maxlinelength} {path}/*.go:PATH:LINE:MESSAGE" \
	--linter="misspell:misspell ${param} -j 1 {path}/*.go:PATH:LINE:COL:MESSAGE" \
	--linter="safesql:safesql ${param} {path}:^- (?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+)$" \
	--linter="staticcheck:staticcheck ${param} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="structcheck:structcheck ${param} {tests=-t} {path}:^(?:[^:]+: )?(?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.+)$" \
	--linter="test:go ${param} test {path}:^--- FAIL: .*$\s+(?P<path>.*?\.go):(?P<line>\d+): (?P<message>.*)$" \
	--linter="testify:go ${param} test {path}:Location:\s+(?P<path>.*?\.go):(?P<line>\d+)$\s+Error:\s+(?P<message>[^\n]+)" \
	--linter="unconvert:unconvert ${param} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="unparam:unparam ${param} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="unused:unused ${param} {path}:PATH:LINE:COL:MESSAGE" \
	--linter="varcheck:varcheck ${param} {path}:^(?:[^:]+: )?(?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.*)$" \
	--linter="vet:go ${param} tool vet {path}/*.go:^(?:vet:.*?\.go:\s+(?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.*))|(?:(?P<path>.*?\.go):(?P<line>\d+):\s*(?P<message>.*))$" \
	--linter="vetshadow:go ${param} tool vet --shadow {path}/*.go:^(?:vet:.*?\.go:\s+(?P<path>.*?\.go):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.*))|(?:(?P<path>.*?\.go):(?P<line>\d+):\s*(?P<message>.*))$"

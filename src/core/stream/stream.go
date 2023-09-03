package stream

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type JsonStream struct {
	r io.ReadSeeker
	s *bufio.Scanner
}

func NewJsonStream(r io.ReadSeekCloser) *JsonStream {
	return &JsonStream{r: r, s: bufio.NewScanner(r)}
}

func (j *JsonStream) Offset(offset int64) error {
	_, err := j.r.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	j.s = bufio.NewScanner(j.r)
	return nil
}

func (j *JsonStream) ScanNext(rst interface{}, cond [][]string) (int64, interface{}, error) {
	for j.s.Scan() {
		text := j.s.Text()
		if j.match(text, cond) {
			offset, err := j.current()
			if err != nil {
				return 0, nil, err
			}
			rstObj, err := j.decode(text, rst)
			if err != nil {
				return offset, nil, err
			}
			return offset, rstObj, nil
		}
	}
	return 0, nil, io.EOF
}

func (j *JsonStream) current() (int64, error) {
	return j.r.Seek(0, io.SeekCurrent)
}

func (j *JsonStream) decode(text string, rst interface{}) (interface{}, error) {
	if err := json.Unmarshal([]byte(text), &rst); err != nil {
		return nil, err
	}
	return rst, nil
}

func (j *JsonStream) match(target string, match [][]string) bool {
	for _, m := range match {
		if j.containsAll(target, m) {
			return true
		}
	}
	return false
}

func (j *JsonStream) containsAll(target string, match []string) bool {
	for _, v := range match {
		if strings.Index(target, v) == -1 {
			return false
		}
	}
	return true
}

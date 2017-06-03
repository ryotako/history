package history

import (
	"bytes"
	"sort"
	"strings"
	tt "text/template"
	"time"

	ltsv "github.com/Songmu/go-ltsv"
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
)

type Record struct {
	ID      uint32
	Date    time.Time
	Command string
	Dir     string
	Branch  string
	Status  int
}

type Records []Record

func NewRecord() *Record {
	return &Record{
		ID:   uuid.New().ID(),
		Date: time.Now(),
	}
}

func (r *Record) SetCommand(arg string) { r.Command = arg }
func (r *Record) SetDir(arg string)     { r.Dir = arg }
func (r *Record) SetBranch(arg string)  { r.Branch = arg }
func (r *Record) SetStatus(arg int)     { r.Status = arg }

func (r *Record) Render(visible []string) (line string) {
	var tmpl *tt.Template
	format := "{{.ID}}"
	for _, v := range visible {
		format += "\t" + v
	}
	t, err := tt.New("format").Parse(format)
	if err != nil {
		return
	}
	tmpl = t
	if tmpl != nil {
		var b bytes.Buffer
		err := tmpl.Execute(&b, map[string]interface{}{
			"ID":      r.ID, // Required when parsing
			"Date":    r.Date.Format("2006-01-02"),
			"Time":    humanize.Time(r.Date),
			"Command": r.Command,
			"Dir":     r.Dir,
			"Branch":  r.Branch,
			"Status":  r.Status,
		})
		if err != nil {
			return
		}
		line = b.String()
	}
	return
}

func (r *Record) Unmarshal(line string) Record {
	ltsv.Unmarshal([]byte(line), r)
	return *r
}

func (r *Record) Marshal() ([]byte, error) {
	b, err := ltsv.Marshal(r)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func (r *Records) Filter(fn func(Record) bool) *Records {
	records := make(Records, 0)
	for _, record := range *r {
		if fn(record) {
			records = append(records, record)
		}
	}
	return &records
}

func (r *Records) Unique() {
	rs := make(Records, 0)
	encountered := map[string]bool{}
	for _, record := range *r {
		if !encountered[record.Command] {
			encountered[record.Command] = true
			rs = append(rs, record)
		}
	}
	*r = rs
}

func (r *Records) Reverse() {
	var rs Records
	for i := len(*r) - 1; i >= 0; i-- {
		rs = append(rs, (*r)[i])
	}
	*r = rs
}

func (r *Records) Grep(words []string) {
	for _, word := range words {
		*r = *r.Filter(func(r Record) bool {
			return strings.HasPrefix(r.Command, word)
		})
	}
}

func (r Records) Len() int           { return len(r) }
func (r Records) Less(i, j int) bool { return r[i].Date.Before(r[j].Date) }
func (r Records) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (r *Records) Sort() {
	sort.Sort(*r)
}

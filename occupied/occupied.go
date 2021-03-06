package occupied

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"text/template"
	"time"
)

type Record struct {
	Occupied bool
	Date     time.Time
}

func init() {
	http.HandleFunc("/", latest_html)
	http.HandleFunc("/latest.json", latest_json)
	http.HandleFunc("/record/opened", opened)
	http.HandleFunc("/record/closed", closed)
}

func get_latest_record(r *http.Request) (rec Record, err error) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Record").Order("-Date").Limit(1)
	records := make([]Record, 0, 1)
	if _, err := q.GetAll(c, &records); err != nil {
        return Record{}, err
	}
    if len(records) == 0 {
        return Record{false, time.Now()}, nil
    }
    return records[0], nil
}

func latest_json(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/json")
    add_standard_headers(w)
    var rec Record
    var err error
	if rec, err = get_latest_record(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := latestJsonTemplate.Execute(w, rec); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func add_standard_headers(w http.ResponseWriter) {
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

func latest_html(w http.ResponseWriter, r *http.Request) {
    add_standard_headers(w)
    var rec Record
    var err error
	if rec, err = get_latest_record(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := latestHtmlTemplate.Execute(w, rec); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var latestJsonTemplate = template.Must(template.New("latest_json").Parse(latestJsonTemplateStr))
var latestHtmlTemplate = template.Must(template.New("latest_html").Parse(latestHtmlTemplateStr))

const latestJsonTemplateStr = `{"occupied": {{.Occupied}}}`
const latestHtmlTemplateStr = `
<html>
<head>
<meta http-equiv="refresh" content="5">
<title>{{if .Occupied}}Occupied{{else}}Available{{end}}</title>
<link rel="icon" 
      type="image/png" 
      href="/static/img/{{if .Occupied}}df-poop-occupied{{else}}df-poop-vacant{{end}}.png">
</head><body>
<img src="/static/img/{{if .Occupied}}occupied{{else}}vacant{{end}}.jpg">
</body>
`

func opened(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	rec := Record{
		Occupied: false,
		Date:     time.Now(),
	}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Record", nil), &rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func closed(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	rec := Record{
		Occupied: true,
		Date:     time.Now(),
	}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Record", nil), &rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Command example runs a sample webserver that uses go-i18n/v2/i18n.
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var page = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<body>

<h1>{{.Title}}</h1>

{{range .Paragraphs}}<p>{{.}}</p>{{end}}

</body>
</html>
`))

var (
	translator *i18n.Bundle
)

func init() {
	translator = &i18n.Bundle{DefaultLanguage: language.English}
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(translator, lang, accept)

	name := r.FormValue("name")
	if name == "" {
		name = "Bob"
	}

	unreadEmailCount, _ := strconv.ParseInt(r.FormValue("unreadEmailCount"), 10, 64)

	helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HelloPerson",
			Other: "Hello {{.Name}}",
		},
		TemplateData: map[string]string{
			"Name": name,
		},
	})

	myUnreadEmails := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "MyUnreadEmails",
			Description: "The number of unread emails I have",
			One:         "I have {{.PluralCount}} unread email.",
			Other:       "I have {{.PluralCount}} unread emails.",
		},
		PluralCount: unreadEmailCount,
	})

	personUnreadEmails := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "PersonUnreadEmails",
			Description: "The number of unread emails a person has",
			One:         "{{.Name}} has {{.UnreadEmailCount}} unread email.",
			Other:       "{{.Name}} has {{.UnreadEmailCount}} unread emails.",
		},
		PluralCount: unreadEmailCount,
		TemplateData: map[string]interface{}{
			"Name":             name,
			"UnreadEmailCount": unreadEmailCount,
		},
	})

	err := page.Execute(w, map[string]interface{}{
		"Title": helloPerson,
		"Paragraphs": []string{
			myUnreadEmails,
			personUnreadEmails,
		},
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	translator.RegisterUnmarshalFunc("json", json.Unmarshal)
	translator.MustLoadMessageFile("active.es.json")

	http.HandleFunc("/", myHandler)

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

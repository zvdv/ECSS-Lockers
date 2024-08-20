package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/router/ioutil"
)

const htmlBase string = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width,initial-scale=1.0" />
        <link rel="stylesheet" type="text/css" href="assets/css/index.css" />
        <link rel="icon" type="image/x-icon" href="assets/favicon.png">
        <script src="assets/js/htmx.min.js" defer></script>
        <title>ECSS' Locker Registration</title>
    </head>

    <body>
        <!-- template goes in here -->
        %s 
    </body>

    <footer class="fixed bottom-0 flex justify-center items-center w-full p-5">
        <span class="">
            Having issues? Contact <a class="link link-info" target="blank_" href="mailto:foobar@uvic.ca">supportemail@email.ca<a/>.
        </span>
    </footer>
</html>
    `

func Base(writer io.Writer, t *template.Template, data any) error {
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		return err
	}

	html := fmt.Sprintf(htmlBase, buf.String())
	_, err := writer.Write([]byte(html))
	return err
}

func Html(w http.ResponseWriter, fileName string, data any) {
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		logger.Fatal("error reading file %s: %v", fileName, err)
		ioutil.WriteResponse(w, http.StatusInternalServerError, nil)
	}

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, data); err != nil {
		logger.Error("error executing template data: %v", err)
		ioutil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	html := fmt.Sprintf(htmlBase, buf.String())
    ioutil.WriteResponse(w, http.StatusOK, []byte(html))
}

func Component(writer http.ResponseWriter, fileName string, data any) {
	var err error

	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		logger.Fatal("error parsing template file: %v", err)
	}

	if err := tmpl.Execute(writer, data); err != nil {
		logger.Error("error writing template: %v", err)
		ioutil.WriteResponse(writer, http.StatusInternalServerError, nil)
	}
}

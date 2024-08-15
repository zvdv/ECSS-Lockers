package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

const htmlBase string = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <link rel="stylesheet" type="text/css" href="assets/css/index.css" />
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

func Html(writer io.Writer, fileName string, data any) error {
	var err error
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, data); err != nil {
		return err
	}

	html := fmt.Sprintf(htmlBase, buf.String())
	_, err = writer.Write([]byte(html))
	return err
}

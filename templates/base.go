package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

func Base(writer io.Writer, t *template.Template, data any) error {
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		return err
	}

	htmlString := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <link rel="stylesheet" type="text/css" href="assets/css/index.css" />
        <script src="assets/js/htmx.min.js" defer></script>
        <title></title>
    </head>

    <body>
        %s
    </body>

    <footer class="fixed bottom-0 flex justify-center items-center w-full p-5">
        <span class="">
            Having issues? Contact <a class="link link-info" target="blank_" href="mailto:foobar@uvic.ca">supportemail@email.ca<a/>.
        </span>
    </footer>
</html>
    `, buf.String())

	_, err := writer.Write([]byte(htmlString))
	return err
}

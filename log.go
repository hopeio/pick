package pick

import (
	"github.com/hopeio/utils/console/style"
	"log"

	stringsi "github.com/hopeio/utils/strings"
)

func Log(method, path, title string) {
	log.Printf(" %s\t %s %s\t %s",
		style.Green("API:"),
		style.Yellow(stringsi.FormatLen(method, 6)),
		style.Blue(stringsi.FormatLen(path, 50)), style.Magenta(title))
}

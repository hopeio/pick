package pick

import (
	"github.com/hopeio/tiga/utils/console/style"
	"github.com/hopeio/tiga/utils/log"
	stringsi "github.com/hopeio/tiga/utils/strings"
)

func Log(method, path, title string) {
	log.Printf(" %s\t %s %s\t %s",
		style.Green("API:"),
		style.Yellow(stringsi.FormatLen(method, 6)),
		style.Blue(stringsi.FormatLen(path, 50)), style.Purple(title))
}

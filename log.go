package pick

import (
	"fmt"
	"github.com/hopeio/lemon/utils/console/concolor"
	stringsi "github.com/hopeio/lemon/utils/strings"
)

func Log(method, path, title string) {
	fmt.Printf(" %s\t %s %s\t %s\n",
		concolor.Green("API:"),
		concolor.Yellow(stringsi.FormatLen(method, 6)),
		concolor.Blue(stringsi.FormatLen(path, 50)), concolor.Purple(title))
}

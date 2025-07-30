/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package pick

import (
	"github.com/hopeio/gox/terminal/style"
	"log"

	stringsi "github.com/hopeio/gox/strings"
)

func Log(method, path, title string) {
	log.Printf(" %s\t %s %s\t %s",
		style.Green("API:"),
		style.Yellow(stringsi.FormatLen(method, 6)),
		style.Blue(stringsi.FormatLen(path, 50)), style.Magenta(title))
}

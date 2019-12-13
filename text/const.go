package text

var eos map[rune]rune = map[rune]rune{
	'.': '.',
	'?': '?',
	'!': '!',
	'…': '.',
}

var sep map[rune]rune = map[rune]rune{
	'-':  '-',
	'–':  '-',
	'—':  '-',
	',':  ',',
	';':  ';',
	':':  ':',
	'"':  '"',
	'\'': '\'',
	'«':  '«',
	'»':  '»',
	'(':  '(',
	')':  ')',
}

var proto map[string]bool = map[string]bool{
	"ftp":   true,
	"http":  true,
	"https": true,
	"sftp":  true,
}

var numRoman map[string]bool = map[string]bool{
	"i":      true,
	"ii":     true,
	"iii":    true,
	"iv":     true,
	"v":      true,
	"vi":     true,
	"vii":    true,
	"viii":   true,
	"ix":     true,
	"x":      true,
	"xi":     true,
	"xii":    true,
	"xiii":   true,
	"xiv":    true,
	"xv":     true,
	"xvi":    true,
	"xvii":   true,
	"xviii":  true,
	"xix":    true,
	"xx":     true,
	"xxi":    true,
	"xxii":   true,
	"xxiii":  true,
	"xxiv":   true,
	"xxv":    true,
	"xxvi":   true,
	"xxvii":  true,
	"xxviii": true,
	"xxix":   true,
}

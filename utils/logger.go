package utils

import "fmt"

const (
	reset     = "\033[0m"
	Bold      = "\033[1m"
	Underline = "\033[4m"
	Strike    = "\033[9m"
	Italic    = "\033[3m"

	CRed    = "\033[31m"
	CGreen  = "\033[32m"
	CYellow = "\033[33m"
	CBlue   = "\033[34m"
	CPurple = "\033[35m"
	CCyan   = "\033[36m"
	CWhite  = "\033[37m"
)

func ColoredPrintln(msg string, color string) {
	fmt.Println(color + msg + reset)
}

func LogError(fileName string, err error, info string) {
	fmt.Println(CYellow + "[fileName]: \t" + fileName + reset)
	fmt.Println(CBlue + "[Info]: \t" + info + reset)
	fmt.Printf("%v[Error]:\t%v\n", CRed, err)
	fmt.Println(reset + "--------------------------------------------")
}

package gotokenize

import (
	"fmt"
)

var _fred = "\033[1;31m%s\033[0m"
var _fgreen = "\033[1;32m%s\033[0m"
var _fyellow = "\033[1;33m%s\033[0m"
var _fpurple = "\033[1;34m%s\033[0m"
var _fmagenta = "\033[1;35m%s\033[0m"
var _fteal = "\033[1;36m%s\033[0m"

func ColorRed(content string) string {

	return fmt.Sprintf(_fred, content)
}
func ColorGreen(content string) string {

	return fmt.Sprintf(_fgreen, content)
}
func ColorYellow(content string) string {

	return fmt.Sprintf(_fyellow, content)
}
func ColorPurple(content string) string {

	return fmt.Sprintf(_fpurple, content)
}
func ColorMagenta(content string) string {

	return fmt.Sprintf(_fmagenta, content)
}
func ColorTeal(content string) string {

	return fmt.Sprintf(_fteal, content)
}

//ColorOffset color the offset
func ColorOffset(offset int) string {

	return fmt.Sprintf("\033[1;32m%3d\033[0m", offset)
}

//ColorIgnore color ignore
func ColorIgnore() string {
	return fmt.Sprintf("\033[1;31mignore\033[0m")
}

//ColorName color name of token
func ColorName(name string) string {
	return fmt.Sprintf("\033[1;34m[%s]\033[0m", name)
}

//ColorType color type of token
func ColorType(tokenType int) string {
	return fmt.Sprintf("\033[1;35m%d\033[0m", tokenType)
}

//ColorContent color the content of token
func ColorContent(content string) string {
	return fmt.Sprintf("'\033[1;36m%s\033[0m'", content)
}

//ColorFail color fail
func ColorFail() string {
	return fmt.Sprintf("\033[1;31mfail\033[0m")
}

//ColorSuccess color success
func ColorSuccess() string {
	return fmt.Sprintf("\033[1;32msuccess\033[0m")
}

//Log to log
type Log struct {
	logs string
}

//Append append log
func (log *Log) Append(content string) {
	log.logs += content
}

//Print print
func (log *Log) Print() {
	fmt.Println(log.logs)
}

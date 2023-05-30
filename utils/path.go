package utils

func CheckPath(path string) {
	if string(path[0]) != "/" {
		PanicWithContext("路径必须以“/”开头")
	}
	if string(path[len(path)-1]) == "/" {
		PanicWithContext("\"路径不能以“/”结束\"")
	}
}

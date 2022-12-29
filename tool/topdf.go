package tool

import (
	"fmt"
	"os/exec"
)

// word转pdf
func WordToPdf(resPathName, outPathName string) error {
	command := exec.Command("java", "-Dsun.jnu.encoding=UTF-8", "-Dfile.encoding=UTF-8", "-jar", "./resources/w2f.jar", resPathName, outPathName)
	fmt.Println(command.String())
	out, err := command.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

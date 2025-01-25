package covplots

import (
	"os"
	"os/exec"
)

func PlotSimple(scriptPath string) func(outpre string) error {
	return func(outpre string) error {
		data := outpre + "_plfmt.bed"
		plotted := outpre + ".pdf"
		cmd := exec.Command(scriptPath, data, plotted)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

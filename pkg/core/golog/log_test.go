package golog

import "testing"

func TestLog(t *testing.T) {
	err := InitZapLog("info", "stdout")
	if err != nil {
		t.Fatal("initZapLog failed", err)
	}

	defer logger.Sync()
	logOutput := []string{"Debug", "log"}
	Debug("output", logOutput)

	logOutput = []string{"Error", "log"}
	Error("output", logOutput)

	logOutput = []string{"Warn", "log"}
	Warn("output", logOutput)

	logOutput = []string{"Info", "log"}
	Info("output", logOutput)

	logOutput = []string{"Fatal", "log"}
	Fatal("output", logOutput)

}

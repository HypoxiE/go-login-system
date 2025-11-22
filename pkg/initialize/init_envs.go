package initialize

import "os"

func InitEnv(enviroment map[string]string) {
	for name, path := range enviroment {
		os.Setenv(name, path)
	}
	os.Setenv("SHELL", "/bin/bash")
}

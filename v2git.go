/*

 */

package main

import (
	"fmt"
	"os"

	"github.com/molisoft/v2git/api"
	"github.com/molisoft/v2git/auth"
	"github.com/molisoft/v2git/cfg"
	"github.com/molisoft/v2git/git"
	"github.com/molisoft/v2git/repo"
)

var config *cfg.Config

func init() {
	var err error
	config, err = cfg.New("config.json")
	if err != nil {
		fmt.Println("Failed to load config.json. ", err.Error())
		os.Exit(1)
	}

	repo.RootDir = config.RepoDir
}

func apiServer() {
	go api.Start(":3344", config.ApiPassword)
}

func verify(request *cfg.AuthRequest) bool {
	b, err := auth.Verify(request)
	if err != nil {
		fmt.Println("Verify error ", err)
		return false
	}
	return b
}

func main() {
	httpAddr := fmt.Sprintf(":%d", config.HttpPort)
	sshAddr := fmt.Sprintf(":%d", config.SshPort)

	apiServer()

	ssh := git.NewSSH(sshAddr, config)
	ssh.RegisterVerify(verify)
	go ssh.Start()

	http := git.NewHTTP(httpAddr, config)
	http.RegisterVerify(verify)
	http.Start()

	//select {}
}

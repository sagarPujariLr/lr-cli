package logout

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	"github.com/loginradius/lr-cli/cmdutil"
	"github.com/loginradius/lr-cli/config"
	"github.com/spf13/cobra"
)

var tempServer *cmdutil.TempServer

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout of LR account",
		Long:  `This commmand logs user out of the LR account`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logout()
		},
	}
	return cmd
}

func logout() error {
	conf := config.GetInstance()
	user, _ := user.Current()
	dirName := filepath.Join(user.HomeDir, ".lrcli")
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		fmt.Println("You are already been logged out")
		return nil
	} else {
		cmdutil.Openbrowser(conf.HubPageDomain + "/auth.aspx?action=logout&return_url=http://localhost:8089/postLogout")
		tempServer = cmdutil.CreateTempServer(cmdutil.TempServer{
			Port:        ":8089",
			HandlerFunc: postLogout,
			RouteName:   "/postLogout",
		})
		tempServer.Server.ListenAndServe()
		dir, err := ioutil.ReadDir(dirName)
		for _, d := range dir {
			os.RemoveAll(path.Join([]string{dirName, d.Name()}...))
		}
		if err != nil {
			return err
		}

	}
	log.Println("You are successfully Logged Out")
	return nil
}

func postLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are Successfully Logged Out, Kindly Close this browser window and go back to CLI")
	time.AfterFunc(1*time.Second, tempServer.CloseServer)
}
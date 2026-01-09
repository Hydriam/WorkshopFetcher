package libLWD

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/codeclysm/extract/v4"
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
)

func GetSteamcmd() error {
	//Thanks for https://gist.github.com/cnu/026744b1e86c6d9e22313d06cba4c2e9

	//Download the archive with steamcmd binaries
	out, err := os.Create("steamcmd.tar.gz")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get("https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	_, err = os.Stat("steamcmd")
	if os.IsNotExist(err) {
		err = os.Mkdir("steamcmd", 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.Open("steamcmd.tar.gz")
	if err != nil {
		return err
	}
	err = extract.Gz(context.TODO(), file, "steamcmd", nil)
	if err != nil {
		defer file.Close()
		defer os.Remove("steamcmd.tar.gz")
		return err
	}
	defer file.Close()
	os.Remove("steamcmd.tar.gz")
	adw.NewAboutDialog()
	return nil
}
func CheckSteamcmd() bool {
	// This function checks if steamcmd is downloaded
	_, err := os.Stat("steamcmd/steamcmd.sh")
	return err == nil
}

func DownloadFromSteamcmd(appID string, workshopIDs []string) error {
	//TODO: proper logging, save all the output to a log file, if there is a error just display a adw alert dialog
	logFile, err := os.Create("steamcmd.log") // Create or overwrite the log file
	if err != nil {
		return err
	}
	defer logFile.Close()
	// Start logging
	logger := log.New(logFile, "", log.LstdFlags)
	logger.Println("=== SteamCMD Download Started ===")
	logger.Printf("App ID: %s\n", appID)
	logger.Printf("Workshop IDs: %v\n", workshopIDs)

	logger.Println("Download is starting:")
	// cmdt = command template
	cmdt := []string{
		"./steamcmd/steamcmd.sh",
		"+login", "anonymous",
	}
	for i := 0; i < len(workshopIDs); i++ {
		cmdt = append(cmdt, "+workshop_download_item", appID, workshopIDs[i])
	}
	cmdt = append(cmdt, "+quit")
	cmd := exec.Command(cmdt[0], cmdt[1:]...) //cmdt[1:]... makes it use every object in array
	// Redirect output of the command to log file
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	// Run the command
	err = cmd.Run()
	if err != nil {
		logger.Println("Steamcmd Failed.")
		logger.Println("Make sure that you have 32bit version of glibc installed on your system.")
		return err
	}
	logger.Println("Please check above lines to know if everything went succesfully.")
	logger.Println("The file should be downloaded.")
	logger.Println("It should be under ~/.local/share/Steam/steamapps/workshop/")
	return nil
}

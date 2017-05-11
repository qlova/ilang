package main

import "os"
import "fmt"
import "io"
import "os/exec"
import "net/http"
import "bufio"

var git = "git.exe"
var goc = "go.exe"
var ext = ".exe"

func downloadFile(filepath string, url string) (err error) {

  // Create the file
  out, err := os.Create(filepath)
  if err != nil  {
    return err
  }
  defer out.Close()

 fmt.Println("Downloading...")
  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  // Writer the body to file
  _, err = io.Copy(out, resp.Body)
  if err != nil  {
    return err
  }

  return nil
}

func SystemChecks() {
	var err error

	err = exec.Command("go.exe", "version").Run()
	if err != nil {
		err = exec.Command(`C:\Go\bin\go.exe`, "version").Run()
		if err != nil {
			fmt.Println("The Go programming language is required for IT and will now be downloaded... Please Wait =D")
			err = downloadFile("./InstallGo1.8.msi", "https://storage.googleapis.com/golang/go1.8.1.windows-amd64.msi")
			if err != nil {
				fmt.Println("The Go programming language cannot be automatically installed on your system.")
				fmt.Println("Please visit http://golang.org/dl/ and install it manually.")
				
				fmt.Println("Press enter to close.")
				reader := bufio.NewReader(os.Stdin)
				reader.ReadString('\n')
				os.Exit(1)
			}
			err = exec.Command("msiexec", "/i", "%s", "/qn", "./InstallGo1.7.msi").Run()
			if err != nil {
				fmt.Println("Go has been downloaded, please run InstallGo1.7.msi")
				
				fmt.Println("Press enter to close.")
				reader := bufio.NewReader(os.Stdin)
				reader.ReadString('\n')
				os.Exit(2)
			}
			fmt.Println("Please rerun IT after you have completed the installation of Go :)")
			
			fmt.Println("Press enter to close.")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
			os.Exit(2)
		} else {
			goc = `C:\Go\bin\go.exe`
		}
	}
	
	/*err = exec.Command("git.exe", "--version").Run()
	if err != nil {
		fmt.Println("Git is required for IT in order to function and will now be downloaded... Please Wait =D")
		err = downloadFile("./InstallGit.exe", "http://github.com/git-for-windows/git/releases/download/v2.10.1.windows.1/Git-2.10.1-32-bit.exe")
		if err != nil {
			fmt.Println("Git cannot be automatically installed on your system.")
			fmt.Println("Please visit https://github.com/git-for-windows/git/releases/tag/v2.10.1.windows.1 and install it manually.")
			os.Exit(1)
		}
		err = exec.Command("./InstallGit.exe").Run()
		if err != nil {
			fmt.Println("Git has been downloaded, please run InstallGit.exe")
			os.Exit(2)
		}
		fmt.Println("Please rerun IT after you have completed the installation of Git :)")
		os.Exit(2)
	}*/
}

package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
)

const (
    Red    = "\033[0;31m"
    Green  = "\033[0;32m"
    Yellow = "\033[0;33m"
    Cyan   = "\033[0;36m"
    NC     = "\033[0m"
)

var (
    destinationDir = os.ExpandEnv("$HOME/Zone01_Desk_cfg")
    downloadLink   = "https://github.com/AmineS530/Config-maker.git"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: config-maker [setup|docker|bash|zsh|git|theme|background|finish]")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "setup":
        runSetup()
    case "docker":
        runDocker()
    case "bash":
        disableZsh()
    case "zsh":
        enableZsh()
    case "git":
        runScript("git_setup.sh")
    case "theme":
        runScript("set_theme.sh")
    case "background":
        runScript("set_background.sh")
    case "finish":
        // runFinish()
    default:
        fmt.Println("Unknown command:", os.Args[1])
    }
}

func runSetup() {
    if _, err := os.Stat(os.ExpandEnv("$HOME/.oh-my-zsh")); err == nil {
        fmt.Println(Green + "Oh-my-zsh is already installed, continuing..." + NC)
    } else {
        fmt.Println(Yellow + "Installing Oh-my-zsh..." + NC)
        exec.Command("sh", "-c", "wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh -O /tmp/install.sh && sh /tmp/install.sh --unattended").Run()
    }

    if _, err := os.Stat(destinationDir); err == nil {
        fmt.Println(Yellow+"Directory already exists. Overwriting..."+NC)
        os.RemoveAll(destinationDir)
    } else {
        fmt.Println(Cyan+"Cloning repository..."+NC)
    }

    if err := exec.Command("git", "clone", "--quiet", "--depth=1", downloadLink, destinationDir).Run(); err != nil {
        log.Fatal(err)
    }

    runScript("setup.sh")
}

func runScript(name string) {
    script := destinationDir + "/" + name
    cmd := exec.Command("zsh", script)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }
}

func runDocker() {
    fmt.Println(Green + "Getting Docker ready for first use..." + NC)
    exec.Command("sh", "-c", "curl -fsSL https://get.docker.com/rootless | sh").Run()
    fmt.Println(Green + "Docker environment is set up for rootless mode." + NC)
}

func disableZsh() {
    fmt.Println(Green + "Removing zsh from default shell" + NC)
    exec.Command("sed", "-i", "/SHELL=\\/bin\\/zsh/d", os.ExpandEnv("$HOME/.bashrc")).Run()
    exec.Command("sed", "-i", "/exec \\/bin\\/zsh -l/d", os.ExpandEnv("$HOME/.bashrc")).Run()
    fmt.Println(Green + "Enjoy Bash :))" + NC)
}

func enableZsh() {
    fmt.Println(Green + "Enabling zsh as default shell" + NC)
    f, err := os.OpenFile(os.ExpandEnv("$HOME/.bashrc"), os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    f.WriteString("SHELL=/bin/zsh\nexec /bin/zsh -l\n")
}
// ai generated for ref 
//todo later
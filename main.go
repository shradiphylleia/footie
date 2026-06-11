package main
import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Match struct {
	Active bool `json:"active"`
	You int `json:"you"`
	Docker int `json:"docker"`
	StartedAt time.Time `json:"started_at"`
	EndedAt time.Time `json:"ended_at,omitempty"`
	Commands []PlayedCommand `json:"commands"`
}

type PlayedCommand struct {
	Command string `json:"command"`
	ExitCode int `json:"exit_code"`
	WhoScored string `json:"who_scored"`
	PlayedAt time.Time `json:"played_at"`
}

var white="\033[97m"
var yellow="\033[33m"
var green="\033[32m"
var cyan="\033[36m"
var reset="\033[0m"

func main() {
	args:=os.Args[1:]

	if len(args)==0{
		help()
		return
	}
	command:= args[0]

	if command=="kickoff" {
		kickoff()
		return
	}
	if command=="docker" {
		runDocker(args[1:])
		return
	}
	if command=="hook" {
		hook(args[1:])
		return
	}

	if command=="score" {
		showScore()
		return
	}

	if command=="fulltime" {
		fulltime()
		return
	}

	if command=="help"||command=="--help"||command=="-h"{
		help()
		return
	}
	// if it is not a footie command, try it as a docker command
	runDocker(args)
}
func help() {
	fmt.Println("footie-docker football in your terminal")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println(" footie kickoff")
	fmt.Println(" footie hook bash")
	fmt.Println(" footie hook powershell")
	fmt.Println(" footie hook cmd")
	fmt.Println(" footie docker <docker args...>")
	fmt.Println(" footie compose up -d")
	fmt.Println(" footie build -t my-app .")
	fmt.Println(" footie score")
	fmt.Println(" footie fulltime")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("  footie kickoff")
	fmt.Println("  eval \"$(footie hook bash)\"")
	fmt.Println("  docker compose up -d")
	fmt.Println("  footie compose up -d")
	fmt.Println("  footie fulltime")
}
func kickoff() {
	oldMatch, err := loadMatch()
	if err == nil && oldMatch.Active {
		fmt.Println("Match already started.")
		printScore(oldMatch)
		return
	}
	match := Match{}
	match.Active = true
	match.You = 0
	match.Docker = 0
	match.StartedAt = time.Now()
	match.Commands = []PlayedCommand{}
	err = saveMatch(match)
	if err != nil {
		fmt.Println("Could not save match:", err)
		os.Exit(1)
	}
	fmt.Println("Kickoff!")
	fmt.Println("you vs docker")
	printScore(match)
	fmt.Println("")
	fmt.Println("hook docker in your terminal:")
	printHookHelp()
}

func printHookHelp(){
	exePath, err:=os.Executable()
	if err!=nil{
		fmt.Println("  Git Bash / bash / zsh:  eval \"$(footie hook bash)\"")
		fmt.Println("  PowerShell:  Invoke-Expression (& footie hook powershell)")
		fmt.Println("  Cmd: footie hook cmd")
		return
	}
	bashExe:=bashPath(exePath)
	powerShellExe:=strings.ReplaceAll(exePath, "'", "''")
	fmt.Println("  Git Bash / bash / zsh:  eval \"$(" + bashExe + " hook bash)\"")
	fmt.Println("  PowerShell:  Invoke-Expression (& '" + powerShellExe + "' hook powershell)")
	fmt.Println("  Cmd: \"" + exePath + "\" hook cmd")
}
func hook(args []string) {
	if len(args)==0{
		fmt.Println("Usage: footie hook bash|zsh|powershell|cmd")
		os.Exit(1)
	}
	exePath, err := os.Executable()
	if err!=nil{
		exePath="footie"
	}
	shell:= args[0]
	if shell=="bash"||shell=="zsh"{
		printBashHook(exePath)
		return
	}
	if shell=="powershell"||shell=="pwsh"{
		printPowerShellHook(exePath)
		return
	}
	if shell == "cmd" {
		printCmdHook(exePath)
		return
	}
	fmt.Println("I don't know this shell:", shell)
	fmt.Println("Use: footie hook bash|zsh|powershell|cmd")
	os.Exit(1)
}
func printBashHook(exePath string) {
	exePath = bashPath(exePath)
	exePath = strings.ReplaceAll(exePath, "'", "'\"'\"'")
	fmt.Println("docker() {")
	fmt.Println("    '" + exePath + "' docker \"$@\"")
	fmt.Println("}")
}
func bashPath(exePath string) string {
	exePath = strings.ReplaceAll(exePath, "\\", "/")
	if len(exePath) > 2 && exePath[1] == ':' {
		drive := strings.ToLower(exePath[0:1])
		return "/" + drive + exePath[2:]
	}
	return exePath
}
func printPowerShellHook(exePath string) {
	exePath = strings.ReplaceAll(exePath, "'", "''")
	fmt.Println("function global:docker {")
	fmt.Println("    & '" + exePath + "' docker @args")
	fmt.Println("}")
}
func printCmdHook(exePath string) {
	fmt.Println("doskey docker=\"" + exePath + "\" docker $*")
}
func runDocker(dockerArgs []string) {
	if len(dockerArgs) == 0 {
		fmt.Println("usage: footie docker <docker args...>")
		os.Exit(1)
	}
	match,err:=loadMatch()
	if err!= nil || match.Active== false {
		exitCode:= runPlainDocker(dockerArgs)
		os.Exit(exitCode)
	}
	exitCode:=runPlainDocker(dockerArgs)
	played:=PlayedCommand{}
	played.Command="docker " + strings.Join(dockerArgs, " ")
	played.ExitCode=exitCode
	played.PlayedAt=time.Now()
	if exitCode==0{
		match.You=match.You+1
		played.WhoScored="you"
		goalAnimation("you score!")
	} else {
		match.Docker = match.Docker + 1
		played.WhoScored = "docker"
		goalAnimation("docker scores")
	}
	match.Commands = append(match.Commands, played)
	err = saveMatch(match)
	if err != nil {
		fmt.Println("could not save score:", err)
		os.Exit(1)
	}
	fmt.Println("")
	printScore(match)
	os.Exit(exitCode)
}

func runPlainDocker(dockerArgs []string) int {
	dockerCommand := exec.Command("docker", dockerArgs...)
	dockerCommand.Stdin = os.Stdin
	dockerCommand.Stdout = os.Stdout
	dockerCommand.Stderr = os.Stderr

	err:=dockerCommand.Run()
	if err == nil {
		return 0
	}
	var dockerError *exec.ExitError
	if errors.As(err, &dockerError) {
		return dockerError.ExitCode()
	}
	fmt.Println(err)
	return 1
}

func showScore(){
	match, err:=loadMatch()
	if err != nil {
		fmt.Println("no match/game found. run `footie kickoff` to get started")
		os.Exit(1)
	}
	printScore(match)
}
func fulltime(){
	match, err:=loadMatch()
	if err != nil || match.Active == false {
		fmt.Println("no active match/game found.")
		os.Exit(1)
	}
	match.Active=false
	match.EndedAt=time.Now()
	err=saveMatch(match)
	if err != nil {
		fmt.Println("could not save match:", err)
		os.Exit(1)
	}
	// only if i knew how to add sound ?
	fmt.Println("full time!")
	printScore(match)
	fmt.Println("docker commands played:",len(match.Commands))
}
func printScore(match Match) {
	fmt.Println("you", match.You, "-", match.Docker, "docker")
}

func matchFile() string {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "footie-session.json"
	}
	return filepath.Join(cache, "footie", "session.json")
}

func loadMatch() (Match, error) {
	fileName := matchFile()
	data, err := os.ReadFile(fileName)
	if err != nil {
		return Match{}, err
	}
	match := Match{}
	err = json.Unmarshal(data, &match)
	return match, err
}

func saveMatch(match Match) error {
	fileName:= matchFile()
	folderName:= filepath.Dir(fileName)
	err := os.MkdirAll(folderName, 0755)
	if err!=nil {
		return err
	}
	data, err:=json.MarshalIndent(match, "", "  ")
	if err!=nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func goalAnimation(words string) {
	fmt.Println("")
	fmt.Print("\033[?25l")
	fmt.Print("\033[s")

	fmt.Println(white + "     ___________________________________")
	fmt.Println("       | . . . . . . . . . . . . . . . . |")
	fmt.Println("       | . . . . . . . . . . . . . . . . |")
	fmt.Println("       | . . . . . . . . . . . . . . . . |")
	fmt.Println("       | . . . . . . . . . . . . . . . . |")
	fmt.Println("       | . . . . . . . . . . . . . . . . |")
	fmt.Println("       |_________________________________|" + reset)
	fmt.Println("                    ")
	fmt.Println("                    ")
	ballPath := [][]int{
		{8, 24},
		{7, 24},
		{6, 24},
		{5, 24},
		{4, 24},
		{3, 24},
		{2, 24},
	}
	oldRow:=-1
	oldColumn:=-1
	for i:=0;i<len(ballPath); i++ {
		if oldRow!= -1 {
			moveCursor(oldRow, oldColumn)
			fmt.Print(" ")
		}
		row := ballPath[i][0]
		column := ballPath[i][1]
		moveCursor(row, column)
		fmt.Print(yellow + "O" + reset)
		oldRow = row
		oldColumn = column
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(150 * time.Millisecond)
	moveCursor(2, 20)
	fmt.Print(cyan + "~ ~ O ~ ~" + reset)
	moveCursor(4, 17)
	fmt.Print(green + words + reset)
	moveCursor(10, 1)
	fmt.Print("\033[?25h")
	fmt.Println("")
}

func moveCursor(row int, column int) {
	fmt.Print("\033[u")
	if row>0{
		fmt.Print("\033["+fmt.Sprint(row)+"B")
	}
	if column>0{
		fmt.Print("\033["+fmt.Sprint(column)+"C")
	}
}

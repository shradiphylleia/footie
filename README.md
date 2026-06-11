# footie
docker but make it a small football match in the terminal.
every docker command gets treated like a shot.
if docker exits with `0`, you score.
if docker exits with anything else, docker scores.
then footie shows a little terminal goal.

## demo:
https://github.com/user-attachments/assets/e551b850-6d30-4a50-a1ff-fd0377095e88

## build
```bash
go build -o footie.exe .
```

on mac/linux it would/should be:

```bash
go build -o footie .
```

## start
start the match:
```bash
./footie.exe kickoff
```
then hook docker for this terminal.
Git Bash / bash / zsh:

```bash
eval "$(./footie.exe hook bash)"
```

PowerShell:
```powershell
Invoke-Expression (& .\footie.exe hook powershell)
```

Cmd:
```cmd
footie.exe hook cmd
```
after that, just use docker like normal:

```bash
docker --version
docker compose up -d
docker build -t my-app .
```

## score
```bash
./footie.exe score
```

## end
```bash
./footie.exe fulltime
```

## without the hook
this also works:

```bash
./footie.exe docker --version
./footie.exe docker compose up -d
```
the hook is just for the nicer version where you can keep typing `docker ...`.


## notes
this is the first version.
it is intentionally simple right now.
also: `footie kickoff` starts the match, but the hook is still needed because a cli cannot change your already-running terminal by itself.

i wrote a small build note here: [learning.md](learning.md)

github release might follow ;)
happy football season!

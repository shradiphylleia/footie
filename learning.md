# how
i wanted a tiny cli where docker commands feel like a football match.
```text
docker command works  -> i score
docker command fails  -> docker scores
```
then after the command, the terminal shows a small goal animation.

## first idea
the first version was basically:
```bash
footie kickoff
footie docker --version
footie score
footie fulltime
```
`footie docker ...` runs the real docker command under the hood.
so this:
```bash
footie docker compose up -d
```
becomes this inside the go code:

```bash
docker compose up -d
```
then footie checks the exit code.
## the score

the score is saved as a small json file in the user cache folder.
it stores things like:

```text
match active?
what is my score?
what is docker's score?
which docker commands were played?
```
so the match can keep going across multiple commands

## the shell problem :/
this was the main thing i had to figure out.
i wanted this:

```bash
footie kickoff
docker compose up -d
```
but a normal cli cannot run once and then magically change how your terminal works.
when `footie kickoff` exits, it cannot rename or intercept `docker` in the parent shell.
so the solution is a shell hook.

after kickoff, you run one hook command for your terminal:
Git Bash / bash / zsh:
```bash
eval "$(footie hook bash)"
```
PowerShell:
```powershell
Invoke-Expression (& footie hook powershell)
```
that creates a temporary `docker` function in the current terminal session.
so when i type:
```bash
docker --version
```
the terminal actually calls:
```bash
footie docker --version
```
and footie still runs the real docker command after that.


## what works now

right now footie can:

```text
start a match
hook docker in bash/zsh/powershell/cmd
run docker commands
score based on exit code
show a goal animation
show the score
end the match
```

## current flow
for me in Git Bash:
```bash
go build -o footie.exe .
./footie.exe kickoff
eval "$(./footie.exe hook bash)"
docker --version
./footie.exe score
./footie.exe fulltime
```

happy football season 2026 :)
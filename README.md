# encabulator

Some golang projects, grouped together.

[About the name](https://www.youtube.com/watch?v=RXJKdh1KZ0w).

## rekey

```bash
go get -u github.com/justjake/encabulator/cmd/rekey
```

This is a tool that ensures your Yubikey is loaded into ssh-agent. If it is
not, rekey will load the key. This is useful if you use Yubikey for your Github
key, and freqently get annoying messages when trying to push.

Add this to your `~/.gitconfig` and forget about running `ssh-add` manually:

```
[alias]
pll = !rekey && git pull
psh = !rekey && git push
ftch = !rekey && git fetch
```

You can also alias `ssh` to ensure your keys are loaded:

```bash
alias ssh="reke && ssh"
```

TODOs for the command-line tool:

- add command line parsing
- add help output
- add support for goldkey
- add support for loading the default identity that `ssh-add` would add

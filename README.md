# encabulator

Some golang projects, grouped together.

[About the name](https://www.youtube.com/watch?v=RXJKdh1KZ0w).

## rekey

```bash
go get -u github.com/justjake/encabulator/cmd/rekey
```

This is a tool that ensures your Yubikey is laoded into ssh-agent. If it is
not, rekey will laod the key. This is useful if you use Yubikey for your Github
key, and freqently get annoying messages when trying to push.

Add this to your `~/.gitconfig` and forget about running `ssh-add` manually:

```
[alias]
pll = !rekey && git pull
psh = !rekey && git push
ftch = !rekey && git fetch
```

TODOs:

- add command line parsing
- add help output
- add support for goldkey
- add support for 

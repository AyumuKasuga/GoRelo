GoRelo
======

Tool for development which automatically reloads your service when some of files changed.

Inspired by [nodemon](https://nodemon.io/) and Django's `python manage.py runserver` but language agnostic.

Example:
```
gorelo -e ".git .vscode" python service.py
```

By default it watches all current directory, most likely you will want to exclude some subdirectories using `-e` parameter.


[![asciicast](https://asciinema.org/a/150624.png)](https://asciinema.org/a/150624)


Feel free to create issue or make PR.

## Build
```
go build
```
# TLDR
LSP server/client for time sheet management.
Each day is new file.
Inspired by [lsp video](https://www.youtube.com/watch?v=YsdlcQoHqPY&t=6296s&pp=ygUDbHNw).

# Using:
- [x] [LSP](https://microsoft.github.io/language-server-protocol/)
- [x] [sqlite3](https://www.sqlite.org/index.html)
- [x] [go migrate]( https://github.com/golang-migrate/migrate)
- [x] [sqlc](https://sqlc.dev/)

# Features:
- [x] Show hover shows the summary of the day.
- [x] Go to definition show daily sorted statistics, weekly and monthly for day represented by file.
- [x] Completion for category filling.
- [x] Colorize category,time and description.

# Example usage
Use 
```
make create_test_project
```
which create and open example project in [neovim](https://neovim.io/)



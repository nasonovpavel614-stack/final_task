package tests

// Port — порт запущенного сервера (должен совпадать с TODO_PORT).
var Port = 7540

// DBFile — путь к SQLite относительно папки tests.
var DBFile = "../scheduler.db"

// FullNextDate = true — тестировать правила повторения w и m.
var FullNextDate = true

// Search = true — тестировать поиск в GET /api/tasks?search=...
var Search = true

// Token — JWT из /api/signin (кука token), если сервер запущен с TODO_PASSWORD.
var Token = ``

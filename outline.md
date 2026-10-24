# Refactor outline: `src_old` → `src`

## Understanding of the plan

- `src` is a rewrite of `src_old` with stricter layering: **web → service → repositories**.
- Dependencies are wired by **constructor injection** in `main.go`: services are constructed there (`service.NewXService()`) and passed into `web.Initrouter(...)`, which hands them to handler factories.
- **Everything gin lives in `web`.** `service` must never import gin, and must not know about `*gin.Context`, sessions, status codes or HTML. It gets plain Go values (`context.Context`, input structs) and returns plain values plus an `error`.
- Handlers follow the pattern on `src/web/routes.go:11`: a function that takes the service(s) it needs and returns a `gin.HandlerFunc` closure.

  ```go
  r.POST("/api/register", handleApiRegister(authService))
  ```

- **Every handler is a factory**, also those without a dependency (they just take no parameters: `r.GET("/about", handleAboutPage())`). This keeps the codebase uniform.
- **Business logic** = all logic that neither depends on the web framework nor talks to the database directly. It belongs in `service`. For example, checking a username's length or that it has no spaces is business logic, and so is comparing a password with the stored hash. The handler only extracts raw values from the request into an input struct and maps the service's result/error to an HTTP response.
- **Scope for now:** structure and boilerplate only. No real service logic, no template rendering, no sessions, no static files, no DB. Every service method just returns the name of its route, e.g. `"POST /api/register"`.
- **Naming** follows Go's initialism convention: `InitRouter`, `handleAPIRegister`, not `Initrouter` / `handleApiRegister`.

## Routes in `src_old` (all of them need to exist in `src`)

From `src_old/main.go:63-87`:

| Method | Path            | Old handler                     | Old dependency    | Kind |
|--------|-----------------|---------------------------------|-------------------|------|
| GET    | `/`             | `delivery.SearchPage(db)`       | DB (search pages) | page |
| GET    | `/login`        | `delivery.LoginPage`            | none              | page |
| GET    | `/register`     | `delivery.RegisterPage`         | none              | page |
| GET    | `/about`        | `delivery.AboutPage`            | none              | page |
| GET    | `/weather`      | `delivery.WeatherPage`          | weather client    | page |
| POST   | `/api/login`    | `loginHandler.APILogin`         | DB + session      | api  |
| POST   | `/api/register` | `registerHandler.APIRegister`   | DB                | api  |
| GET    | `/api/logout`   | `delivery.LogoutHandler`        | session only      | api  |
| GET    | `/api/search`   | `delivery.SearchAPI(db)`        | DB (search pages) | api  |
| GET    | `/api/weather`  | `delivery.WeatherAPI`           | weather client    | api  |
| GET    | `/health`       | `delivery.HealthHandler`        | none              | api  |
| GET    | `/static/*`     | `router.Static(...)`            | filesystem        | infra |

Old infrastructure set up in `main.go` that is **not** a route: HTML renderer (`multitemplate`), cookie session store + `sessions.Sessions` middleware, DB open/init.

## Proposed structure in `src`

```
src/
├── main.go                 construct services → web.Initrouter(...) → Run
├── web/                    gin only: routing, request parsing, responses
│   ├── routes.go           InitRouter: registers every route
│   ├── auth.go             handleLoginPage, handleRegisterPage,
│   │                       handleAPILogin, handleAPIRegister, handleAPILogout
│   ├── search.go           handleSearchPage, handleAPISearch
│   ├── weather.go          handleWeatherPage, handleAPIWeather   (new file)
│   ├── about.go            handleAboutPage                       (new file)
│   └── health.go           handleHealth                          (new file)
├── service/                no gin imports
│   ├── auth.go             AuthService   (exists; extend)
│   ├── search.go           SearchService                         (new file)
│   └── weather.go          WeatherService                        (new file)
├── repositories/           untouched in this phase
└── models/                 untouched in this phase
```

### Service ↔ route mapping (first pass)

| Route                | Handler              | Service method                      |
|----------------------|----------------------|-------------------------------------|
| `GET /`              | `handleSearchPage`   | `SearchService.SearchPage`          |
| `GET /api/search`    | `handleAPISearch`    | `SearchService.Search`              |
| `GET /weather`       | `handleWeatherPage`  | `WeatherService.WeatherPage`        |
| `GET /api/weather`   | `handleAPIWeather`   | `WeatherService.Forecast`           |
| `POST /api/login`    | `handleAPILogin`     | `AuthService.Login`                 |
| `POST /api/register` | `handleAPIRegister`  | `AuthService.Register` (exists)     |
| `GET /api/logout`    | `handleAPILogout`    | `AuthService.Logout`                |
| `GET /login`         | `handleLoginPage`    | — (pure page, no service)           |
| `GET /register`      | `handleRegisterPage` | — (pure page, no service)           |
| `GET /about`         | `handleAboutPage`    | — (pure page, no service)           |
| `GET /health`        | `handleHealth`       | — (web-level liveness check)        |
| `GET /static/*`      | `r.Static`           | — (deferred until templates move)   |

Handlers without a service return their own route name directly from the handler.

### Order of work

- **Step 4:** register every route in `routes.go`. Every handler is a parameterless factory, except `handleAPIRegister`, which already takes `authService`. Handlers that will get a service in step 6 respond `"OK"` for now. Service-less handlers respond with their route name, which is final for this phase.
- **Step 6:** add `SearchService` and `WeatherService`, extend `AuthService`, add their parameters to the handler factories and `InitRouter`, and construct them in `main.go`.
- **Step 7:** every service method returns `(string, error)` with its route name. Handlers write that string with `c.String`.

### Wiring

```go
// main.go
authService := service.NewAuthService()
searchService := service.NewSearchService()
weatherService := service.NewWeatherService()

router := web.InitRouter(authService, searchService, weatherService)
```

### Placeholder service methods (step 7)

```go
func (s *SearchService) Search(ctx context.Context, input SearchInput) (string, error) {
	return "GET /api/search", nil
}
```

The handler calls the service, handles the error, and writes the returned string with `c.String(http.StatusOK, result)`.

## Things noticed along the way

- `src/models/Model.go` is not valid Go (no `package` clause). `go build .` still succeeds because nothing imports `models`, but `go vet ./...` and the lint workflow will fail on it.
- `src/web/auth.go:handleApiLogin` and `src/web/search.go` handlers are plain functions, not factories, and aren't registered yet. They'll be converted to the `routes.go:11` pattern.
- Old `/api/logout` only touches the session, so it has no business logic today. It still gets `AuthService.Logout` as an agreed placeholder.

## Decisions (answers to step 2)

1. Pages and `/health` stay web-only. `AuthService.Logout` gets a method.
2. One service method per route for now.
3. Service methods return `(string, error)`. `Register` changes accordingly.
4. `InitRouter` takes services as separate parameters.
5. Static files and sessions wait until templates move. Names are changed to Go's initialism convention.

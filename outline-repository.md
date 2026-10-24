# Refactor outline: repository layer

## Understanding of the plan

- Add the third layer: **web → service → repositories**. The repository layer is the *only* place that talks to the database.
- Same wiring style as before: **constructor injection in `main.go`**. `main.go` opens the DB, builds the repositories with it, builds the services with the repositories, and builds the router with the services.
- The DB file stays at `/data/whoknows.db` (repo root, git-ignored via `*.db`). The path comes from `DATABASE_PATH` in `.env`, with `../data/whoknows.db` as the fallback, like in `src_old`. The path is relative to the working directory, so the server must be started from `src/`.
- Repositories hold **no business logic**: no validation, no hashing, no defaults. They run SQL and map rows to and from model structs.
- `service` must not need to import `database/sql`. Repositories translate driver errors into the layer's own errors (see question 3).

## What exists in `src_old`

`src_old/internal/repository/`:

| File                | Contents |
|---------------------|----------|
| `database.go`       | `Open(path)`: creates the dir, opens SQLite (`modernc.org/sqlite`), `MaxOpenConns(1)`, ping, `PRAGMA foreign_keys = ON`. `Initialize(ctx, db)`: runs the embedded `schema.sql`. |
| `schema.sql`        | Tables `users (id, username UNIQUE, email UNIQUE, password_hash)` and `pages (title PK, url UNIQUE, language 'en'/'da', last_updated, content)`. |
| `userRepository.go` | `CreateUser`, `FindUserByUsername`, `IsDuplicateUserError` |
| `pageRepository.go` | `CreatePage`, `SearchPages(query, language)` |

`src_old/domain/`: the structs `User` and `Page`.

All the old functions are **free functions that take `db *sql.DB` as a parameter**. The caller (handler or service) therefore had to hold the DB itself, which is why `delivery` and `service` both imported `database/sql`.

## Proposed structure in `src`

```
src/
├── main.go                    open DB → Initialize → repos → services → router
├── models/                    plain structs, no dependencies
│   ├── user.go                User   (replaces the invalid Model.go)
│   └── page.go                Page
├── repositories/
│   ├── database.go            Open, Initialize (replaces Repository.go)
│   ├── schema.sql             embedded with //go:embed
│   ├── errors.go              ErrNotFound, ErrDuplicate
│   ├── user.go                UserRepository: Create, FindByUsername
│   └── page.go                PageRepository: Search
├── service/
│   ├── auth.go                AuthService   ← *UserRepository
│   ├── search.go              SearchService ← *PageRepository
│   └── weather.go             WeatherService (unchanged, no DB)
└── web/                       unchanged
```

### Repositories as structs with an injected DB

```go
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (int64, error)
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (models.User, error)
```

```go
type PageRepository struct {
	db *sql.DB
}

func NewPageRepository(db *sql.DB) *PageRepository

func (r *PageRepository) Search(ctx context.Context, query, language string) ([]models.Page, error)
```

`CreatePage` isn't used by any route in `src_old`, so it's left out for now (see question 1).

### Services receive repositories through interfaces

Each service declares the interface it needs **in the `service` package**, next to the service that uses it. This is the Go convention: the consumer defines the interface, and the concrete type satisfies it implicitly, with no `implements` keyword. As a result `service` doesn't import `repositories` at all. Only `main.go` knows both sides.

```go
// service/auth.go
type UserRepository interface {
	Create(ctx context.Context, user models.User) (int64, error)
	FindByUsername(ctx context.Context, username string) (models.User, error)
}

type AuthService struct {
	users UserRepository
}

func NewAuthService(users UserRepository) *AuthService
```

```go
// service/search.go
type PageRepository interface {
	Search(ctx context.Context, query, language string) ([]models.Page, error)
}

type SearchService struct {
	pages PageRepository
}

func NewSearchService(pages PageRepository) *SearchService
```

This follows the Go proverb "accept interfaces, return structs": constructors take interfaces and return concrete pointers.

The dependency graph becomes:

```
main ──► web ──► service ──► models
  │                 ▲
  └──► repositories ┘ (satisfies service's interfaces; imports models only)
```

### Wiring in `main.go`

```go
databasePath := os.Getenv("DATABASE_PATH")
if databasePath == "" {
	databasePath = "../data/whoknows.db"
}

db, err := repositories.Open(databasePath)
if err != nil {
	log.Fatalf("could not open database: %v", err)
}
defer db.Close()

if err := repositories.Initialize(context.Background(), db); err != nil {
	log.Fatalf("could not initialize database: %v", err)
}

userRepository := repositories.NewUserRepository(db)
pageRepository := repositories.NewPageRepository(db)

authService := service.NewAuthService(userRepository)
searchService := service.NewSearchService(pageRepository)
weatherService := service.NewWeatherService()

router := web.InitRouter(authService, searchService, weatherService)
```

### Dependencies

- `modernc.org/sqlite` must be added to `src/go.mod`, at the same version as `src_old` (`v1.58.0`). It's a pure-Go driver, so no cgo or C compiler is needed on Windows.

### Order of work

- **Now:** `models`, `repositories` (real SQL ported from `src_old`), the `modernc.org/sqlite` dependency. Nothing calls the repositories yet.
- **After confirmation:** add the interfaces to `service`, inject the repositories into `AuthService` and `SearchService`, and wire everything in `main.go`. The service methods keep returning route names. Business logic is ported in a later round.

## Decisions (answers to step 2)

1. Port the real SQL. Inject the repositories, but the services keep returning route names. `CreatePage` is left out.
2. Services depend on interfaces declared in `service`.
3. Repositories return `ErrNotFound` / `ErrDuplicate`. `service` never imports `database/sql`.
4. Weather is not touched.
5. `User` and `Page` are ported as they are, including the JSON tags.

## Things noticed along the way

- `router.Run(":8080")` ignores its error. With `defer db.Close()` added, that's worth checking: if the port is taken, `main` returns silently.
- Go file names are conventionally lower-case (`model.go`, not `Model.go`). The two placeholder files get replaced anyway.
- `IsDuplicateUserError` in `src_old` works by matching on the error *string*. That's fragile, but it's also the simplest option with `modernc.org/sqlite`, so it stays, just inside the repository.

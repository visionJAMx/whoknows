# Migration outline: register

## Task

Move the behaviour of `POST /api/register` from `src_old` into the new layers, following the *functionality* rather than the old implementation. Business logic (anything that neither depends on gin nor talks to the DB directly) goes into `AuthService.Register`. Hashing uses `src/crypto`.

## Principle: every layer translates errors at its boundary

```
repositories  speaks database  →  ErrUsernameTaken / ErrEmailTaken (both wrap ErrDuplicate)
service       speaks domain    →  *ValidationError{Fields: [{username, taken}]}, or a wrapped error
web           speaks HTTP      →  422 / 500 / 200 + JSON per docs/openapi/swagger.json
```

No layer passes the errors of the layer below on to the layer above. `service` imports `repositories` only for its error values. That's the allowed direction in strict layering. `web` never imports `repositories`.

## Implemented design

### repositories

- `errors.go`: `ErrUsernameTaken` and `ErrEmailTaken`, both wrapping `ErrDuplicate`.
- `user.go`: `duplicateUserError` checks for `SQLITE_CONSTRAINT_UNIQUE`, then reads the failing column from the driver message (`users.username` / `users.email`). The driver doesn't expose the column any other way.
- Uniqueness stays enforced by the UNIQUE constraint, because only the DB can do that atomically. Case-insensitivity comes from `COLLATE NOCASE`, which is ASCII-only.

### service

- `validation.go`: `Problem` (typed string constants), `FieldError{Field, Problem, Message}`, `ValidationError{Fields}`. `Field` is the form field name (`"password2"` included). The message describes the rule, so it's written next to the rule. There are no HTTP contract strings here.
- `auth.go`:
  - `UsernameMinLength`, `UsernameMaxLength`, `PasswordMinLength`: the rule values live only here.
  - `RegisterInput{Username, Email, Password, PasswordConfirm}`.
  - `Register(ctx, input) error`: normalise (trim username and email), then `validateRegistration` (a pure function that collects all problems), then `crypto.HashPassword`, then `users.Create`. A taken username or email becomes a `ValidationError` on that field.
  - Lengths are counted in characters (`utf8.RuneCountInString`), not bytes.
  - `PasswordConfirm` is optional: an empty value means not sent, and anything else must match.

### web

- `responses.go`: the contract's JSON shapes (`AuthResponse`, `ValidationErrorDetail`, `ValidationErrorResponse`) and the `errorTypes` map from `Problem` to the contract's `type` string.
- `auth.go`: `handleAPIRegister` owns its responses and reads top to bottom:
  1. Read the form and call `Register`.
  2. `errors.As(err, &validationErr)` → 422 with one detail per field.
  3. Any other error → 500 `"Could not create account"`, plus `c.Error(err)` so gin logs it.
  4. Success → 200 + `AuthResponse` (per the contract, not 201).
- There's no shared error helper on purpose ("a little copying is better than a little dependency"). Extract one only when several handlers show what is actually shared.

## Decisions

1. Errors stay in `repositories`, and `service` imports them (strict layering, the practical definition of isolation).
2. Contract error types belong to web (the `errorTypes` map). Messages are written in the service, next to the rules they describe.
3. JSON per the contract: 200 / 422 / 500.
4. `password2`: an empty value counts as not sent. There's no `*string`.
5. Lengths are counted in characters.
6. Duplicate username and duplicate email are reported separately.

## Known limitations / follow-ups

- `NOCASE` is ASCII-only: `Søren` and `SØREN` can both register. Fixing it means normalising (e.g. lower-casing) in the service before storing and before looking up.
- `data/whoknows.db` predates the schema changes and must be deleted to pick them up.
- `src/static/auth.js` is missing and no template loads a script, so in the browser the form shows raw JSON instead of inline errors / redirecting.
- `content LIKE '?'` in `repositories/page.go:28` is a literal string, so search only matches titles.
- The email check (contains `@`) is deliberately minimal.

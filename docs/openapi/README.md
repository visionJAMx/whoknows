# WhoKnows OpenAPI

`swagger.json` er den genererede **OpenAPI 3.1.0**-specifikation. Filnavnet betyder ikke, at formatet er Swagger 2.0.

## Generér fra projektets rod

```bash
go -C src run ./cmd/openapi
```

Det kræver Go-versionen fra `src/go.mod` og internet ved første download. Serveren skal ikke køre; der kræves hverken `.env`, sessionsnøgle eller database. Kommandoen ændrer ikke databasen.

Generatoren er fastlåst til `github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc5` i `src/cmd/openapi/main.go`. Den skriver resultatet til `docs/openapi/swagger.json`. Gentagen generering uden kildeændringer skal give samme fil.

**Ret ikke JSON-filen manuelt.** Ret annotationerne og dokumentationstyperne ved handlerne i `src/internal/delivery`, og generér igen. Commit både kildeændringer og den genererede fil.

## Hvordan dokumentationen bliver til

1. Swag læser `@Router`, `@Param`, `@Success`, `@Failure` m.fl. samt Go-typerne. Overordnede oplysninger læses fra `src/main.go`.
2. Go-kommandoen normaliserer det midlertidige output til underviserens kontrakt. Den fjerner tom `externalDocs` og utilsigtede `form`-felter samt tilføjer de union-typer, rc5 ikke udleder: nullable auth-felter/fejlbesked/sprog og string/integer i `loc`.
3. Først efter vellykket generering og normalisering skrives resultatet. Ved manglende forventede modeller afbryder kommandoen med fejl i stedet for at levere en delvist normaliseret specifikation.

Efterbehandlingen er versionsstyret kode, ikke manuelle rettelser eller en kopi af underviserens JSON. Ved opgradering af Swag skal dette trin genvurderes. Brug teamets Go-kommando ovenfor, ikke den tidligere direkte `swag init`-kommando.

Formularer dokumenteres med én request-type pr. endpoint. I denne prerelease giver separate `formData`-parametre sammen med `@Accept` forkert `oneOf`. Derfor benyttes én `@Param ... formData RequestType ...` uden `@Accept`.

Health-handleren ligger i delivery-laget, så dens operationstekst ikke overskriver projektets `info.description` i `main.go`.

## Omfang og kontrakt

| Metode | Sti | Svarformat |
| --- | --- | --- |
| GET | `/` | HTML |
| GET | `/login` | HTML |
| GET | `/register` | HTML |
| GET | `/about` | HTML |
| POST | `/api/login` | JSON |
| POST | `/api/register` | JSON |
| GET | `/api/logout` | JSON |
| GET | `/api/search` | JSON |
| GET | `/health` | JSON |

Udgangspunkt: [underviserens kontrakt](https://github.com/who-knows-inc/EK_DAT_DevOps_2026_Autumn/blob/main/00._Course_Material/02._Conventions_OpenAPI_DotEnv/01._Assignments/02._After/openapi.json) og [opgavebeskrivelsen](https://github.com/who-knows-inc/EK_DAT_DevOps_2026_Autumn/blob/main/00._Course_Material/02._Conventions_OpenAPI_DotEnv/01._Assignments/02._After/generate_openapi_specification.md).

- `/weather` og `/api/weather` er udtrykkeligt udskudt i opgaven. Statiske filer dokumenteres ikke enkeltvis.
- `/about`, `/health` og yderligere fejlsvar er tilladte udvidelser.
- API-søgning kræver `q`; `?q=` er tilladt og giver `{"data":[]}`. HTML-forsiden kræver ikke `q`.
- `language` er valgfrit. Vores kode bruger `en`, når parameteren udelades; denne standard er dokumenteret som en implementationsdetalje. En tom værdi søger på tværs af sprog. JSON Schema `null` i kontrakten betyder ikke, at teksten `language=null` sendes som JSON-null.
- Login kræver username/password. Registrering kræver username/email/password; password2 er valgfrit, men kontrolleres, hvis det medsendes.
- Kontrakten tillader null og numeriske lokationsindekser. De nuværende handlers sender konkrete auth-værdier og lokationer med tekst. Specifikationen bevarer kontraktens tilladte typer uden at ændre runtime-svarene.
- Registreringens eksisterende ekstra validering (brugernavn 3–50 bytes, password mindst 10 bytes, email indeholder @ og password2-match) er beskrevet i schemaet. Opgaven her ændrer ikke disse forretningsregler. `len` i Go måler bytes, mens JSON Schema `minLength` måler tegn; derfor er disse længder beskrevet i tekst.
- Modelnavne, danske beskrivelser og tags kan afvige fra underviserens. `type: [string, null]` og `anyOf` med string/null udtrykker samme tilladte typer i OpenAPI 3.1.

## Åbn i Swagger Editor

Åbn [Swagger Editor](https://editor.swagger.io/), vælg **File → Import file**, og vælg `docs/openapi/swagger.json`. Kontrollér alle ni stier og fejl-/advarselspanelet. Import alene beviser ikke, at serveren følger kontrakten.

Der er ikke hardcodet en VM-adresse eller et produktionsmiljø i dokumentet. Swagger Editors **Try it out** kræver særskilt serveradresse og browseradgang (HTTPS/CORS/cookies). Import/visning virker uden server. Brug ikke rigtige passwords eller sessionscookies i screenshots eller PR-beskrivelser.

## Kontrol inden PR

```bash
go -C src run ./cmd/openapi
go -C src build ./...
go -C src vet ./...
git diff --check
```

Som ekstra OpenAPI-validering, hvis `uv` er installeret:

```bash
uv tool run --from openapi-spec-validator==0.7.2 openapi-spec-validator docs/openapi/swagger.json
```

Python-værktøjet er kun til dokumentvalidering og er ikke en applikationsafhængighed. Go-kommandoen ovenfor er tilstrækkelig til generering. Sammenlign desuden relevante statuskoder, Content-Type og JSON-svar fra serveren med dokumentationen; schema-validering er ikke en runtime-test.

## Verifikation 24. september 2026

- Build, `go vet`, OpenAPI 3.1-validator og strukturel sammenligning af de syv ikke-weather-operationer i underviserens kontrakt gennemført. Beskrivelser, navne og dokumenterede ekstra defaults er ikke krav om tekstlig identitet.
- En midlertidig handlerkontrol med isoleret database gennemførte HTML-sider, health, søgning med/uden q og med/uden resultater, login/registreringsvalidering, forkert login, succesfuld registrering/login og logout. Login satte sessionens bruger-ID; logout fjernede det i den efterfølgende request. Søgefejl på lukket testdatabase gav 500. Kontrolfiler og testdatabase er ikke del af leverancen.
- Dette er ikke en fuld sikkerhedstest eller bevis for tilbagekaldelse af tidligere kopierede signerede cookies. Alle mulige 500-fejlstier er heller ikke fremprovokeret.
- `go -C src test ./...` blev også kørt: den eksisterende `TestLoginPage` fejler, fordi testen stadig indlæser den enkelte template i stedet for applikationens fælles multitemplate-layout. Den målrettede kontrol med korrekt renderer bestod. Denne gamle testopsætning er ikke ændret som del af OpenAPI-opgaven; den fulde testsuite er derfor ikke grøn.

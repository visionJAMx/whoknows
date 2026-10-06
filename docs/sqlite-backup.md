# SQLite: backup og sikker gendannelsesprøve

Vi beholder SQLite. Underviserens PostgreSQL-eksempel skal ikke køres mod
projektets data. Scriptet ændrer hverken skema, brugere eller søgesider.

## Krav og databasevalg

Kør fra projektets rod med Bash og `sqlite3` installeret på værten.
Angiv den database, appen faktisk bruger. Scriptet læser ikke `.env`.
Den lokale standard er `data/whoknows.db`; VM-stien kan være anderledes.

## Tag backup

```bash
bash scripts/sqlite-backup.sh data/whoknows.db backups
```

Den udskrevne sti er den nye backup. Hver kørsel opretter en unik privat mappe
med `snapshot.db` og en `SUCCESS`-fil, når integritetskontrollen er bestået.
Backuppen indeholder hele databasen, inklusive brugere og password-hashes.
Del eller commit den ikke. En mappe uden `SUCCESS` er ikke en godkendt backup.

SQLite `.backup` bruges frem for en filkopi, så også en aktiv database med
WAL kan sikkerhedskopieres konsistent. Ved låsning/andre fejl stopper scriptet;
kontrollér fejlen og prøv igen. Se https://www.sqlite.org/backup.html.

## Gendan til en separat kopi — ikke oven i den aktive database

Brug stien udskrevet ved backup som første argument:

```bash
bash scripts/sqlite-backup.sh backups/DIT-BACKUPNAVN/snapshot.db restore-check
```

Dette genskaber databasen i en ny mappe med SQLite-backupfunktionen.
Kontrollér den udskrevne sti:

```bash
sqlite3 -readonly restore-check/DIT-NYE-NAVN/snapshot.db 'PRAGMA integrity_check; SELECT COUNT(*) FROM pages;'
```

Sammenlign også forventede tabeller og data med backuppen. Integritetskontrol
alene beviser ikke, at de ønskede data findes. En app-prøve på den gendannede
kopi bør indgå før en faktisk driftsgendannelse.

## Faktisk driftsgendannelse (manuel procedure)

1. Stop appen og alle andre skrivere til databasen.
2. Bevar den nuværende database og eventuelle `-wal`, `-shm` og `-journal`-filer
   samlet til fejlsøgning; slet dem ikke og bland dem ikke med en anden database.
3. Gendan til en separat mappe med scriptet, og verificér indholdet.
4. Peg appens `DATABASE_PATH` på den kontrollerede kopi med korrekte rettigheder.
5. Start appen og kontrollér `/health`, søgning og login.

## Docker og begrænsninger

SQLite skal ligge i en persistent mappe/volume uden for applikationsimaget.
Mount hele datamappen, ikke kun `.db`-filen, så journal/WAL-filer følger med.
Et volume er ikke en backup. Den konkrete Docker-integration er endnu ikke
implementeret her, da denne projektversion ikke indeholder Compose/Dockerfile.

Denne første batch har ingen automatisk tidsplan, retention eller ekstern
backupdestination. Der slettes ingen gamle backups automatisk. En lokal backup
beskytter ikke mod tab af hele maskinen: næste batch skal vælge sikker ekstern
opbevaring, adgangskontrol, tidsplan og overvågning af mislykkede backups.

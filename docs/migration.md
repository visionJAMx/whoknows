# Cutover — one-time instructions

> Temporary doc: delete after the cutover is done and verified. Rehearsed 2026-11-22 (staging, all green).

> **Set `GIN_MODE=release`** on the whoknows resource env before the merge. Session cookies are only marked `Secure` in release mode (`gin.Mode() == gin.ReleaseMode`, `src/main.go`); unset, the app runs in debug mode behind the proxy — it works, but logs verbosely and cookies lack the Secure flag.

**Values** — volume `postgres-data-lfrp5wut8aksukr6dlki5mjp` (never delete) · expected counts **4070 / 49** · migrations `DB_HOST` `lfrp5wut8aksukr6dlki5mjp` (dead after merge — re-point) · SQLite backups on VM: `~/whoknows-sqlite-backup-*.db`, `~/seed_*.sql`.

**Preconditions — all verified**: schema head `d489fff393dd` · `external: true` volume mechanism proven in cloud · staging retired, volume intact · `migration` pushed.

**Day of:**

1. **Drift call** — optional re-seed (recipe below) or accept the gap since 2026-11-21.
2. whoknows resource env: `POSTGRES_VOLUME_NAME` ✓, `SESSION_SECRET` ✓, creds = PG17 resource's.
3. PG17 resource stays stopped — two postgres on one volume = corruption.
4. **Merge `migration` → `prod`** → deploy.yml → webhook → stack redeploys; postgres attaches the external volume, app boots on Postgres. **No initdb** in the postgres log (initdb = wrong volume → fix env, redeploy).
5. Verify: `/health`, counts 4070 / 49, one real login, one search.
6. Re-point migrations `DB_HOST` → stack postgres hostname (Coolify) — same session.
7. Delete the PG17 resource **without** "delete associated volumes".
8. After: branch protection on `DB-prod` (optional re-align first); delete this doc.

**Rollback**: revert the merge on `prod` → the SQLite-era stack redeploys and runs as before; the Postgres volume is untouched throughout.

**Re-seed (~10 min, VM — new stack not up yet, so a throwaway postgres on the volume):**

```bash
docker cp <sqlite-app>:/data/whoknows.db ~/bu.db
docker run --rm -v ~/bu.db:/d.db keinos/sqlite3 sqlite3 /d.db ".dump users" | grep '^INSERT' > ~/seed_users.sql
docker run --rm -v ~/bu.db:/d.db keinos/sqlite3 sqlite3 /d.db ".dump pages"  | sed -n '/^INSERT/,$p' > ~/seed_pages.sql   # pages rows span lines
sudo docker run -d --name reseed -v postgres-data-lfrp5wut8aksukr6dlki5mjp:/var/lib/postgresql/data -p 127.0.0.1:5433:5432 postgres:17   # must NOT initdb
sudo docker exec reseed psql -U postgres -d whoknows -c "DELETE FROM users; DELETE FROM pages;"   # old seed has explicit ids — replace, don't append
sudo docker exec -i reseed psql -U postgres -d whoknows < ~/seed_users.sql
sudo docker exec -i reseed psql -U postgres -d whoknows < ~/seed_pages.sql
sudo docker exec reseed psql -U postgres -d whoknows -c "SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));"
# verify counts + MAX(id) + SUM(LENGTH(content)) against the sqlite dump, then:
sudo docker rm -f reseed   # free the volume for the merge
```

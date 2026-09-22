# Deployment (temporary VM setup)

> **Temporary:** This setup will be replaced by a Docker-based deployment. Before switching, follow [Removing this setup](#removing-this-setup), otherwise the old service will keep port 8080 busy.

The app runs on our VM as a systemd service. Every push to `main` builds and deploys it automatically through [.github/workflows/deploy.yml](../.github/workflows/deploy.yml).

## How it works

```
push to main
  > GitHub Actions builds a Linux binary
  > uploads binary + templates/ + static/ to /opt/whoknows/incoming over SSH as user "deploy"
  > after files have been transferred, they are moved into place and the following command is run: sudo systemctl restart whoknows

```

Tests are currently **not** run before deploying.

The app runs from ``/opt/whoknows`` on the VM. 

## Useful commands (on the VM)

```bash
systemctl status whoknows              # Check if service is running
sudo journalctl -u whoknows -f         # follow the app log
sudo journalctl -u whoknows -n 100     # last 100 log lines
sudo systemctl restart whoknows        # restart manually
sudo -u deploy nano /opt/whoknows/.env # edit .env file, (remember to restart service afterwards)
```

## Removing this setup

When we eventually move to a better deployment method, we fully undo the setup for this deployment with the following steps:

1. Delete (or disable) `.github/workflows/deploy.yml` so pushes stop deploying.
2. **Back up the database** if it should be kept: `/opt/whoknows/data/whoknows.db`.
3. On the VM:
   ```bash
   sudo systemctl disable --now whoknows          # stop it and free port 8080
   sudo rm /etc/systemd/system/whoknows.service
   sudo systemctl daemon-reload
   sudo rm /etc/sudoers.d/whoknows
   sudo userdel -r deploy                         # also removes the GitHub key
   sudo rm -rf /opt/whoknows                      # only after backing up the database
   ```
4. Delete the `SSH_HOST`, `SSH_USER` and `SSH_PRIVATE_KEY` secrets in GitHub if the new setup doesn't reuse them.
5. Delete this file.

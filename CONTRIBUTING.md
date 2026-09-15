# Sådan bidrager du til WhoKnows

Vi arbejder struktureret med branches for at undgå konflikter:
1. **Aldrig commit direkte til `main` eller `dev`** lokalt.
2. Opret altid en ny feature-branch ud fra den seneste kode:
   ```bash
   git checkout main (eller dev)
   git pull origin main
   git checkout -b feature/dit-feature-navn
   
Klon repositoriet

Generer jeres ssh-nøgle og knyt den til jeres Github-konto.

Opsæt miljøvariabler ved at bruge jeres SSH-forbindelse til at hoppe på den delte VM: ssh whoknows-vm.

Gå til VM og tilføj jeres miljøvariabler i en env. fil.

Kør applikationen

# Pull Request Processen

Push din branch til GitHub:

Sørg for at din kode er fri for mergekonflikter og at router-opsætningen fungerer korrekt med de nye ændringer.

Få en anden fra holdet til at gennemse koden, før den merges.
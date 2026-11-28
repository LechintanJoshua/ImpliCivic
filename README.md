# ImpliCivic

## 1. Descrierea proiectului si problema abordata
**ImpliCivic** este o platforma web prin care aducem datele publice locale intr-un singur loc si oferim un mod simplu prin care locuitorii din Timisoara pot raporta probleme intalnite in oras.

In momentul de fata, datele puse la dispozitie la nivel local, cum ar fi locatiile parcarilor, punctele de colectare a deseurilor sau statiile de transport, se gasesc in fisiere CSV greu de parcurs. Totodata, sesizarile oamenilor legate de gropi, iluminat stradal defect sau deseuri ajung de obicei pe retelele sociale, unde se pierd rapid si nu au un flux clar de rezolvare.

## 2. Obiectivele proiectului
- Parsarea si afisarea datelor din fisiere CSV intr-o interfata web usor de citit si filtrat.
- Oferirea unui sistem prin care cetatenii autentificati pot trimite sesizari si pot urmari stadiul acestora.
- Crearea unui panou de administrare pentru actualizarea statusurilor si exportul sesizarilor inapoi in format CSV.

## 3. Utilizatori vizati
- **Vizitatori anonimi:** Oricine vrea doar sa consulte datele publice sau sa vada incidentele deja raportate, fara sa aiba cont.
- **Cetateni autentificati:** Locuitorii care vor sa trimita sesizari concrete, sa isi urmareasca istoricul propriu si sa poata sterge sesizarile inca nepreluate.
- **Administratori:** Cei care modereaza sesizarile, le actualizeaza stadiul cu observatii si descarca rapoarte centralizate in format CSV.

## 4. Echipa si responsabilitati
- **Lechintan Joshua-Leonardo-Augustin** – *Arhitectura Backend & Persistenta Datelor*: Configurarea serverului in Go, structura modulara, schema si persistenta in PostgreSQL, procesarea concurenta a fisierelor CSV, exportul CSV, securitatea server-side (validari, prevenire SQLi/XSS, autentificare si roluri) si administrarea repository-ului de GitHub.
- **Losnariu Marius-Madalin** – *Frontend & UI/UX Design*: Structura paginilor in HTML5, stilizarea responsive cu CSS pentru mobil si desktop, designul interfetei, interactivitatea acesteia si respectarea cerintelor de accesibilitate.
- **Teodorovits Andras-Karoly** – *Integrare Full-Stack & Logica Client*: Conectarea frontend-ului la API prin cereri asincrone, validarea formularelor pe client si testarea fluxurilor end-to-end.

## 5. Functionalitati implementate
1. **Autentificare si roluri (`citizen` / `admin`):** inregistrare cu hash-uire `bcrypt`, autentificare cu token JWT semnat `HMAC-SHA256`, verificare token si restrictii de acces pe roluri (`401` / `403`).
2. **Import concurent si consultare date publice CSV:** citirea in streaming a fisierelor din `data/` la pornirea serverului (cand tabela este goala), cu worker pool de 4 goroutines si parsare defensiva; consultare cu filtrare pe categorie, cautare text si paginare.
3. **Gestiunea sesizarilor urbane:** creare sesizare cu sanitizare anti-XSS si validare lungime titlu, listare publica cu filtre, consultare detalii dupa ID, listare sesizari proprii si stergere permisa cetateanului doar daca este proprietar si statusul este `INREGISTRAT`.
4. **Moderare si export CSV pentru administratori:** actualizarea stadiului unei sesizari (`INREGISTRAT`, `IN_LUCRU`, `REZOLVAT`, `RESPINS`) cu note administrative si exportul sesizarilor filtrate direct in format `.csv`.

## 6. Tehnologii utilizate
- **Backend:** Go 1.26 (`net/http`, `database/sql`, `encoding/csv`, `encoding/json`, `crypto/hmac`, `golang.org/x/crypto/bcrypt`, `github.com/joho/godotenv`).
- **Baza de date:** PostgreSQL + driverul `github.com/jackc/pgx/v5`.
- **Frontend:** HTML5, CSS3, TypeScript, React.
- **Unelte:** Git, GitHub, `psql`, `curl`.

## 7. Structura repository-ului
```text
ImpliCivic/
├── cmd/
│   └── api/                  # Punctul de intrare, server HTTP, middleware si handlere
├── internal/
│   ├── auth/                 # Hashing parole (bcrypt) si token-uri JWT HS256
│   ├── csvproc/              # Import streaming CSV cu goroutines si export CSV
│   ├── db/                   # Schema SQL si pool-ul de conexiuni PostgreSQL
│   └── store/                # Interogari SQL parametrizate pe tabele
├── data/                     # Fisiere CSV cu date publice din Timisoara
├── docs/                     # Documentatia tehnica si scenariile de testare
├── architecture/             # Diagrame PDF de arhitectura si baza de date
├── mockup/                   # Mockup-uri PDF pentru interfata
├── .env.example              # Model de configurare pentru variabilele de mediu
└── go.mod                    # Modulul Go si dependintele proiectului
```

## 8. Instalare, configurare si rulare locala

### Cerinte preliminare
- **Go** 1.22 sau mai nou (proiectul foloseste Go `1.26.8`)
- **PostgreSQL** instalat si pornit local

### Pasi de configurare
1. **Creeaza utilizatorul si baza de date in PostgreSQL:**
   ```bash
   sudo -u postgres psql -c "CREATE USER implicivic WITH PASSWORD 'alege_o_parola_locala';"
   sudo -u postgres psql -c "CREATE DATABASE implicivic OWNER implicivic;"
   ```
2. **Aplica schema SQL pentru tabele (`users`, `incidents`, `public_datasets`):**
   ```bash
   psql -h localhost -U implicivic -d implicivic -f internal/db/schema.sql
   ```
3. **Configureaza variabilele de mediu locale:**
   ```bash
   cp .env.example .env
   ```
   Editeaza `.env` cu parola aleasa pentru PostgreSQL si un secret aleatoriu pentru `JWT_SECRET` (ex. generat cu `openssl rand -hex 32`):
   ```text
   DATABASE_URL=postgres://implicivic:alege_o_parola_locala@localhost:5432/implicivic?sslmode=disable
   PORT=8080
   JWT_SECRET=secret_generat_local
   ALLOWED_ORIGIN=http://localhost:5173
   ```
4. **Instaleaza dependintele si porneste serverul backend:**
   ```bash
   go mod download
   go run ./cmd/api
   ```
   La prima pornire, serverul verifica conexiunea cu PostgreSQL, importa automat randurile valide din `data/*.csv` in `public_datasets` si asculta pe `http://localhost:8080`.
5. **Verificare rapida:**
   ```bash
   curl -i http://localhost:8080/api/health
   ```

## 9. Documentatia proiectului
Toate detaliile tehnice si scenariile de testare se gasesc in proiect:
- [Cerinte Functionale si Nefunctionale](docs/requirements.md)
- [Scenarii de Utilizare (Use Cases)](docs/use-cases.md)
- [Arhitectura Aplicatiei](docs/architecture.md)
- [Modelul de Date si Specificatia API](docs/data-model.md)
- [Testare si Rezultate (Testing)](docs/testing.md)
- [Planul de Implementare](docs/implementation-plan.md)
- [Diagrama Arhitecturala (PDF)](architecture/Arhitectura.pdf)
- [Diagrama Bazei de Date (PDF)](architecture/Database.pdf)
- [Mockup Interfata (PDF)](mockup/ImpliCivic_mockup.pdf)

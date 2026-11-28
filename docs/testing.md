# Scenarii de Testare si Rezultate – ImpliCivic

Acest document centralizeaza scenariile de testare rulate asupra API-ului de backend si a bazei de date PostgreSQL, datele de proba folosite, rezultatele obtinute si limitarile identificate.

---

## 1. Mediul de Testare si Datele de Proba

- **Server API:** Go `1.26.8` rulat local pe `http://localhost:8080`
- **Baza de date:** PostgreSQL local (`implicivic` pe `localhost:5432`)
- **Unealta de testare HTTP/DB:** `curl` si clientul `psql`
- **Date de proba CSV (`data/`):**
  - `parcari_publice_tm.csv` (18 linii: 1 antet, 15 randuri valide, 2 randuri corupte intentionat — nume lipsa si coordonate `abc,xyz`)
  - `puncte_reciclare_tm.csv` (18 linii: 1 antet, 15 randuri valide, 2 randuri corupte intentionat — coloane lipsa si coordonata `not_a_lat`)
  - `statii_transport_tm.csv` (18 linii: 1 antet, 15 randuri valide, 2 randuri corupte intentionat — nume lipsa si coloane lipsa)
- **Conturi de proba:**
  - `joshua@test.ro` (rol `citizen`)
  - `admin@test.ro` (inregistrat initial ca `citizen`, promovat apoi in PostgreSQL la rol `admin` pentru testarea rutelor administrative)

---

## 2. Scenarii Testate si Rezultate Obtinute

### 2.1. Pornire server, conexiune PostgreSQL si import CSV concurent
| Nr. | Scenariu | Comanda / Actiune | Rezultat asteptat | Rezultat obtinut |
| :--- | :--- | :--- | :--- | :--- |
| T1 | Health check server + DB | `GET /api/health` | `200 OK`, JSON `{"db":"connected","status":"ok"}` | **Trecut** (`200 OK`, headere CORS prezente) |
| T2 | Import CSV la pornire cu randuri corupte | Pornire `go run ./cmd/api` cu tabela `public_datasets` goala, apoi interogare `SELECT category, COUNT(*) FROM public_datasets GROUP BY category;` | Import doar pentru cele 45 de randuri valide; cele 6 randuri corupte sarite fara crash | **Trecut** (`parcare: 15`, `reciclare: 15`, `transport: 15`, total `45`) |

### 2.2. Autentificare si Protectia Datelor
| Nr. | Scenariu | Actiune / Input | Rezultat asteptat | Rezultat obtinut |
| :--- | :--- | :--- | :--- | :--- |
| T3 | Inregistrare cont nou valid | `POST /api/auth/register` cu `username`, `email`, `password` | `201 Created`, `role: "citizen"`, fara `password_hash` | **Trecut** (`201 Created`, `id: 1`, fara campul `password_hash`) |
| T4 | Inregistrare cu email/username duplicat | `POST /api/auth/register` cu aceleasi date | `409 Conflict` | **Trecut** (`409 Conflict`, `{"error":"username or email already exists"}`) |
| T5 | Login cu parola gresita | `POST /api/auth/login` cu parola gresita | `401 Unauthorized` | **Trecut** (`401 Unauthorized`, `{"error":"invalid credentials"}`) |
| T6 | Login valid | `POST /api/auth/login` cu date corecte | `200 OK`, token JWT + obiect `user` | **Trecut** (`200 OK`, token HS256 emis) |

### 2.3. Consultare si Filtrare Date Publice (`GET /api/datasets`)
| Nr. | Scenariu | Cerere | Rezultat asteptat | Rezultat obtinut |
| :--- | :--- | :--- | :--- | :--- |
| T7 | Paginare date publice | `GET /api/datasets?limit=5&offset=0` | `200 OK`, exact 5 inregistrari | **Trecut** (`200 OK`, 5 inregistrari, inclusiv campuri `address: null` pastrate corect) |
| T8 | Filtrare categorie + cautare text | `GET /api/datasets?category=reciclare&search=Martirilor` | `200 OK`, punctele potrivite | **Trecut** (`200 OK`, 1 rezultat: `Insulă Ecologică Calea Martirilor`) |
| T9 | Cautare fara potriviri | `GET /api/datasets?category=reciclare&search=Cetatii` | `200 OK` cu array gol `[]` (nu `null`) | **Trecut** (`200 OK`, `[]`) |
| T10 | Parametru invalid de paginare | `GET /api/datasets?limit=-1` | `400 Bad Request` | **Trecut** (`400 Bad Request`, `{"error":"invalid limit"}`) |

### 2.4. CRUD Sesizari, Anti-XSS si Control Acces (RBAC)
| Nr. | Scenariu | Cerere / Input | Rezultat asteptat | Rezultat obtinut |
| :--- | :--- | :--- | :--- | :--- |
| T11 | Listare sesizari pe tabela goala | `GET /api/incidents` | `200 OK` cu `[]` | **Trecut** (`200 OK`, `[]`) |
| T12 | Creare sesizare fara token | `POST /api/incidents` fara `Authorization` | `401 Unauthorized` | **Trecut** (`401 Unauthorized`, `{"error":"missing or invalid token"}`) |
| T13 | Creare sesizare cu payload XSS | `POST /api/incidents` cu token `citizen` si `<script>alert(1)</script>` in descriere | `201 Created`, status `INREGISTRAT`, tag-uri HTML escapate | **Trecut** (`201 Created`, `<script>` transformat in `&lt;script&gt;` / `\u0026lt;script\u0026gt;`) |
| T14 | Consultare, lista proprie si stergere | `GET /api/incidents/1`, `GET /api/users/me/incidents`, `DELETE /api/incidents/1`, apoi `GET /api/incidents/1` | `200 OK`, `200 OK`, `204 No Content`, apoi `404 Not Found` | **Trecut** |
| T15 | Acces `citizen` pe rute de admin | `PATCH /api/admin/incidents/2` si `GET /api/admin/incidents/export` cu token `citizen` | `403 Forbidden` pe ambele | **Trecut** (`403 Forbidden`, `{"error":"admin access required"}`) |
| T16 | Moderare sesizare de catre admin | `PATCH /api/admin/incidents/2` cu token `admin` (`status: "IN_LUCRU"`, `admin_notes`) | `200 OK`, status si note actualizate + `updated_at` | **Trecut** (`200 OK`, status `IN_LUCRU`) |
| T17 | Listare sesizari dupa actualizare (fara filtru) | `GET /api/incidents` | `200 OK` cu lista sesizarilor | **Trecut** (`200 OK`, returneaza sesizarea #2) |
| T18 | Export CSV cu date si export gol | `GET /api/admin/incidents/export` si `?status=REZOLVAT` cu token `admin` | `200 OK` (`text/csv`), cu antet + randuri; la 0 potriviri doar antetul | **Trecut** (`251` bytes cu date; `106` bytes doar linia de antet la `?status=REZOLVAT`) |

---

## 3. Erori Identificate pe Parcurs si Remediate
1. **Eroare SQL in `ListDatasets`:** wildcard-urile pentru `ILIKE` au fost scrise initial cu ghilimele duble (`"%"`) in locul ghilimelelor simple (`'%'`), provocand eroare PostgreSQL la prima interogare; corectat in `internal/store/datasets.go`.
2. **Filtru optional in `listIncidentsHandler`:** dupa introducerea helper-ului `isValidStatus`, apelul `GET /api/incidents` fara parametrul `status` returna `400 Bad Request`; remediat prin conditia `status != "" && !isValidStatus(status)` pe branch-ul `fix/list-incidents-status`.

---

## 4. Limitari Cunoscute
- **Sensibilitate la diacritice in cautare:** operatorul `ILIKE` din PostgreSQL nu echivaleaza automat caracterele fara diacritice cu cele cu diacritice (ex. `search=Cetatii` nu potriveste `Strada Cetății 19`, motiv pentru care testul T9 a intors `[]`).
- **Incarcare CSV din UI:** importul CSV ruleaza automat la pornirea serverului cand tabela `public_datasets` este goala; un endpoint de upload CSV din panoul de administrare ramane planificat pentru versiunea finala.

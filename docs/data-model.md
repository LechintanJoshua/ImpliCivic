# Modelul de Date si Specificatia API – ImpliCivic

## 1. Schema Bazei de Date (PostgreSQL)

Schema este definita idempotent (`CREATE TABLE IF NOT EXISTS`) in [internal/db/schema.sql](../internal/db/schema.sql) si foloseste tipuri cu fus orar (`TIMESTAMPTZ`) si precizie zecimala fixa pentru coordonate (`NUMERIC(9, 6)`).

### Tabela `users`
Conturile utilizatorilor (cetateni si administratori).
* `id` (`BIGSERIAL PRIMARY KEY`)
* `username` (`VARCHAR(50) UNIQUE NOT NULL`)
* `email` (`VARCHAR(255) UNIQUE NOT NULL`)
* `password_hash` (`TEXT NOT NULL`) – stocat exclusiv prin hash `bcrypt` (exclus din JSON prin tag-ul `json:"-"`).
* `role` (`VARCHAR(50) NOT NULL DEFAULT 'citizen' CHECK (role IN ('citizen', 'admin'))`)
* `created_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)

### Tabela `incidents`
Sesizarile urbane trimise de cetateni.
* `id` (`BIGSERIAL PRIMARY KEY`)
* `user_id` (`BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE`)
* `title` (`VARCHAR(200) NOT NULL`)
* `category` (`VARCHAR(50) NOT NULL`) – ex: `infrastructura`, `deseuri`, `iluminat`, `parcare`.
* `description` (`TEXT NOT NULL`)
* `address` (`VARCHAR(255) NOT NULL`)
* `latitude` (`NUMERIC(9, 6) NULL`)
* `longitude` (`NUMERIC(9, 6) NULL`)
* `status` (`VARCHAR(30) NOT NULL DEFAULT 'INREGISTRAT' CHECK (status IN ('INREGISTRAT', 'IN_LUCRU', 'REZOLVAT', 'RESPINS'))`)
* `admin_notes` (`TEXT NULL`) – observatii adaugate de administrator la schimbarea starii.
* `created_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)
* `updated_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)

### Tabela `public_datasets`
Punctele de interes public importate din fisierele CSV din `data/`.
* `id` (`BIGSERIAL PRIMARY KEY`)
* `dataset_name` (`VARCHAR(100) NOT NULL`) – numele fisierului sursa (ex: `parcari_publice_tm.csv`).
* `category` (`VARCHAR(50) NOT NULL`)
* `name` (`VARCHAR(200) NOT NULL`)
* `address` (`VARCHAR(255) NULL`) – optional, unele puncte publice nu au adresa stradala.
* `latitude` (`NUMERIC(9, 6) NULL`)
* `longitude` (`NUMERIC(9, 6) NULL`)
* `imported_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)

---

## 2. Relatii, Constrangeri si Mecanismul de Persistenta
* **Relatia `users` 1 : N `incidents`:** realizata prin cheia externa `incidents.user_id` catre `users(id)` cu `ON DELETE CASCADE`, astfel incat stergerea unui cont sterge automat sesizarile asociate.
* **`public_datasets` decuplata:** nu are cheie externa catre `users`, fiind destinata consultarii publice anonime.
* **Connection Pool (`database/sql` + `pgx/v5/stdlib`):** serverul pastreaza un pool comun de conexiuni (`MaxOpenConns = 25`, `MaxIdleConns = 5`, `ConnMaxLifetime = 5 minute`) reutilizat intre cererile concurente si verificat la pornire prin `PingContext` cu timeout de 5 secunde.
* **Protectie SQL Injection si erori de unicitate:** 100% din interogarile din `internal/store/` folosesc placeholders parametrizate (`$1, $2, ...`). Erorile de duplicat din PostgreSQL (codul `23505`) sunt mapate la `409 Conflict`.
* **Import initial idempotent:** la pornire, serverul verifica numarul de randuri din `public_datasets`; importul CSV din `data/` ruleaza doar daca tabela este goala (`COUNT(*) == 0`), evitand dublarea datelor la restart.

---

## 3. Specificatia Endpoint-urilor API REST

Toate rutele au prefixul `/api` si intorc JSON (cu exceptia exportului CSV si a raspunsului `204 No Content`).

| Metoda | Endpoint | Acces | Descriere & Query Params | Raspuns Succes | Erori asteptate |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/health` | Public | Verifica disponibilitatea serverului si a bazei de date. | `200 OK` | `503` (DB inaccesibil) |
| `POST` | `/api/auth/register` | Public | Inregistreaza un utilizator nou (`role = citizen`), parola hash-uita cu `bcrypt`. | `201 Created` | `400` (date invalide), `409` (email/username existent) |
| `POST` | `/api/auth/login` | Public | Autentifica utilizatorul si returneaza token-ul JWT + datele contului. | `200 OK` | `400` (campuri lipsa), `401` (date gresite) |
| `GET` | `/api/datasets` | Public | Listeaza datele publice din CSV. Suporta `?category=...&search=...&limit=20&offset=0`. | `200 OK` | `400` (limit/offset invalid) |
| `GET` | `/api/incidents` | Public | Listeaza sesizarile publice. Suporta `?status=...&category=...`. | `200 OK` | `400` (status invalid) |
| `GET` | `/api/incidents/{id}` | Public | Returneaza detaliile unei sesizari dupa ID. | `200 OK` | `400` (ID invalid), `404` (inexistent) |
| `POST` | `/api/incidents` | Citizen (`Bearer`) | Creeaza o sesizare noua (status fortat `INREGISTRAT`, `user_id` din token, sanitizare anti-XSS). | `201 Created` | `400` (campuri goale / titlu > 200 caractere), `401` |
| `GET` | `/api/users/me/incidents` | Citizen (`Bearer`) | Returneaza doar sesizarile create de utilizatorul autentificat. | `200 OK` | `401` |
| `DELETE` | `/api/incidents/{id}` | Citizen / Admin (`Bearer`) | Sterge o sesizare. Cetateanul o poate sterge doar daca este proprietar si statusul este `INREGISTRAT`. | `204 No Content` | `400`, `401`, `403`, `404` |
| `PATCH` | `/api/admin/incidents/{id}` | Admin (`Bearer`) | Actualizeaza `status`, `admin_notes` (sanitizat anti-XSS) si `updated_at`. | `200 OK` | `400` (status invalid), `401`, `403`, `404` |
| `GET` | `/api/admin/incidents/export` | Admin (`Bearer`) | Descarca direct sesizarile filtrate (`?status=...&category=...`) in format `.csv`. | `200 OK` (`text/csv`) | `400`, `401`, `403` |

### Reguli de format si cazuri-limita
- **Liste fara rezultate:** `GET /api/datasets`, `GET /api/incidents` si `GET /api/users/me/incidents` returneaza intotdeauna un array JSON gol `[]` (niciodata `null`).
- **Export CSV fara potriviri:** `GET /api/admin/incidents/export` returneaza `200 OK` (`text/csv`) cu un fisier CSV valid care contine doar linia de antet (`12` coloane).
- **Validare lungime titlu:** titlul gol sau mai lung de `200` de caractere Unicode (`utf8.RuneCountInString`) este respins cu `400 Bad Request`.

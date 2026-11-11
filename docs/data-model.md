# Modelul de Date si Specificatia API – ImpliCivic

## 1. Schema Bazei de Date (PostgreSQL)

### Tabela `users`
Conturile utilizatorilor (cetateni si administratori).
* `id` (`BIGSERIAL PRIMARY KEY`)
* `username` (`VARCHAR(50) UNIQUE NOT NULL`)
* `email` (`VARCHAR(255) UNIQUE NOT NULL`)
* `password_hash` (`TEXT NOT NULL`) – stocat exclusiv prin hash bcrypt.
* `role` (`VARCHAR(20) NOT NULL DEFAULT 'citizen'`) – valori permise: `citizen`, `admin`.
* `created_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)

### Tabela `incidents`
Sesizarile urbane trimise de cetateni.
* `id` (`BIGSERIAL PRIMARY KEY`)
* `user_id` (`BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE`)
* `title` (`VARCHAR(200) NOT NULL`)
* `category` (`VARCHAR(50) NOT NULL`) – ex: `infrastructura`, `deseuri`, `iluminat`, `parcare`.
* `description` (`TEXT NOT NULL`)
* `address` (`VARCHAR(255) NOT NULL`)
* `latitude` (`NUMERIC(9,6) NULL`)
* `longitude` (`NUMERIC(9,6) NULL`)
* `status` (`VARCHAR(30) NOT NULL DEFAULT 'INREGISTRAT'`) – valori permise: `INREGISTRAT`, `IN_LUCRU`, `REZOLVAT`, `RESPINS`.
* `admin_notes` (`TEXT NULL`) – observatii adaugate de administrator la schimbarea starii.
* `created_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)
* `updated_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)

### Tabela `public_datasets`
Punctele de interes public parsate din fisierele CSV.
* `id` (`BIGSERIAL PRIMARY KEY`)
* `dataset_name` (`VARCHAR(100) NOT NULL`) – sursa CSV-ului (ex: `parcari_publice_tm.csv`).
* `category` (`VARCHAR(50) NOT NULL`)
* `name` (`VARCHAR(200) NOT NULL`)
* `address` (`VARCHAR(255) NULL`)
* `latitude` (`NUMERIC(9,6) NULL`)
* `longitude` (`NUMERIC(9,6) NULL`)
* `imported_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)


## 2. Relatii si Constrangeri
* `users` 1 : N `incidents` (un utilizator poate raporta mai multe sesizari; fiecare sesizare apartine unui singur utilizator). La stergerea unui cont, sesizarile asociate se sterg automat (`ON DELETE CASCADE`).
* Constrangeri de unicitate pe `email` si `username` pentru a preveni conturile duplicate.
* Toate interogarile folosesc parametri (`$1, $2, ...`) pentru a preveni SQL Injection.



## 3. Specificatia Endpoint-urilor API REST

Toate rutele returneaza raspunsuri in format JSON, cu exceptia descarcarii raportului CSV. Prefixul general este `/api`.

| Metoda | Endpoint | Acces | Descriere | Raspuns Succes | Erori uzuale |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/auth/register` | Public | Inregistreaza un utilizator nou (`role = citizen`). | `201 Created` | `400` (date invalide), `409` (email existent) |
| `POST` | `/api/auth/login` | Public | Autentifica utilizatorul si returneaza token-ul/sesiunea. | `200 OK` | `400` (campuri lipsa), `401` (date gresite) |
| `GET` | `/api/datasets` | Public | Returneaza datele din CSV. Suporta query params: `?category=...&search=...&limit=20&offset=0`. | `200 OK` | `400` (parametri invalizi) |
| `GET` | `/api/incidents` | Public | Listeaza sesizarile din oras. Suporta filtre: `?status=...&category=...`. | `200 OK` | `400` (filtru invalid) |
| `GET` | `/api/incidents/{id}` | Public | Returneaza detaliile unei sesizari dupa ID. | `200 OK` | `404` (inexistent) |
| `POST` | `/api/incidents` | Citizen | Creeaza o sesizare noua (status automat `INREGISTRAT`). | `201 Created` | `400` (campuri lipsa), `401` (neautentificat) |
| `GET` | `/api/users/me/incidents` | Citizen | Returneaza doar tichetele create de utilizatorul logat. | `200 OK` | `401` (neautentificat) |
| `DELETE`| `/api/incidents/{id}` | Citizen / Admin | Sterge o sesizare (permis cetateanului doar daca este `INREGISTRAT`). | `204 No Content` | `401` (nelogat), `403` (nepermis), `404` |
| `PATCH` | `/api/admin/incidents/{id}` | Admin | Modifica `status` si adauga explicatii in `admin_notes`. | `200 OK` | `400` (status invalid), `403` (nu e admin), `404` |
| `GET` | `/api/admin/incidents/export` | Admin | Descarca direct un fisier `.csv` cu sesizarile filtrate. | `200 OK` (`text/csv`)| `401` (nelogat), `403` (nu e admin) |
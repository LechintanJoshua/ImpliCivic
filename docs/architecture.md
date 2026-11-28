# Arhitectura Aplicatiei – ImpliCivic

## 1. Structura Generala
Aplicatia urmeaza o arhitectura modulara de tip **Client–Server pe trei niveluri**, conceputa pentru separarea clara a responsabilitatilor:
- **Client (Frontend):** Interfata web SPA in browser, complet decuplata de server, care gestioneaza afisarea datelor, validarea initiala a formularelor si comunicarea asincrona cu serverul in format JSON.
- **Server (Backend in Go):** Gestioneaza cererile concurente, politica de acces cross-origin (CORS), autentificarea si controlul accesului pe roluri, sanitizarea datelor, procesarea fisierelor CSV si comunicarea cu baza de date.
- **Baza de date (PostgreSQL) si Stocarea locala CSV:** PostgreSQL asigura persistenta utilizatorilor, a sesizarilor si a punctelor de interes public; directorul local `data/` stocheaza seturile deschise in format CSV pentru import.

---

## 2. Diagrame Arhitecturale
Reprezentarile vizuale ale sistemului si ale structurii bazei de date sunt disponibile in format PDF:
- [Diagrama Arhitecturala (PDF)](../architecture/Arhitectura.pdf)
- [Diagrama Bazei de Date (PDF)](../architecture/Database.pdf)

---

## 3. Descrierea Componentelor si Comunicarea dintre Ele

### 3.1. Stratul Client (Frontend)
- **Module de Interfata:** Componente dedicate pentru consultarea catalogului de date publice, raportarea unei sesizari noi, vizualizarea istoricului personal de tichete si panoul de moderare al administratorului.
- **Gestionarea sesiunii:** Pastreaza token-ul de autentificare primit de la server si il ataseaza cererilor protejate; adapteaza optiunile vizibile in interfata in functie de rolul utilizatorului (`citizen` sau `admin`).
- **Comunicare cu serverul:** Toate schimburile de date se fac prin cereri asincrone in format JSON (respectiv descarcare directa de fisiere pentru raportul CSV).

### 3.2. Stratul Server (Backend in Go)
Backend-ul este organizat in pachete interne independente, astfel incat logica de retea (HTTP), logica criptografica, procesarea de fisiere si accesul la baza de date sa nu fie amestecate:

- **Stratul HTTP si Middleware (`cmd/api`):**
  - Porneste serverul cu timeout-uri configurate pentru citire, scriere si conexiuni inactive, tratand fiecare cerere intr-o goroutine separata.
  - Aplica filtrul **CORS** pentru a permite cererile provenite din originea configurata a frontend-ului si a raspunde direct la cererile de verificare preflight trimise de browser.
  - Aplica filtrele de **autentificare si autorizare pe roluri**, atasand identitatea utilizatorului in contextul cererii si respingand cererile neautorizate sau pe cele de administrare venite de la conturi obisnuite.
  - Limiteaza dimensiunea datelor JSON primite, valideaza campurile obligatorii si lungimile textelor, si aplica **curatarea anti-XSS** inainte de salvare.

- **Stratul de Securitate (`internal/auth`):**
  - Complet decuplat de baza de date si de HTTP.
  - Se ocupa de hash-uirea parolelor cu salt aleatoriu si cost ridicat (`bcrypt`) si de generarea, semnarea si verificarea token-urilor de sesiune (`HMAC-SHA256`), folosind comparatii in timp constant si verificare a expirarii.

- **Stratul de Procesare CSV (`internal/csvproc`):**
  - **Import concurent la pornire:** Citeste fisierele CSV linie cu linie (streaming, fara a incarca intregul fisier in memorie) si distribuie randurile printr-un canal catre un grup de goroutines (worker pool) sincronizate. Randurile corupte sau cu campuri lipsa sunt logate si sarite defensiv fara a opri serverul.
  - **Export direct:** Genereaza continutul CSV direct in fluxul de raspuns HTTP, fara a crea fisiere temporare pe disc.

- **Stratul de Acces la Date (`internal/store`):**
  - Izoleaza complet interogarile SQL de stratul HTTP, lucrand exclusiv cu modele de date si parametri SQL legati (`$1, $2, ...`) pentru protectia impotriva SQL Injection.
  - Traduce erorile specifice bazei de date (ex. inregistrari inexistente sau duplicate de unicitate) in erori clare pentru stratul HTTP.

- **Stratul de Baza de Date (`internal/db`):**
  - Initializeaza si configureaza un **pool comun de conexiuni** reutilizabile intre cererile concurente, cu limite de conexiuni active, conexiuni inactive si durata maxima de viata.
  - Verifica disponibilitatea serverului PostgreSQL la pornire cu timeout scurt, oprind pornirea daca baza de date nu raspunde.

### 3.3. Stratul de Persistenta (PostgreSQL)
- Stocheaza permanent tabelele `users`, `incidents` si `public_datasets`, impunand la nivel de schema constrangeri de cheie externa cu stergere in cascada, unicitate pe conturi, precizie fixa pentru coordonate GPS si validari pe valorile permise de rol si status.

---

## 4. Structura de Directoare Implementata

```text
ImpliCivic/
├── cmd/
│   └── api/                  # Punctul de intrare, server HTTP, middleware si handlere
├── internal/
│   ├── auth/                 # Hash-uire parole si token-uri JWT
│   ├── csvproc/              # Import concurent CSV si generator export CSV
│   ├── db/                   # Schema SQL si pool-ul de conexiuni PostgreSQL
│   └── store/                # Stratul de interogari SQL izolat pe entitati
├── data/                     # Fisiere CSV locale cu date publice din Timisoara
├── docs/                     # Documentatia tehnica, modele si testare
├── architecture/             # Diagramele PDF ale arhitecturii si bazei de date
└── mockup/                   # Mockup-ul UI in format PDF
```

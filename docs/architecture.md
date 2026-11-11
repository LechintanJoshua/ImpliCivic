# Arhitectura Aplicatiei – ImpliCivic

## 1. Structura Generala
Aplicatia are o arhitectura simpla de tip **Client-Server**:
- **Client (Frontend):** Interfata web ruleaza integral in browserul utilizatorului si gestioneaza afisarea, interactiunea si cererile asincrone prin API REST.
- **Server (Backend in Go):** Gestioneaza cererile HTTP, verificarile de acces, logica de prelucrare si comunicarea cu baza de date.
- **Baza de date (PostgreSQL):** Asigura salvarea persistenta a utilizatorilor, a sesizarilor si a punctelor publice importate.
- **Stocare locala seturi deschise (CSV):** Fisiere locale cu date publice din Timisoara, procesate la nivel de server.



## 2. Diagrama Arhitecturala
Schema vizuala a arhitecturii aplicatiei (interactiunea dintre client, server si baza de date) este disponibila in format PDF in directorul dedicat:
- 📄 **[Diagrama Arhitecturala (PDF)](diagrams/arhitectura.pdf)**



## 3. Descrierea Componentelor

### 3.1. Frontend
- **Module de Interfata:** Componente separate care corespund paginilor din aplicatie:
  - *Catalog Date Deschise:* Afisare tabelara, filtrare dupa categorie si cautare dupa strada/denumire.
  - *Formular Raportare:* Formular pentru sesizari, preluare coordonate GPS si validare campuri.
  - *Tichetele Mele:* Lista sesizarilor trimise de utilizator, cu eticheta pentru fiecare stare (`Inregistrat`, `In lucru`, `Rezolvat`).
  - *Panou Administrator:* Tabel pentru moderarea sesizarilor, schimbare de status si buton de export raport.
- **Validare pe Client:** Verifica datele inainte de trimitere (campuri obligatorii completate, text existent, coordonate valide) pentru a da erori rapide direct in pagina.
- **Sesiune si Autentificare:** Retine token-ul primit la login si rolul utilizatorului, `citizen` sau `admin`, afisand optiunile potrivite in meniu.
- **Serviciu API:** Modul comun care se ocupa de apelurile `fetch()`, trimite datele in format JSON si prinde raspunsurile sau erorile venite de la server.

### 3.2. Backend
- Toate rutele aplicatiei incep cu prefixul `/api`.
- **Rutare HTTP:** Trimite cererile primite catre functia potrivita pe baza metodei (GET, POST, PATCH) si a URL-ului.
- **Autentificare si Drepturi:** Verifica daca utilizatorul este autentificat si ce rol are. Daca un cont obisnuit de cetatean incearca sa apeleze rute de administrare (`/api/admin/*`), serverul refuza cererea direct cu codul `403 Forbidden`.

### 3.3. Handlere si Logica de Aplicatie
- **Handlere HTTP:** Preiau datele din cerere (JSON), verifica prezenta campurilor obligatorii, curata textul de tag-uri periculoase pentru protectia de baza anti-XSS si apeleaza functiile corespunzatoare pentru salvare.
- **Procesare fisiere CSV:** Citirea fisierelor mari cu date locale se face linie cu linie, folosind procesare paralela cu goroutines pentru a nu bloca serverul. Randurile corecte sunt salvate in baza de date, iar cele cu erori sunt sarite.
- **Export CSV:** Citeste sesizarile filtrate din PostgreSQL si genereaza pe loc un fisier `.csv` pe care administratorul il poate descarca direct din browser.

### 3.4. Baza de Date
- **Organizarea interogarilor:** Functiile care lucreaza direct cu tabelele (`users`, `incidents`, `public_datasets`) sunt tinute separat de logica de rute pentru a pastra codul ordonat.
- **Protectie SQL Injection:** Toate interogarile folosesc parametri (`$1, $2, ...`), astfel incat textul trimis de utilizatori sa nu poata modifica comanda SQL.
- **Conexiune stabila:** Serverul pastreaza un pool comun de conexiuni catre PostgreSQL, le refoloseste fara a deschide o conexiune noua la fiecare click.



## 4. Structura de Directoare a Proiectului

```text
implicivic/
├── cmd/
│   └── api/
├── internal/
│   ├── auth/
│   ├── csvproc/
│   ├── db/
│   └── store/
├── web/
│   ├── css/
│   └── js/
├── data/
└── docs/
    ├── diagrams/
    └── wireframes/
```
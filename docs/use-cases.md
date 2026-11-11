# Scenarii de Utilizare

### 1. Raportarea unui incident urban
**Actor:** Cetatean (utilizator autentificat).  
**Preconditii:** utilizator autentificat cu rol de `citizen`; formularul de sesizare este accesibil.  
**Pasi:**
1. Cetateanul acceseaza formularul de sesizare noua.
2. Completeaza titlul, categoria (ex: carosabil, deseuri, iluminat), adresa/locatia si descrierea problemei, optional coordonate GPS.
3. Trimite formularul prin `POST /api/incidents`.
4. Backend-ul verifica token-ul, verifica textul sa nu existe XSS si valideaza campurile obligatorii.
5. Sesizarea este salvata in PostgreSQL cu statusul implicit `INREGISTRAT` si asociata contului curent.
6. Interfata afiseaza confirmarea si redirectioneaza catre istoricul personal de tichete.  
**Rezultat:** o inregistrare noua in tabela `incidents` legata prin cheie externa de `user_id`, cu statusul initial `INREGISTRAT`.  
**Verificare:** la trimiterea unei sesizari valide, serverul raspunde cu `201 Created`, iar tichetul apare imediat in `GET /api/users/me/incidents` avand statusul `INREGISTRAT`.  
**Alternative/erori:** titlu sau descriere goala: `400 Bad Request`, fara salvare in baza de date. Titlu mai lung de 200 de caractere: `400 Bad Request`.
Utilizator neautentificat: `401 Unauthorized`.

---

### 2. Consultarea si filtrarea datelor publice (CSV)
**Actor:** Vizitator anonim sau Cetatean.  
**Preconditii:** fisierele CSV au fost procesate la pornirea serverului sau sincronizate; baza de date contine inregistrari publice.  
**Pasi:**
1. Utilizatorul acceseaza sectiunea publica de date deschise.
2. Interfata solicita datele initiale prin `GET /api/datasets`.
3. Utilizatorul selecteaza o categorie (ex: "reciclare") si introduce o strada in campul de cautare (ex: "Cetatii").
4. Frontend-ul trimite cererea filtrata: `GET /api/datasets?category=reciclare&search=Cetatii`.
5. Backend-ul ruleaza interogarea si returneaza lista rezultatelor in format JSON.
6. Interfata actualizeaza dinamic tabelul cu punctele gasite, fara a reincarca pagina.  
**Rezultat:** utilizatorul vizualizeaza in tabel datele publice extrase din CSV conform filtrelor alese.  
**Verificare:** daca in baza de date exista 8 puncte de reciclare pe strada specificata, raspunsul serverului este `200 OK`, iar tabelul din pagina afiseaza exact 8 randuri.  
**Alternative/erori:** niciun rezultat gasit: raspuns `200 OK` cu lista goala `[]` si mesaj in interfata "Nu s-au gasit date conform filtrelor".
Parametri incorecti: `400 Bad Request`.

---

### 3. Moderarea unei sesizari si actualizarea statusului
**Actor:** Administrator.  
**Preconditii:** utilizator autentificat cu rol de `admin`; sesizarea vizata exista in baza de date.  
**Pasi:**
1. Administratorul acceseaza panoul de gestiune si selecteaza un incident din lista.
2. Modifica statusul, din `INREGISTRAT` in `IN_LUCRU`, `REZOLVAT` sau `RESPINS`, si adauga o nota explicativa.
3. Trimite datele prin `PATCH /api/admin/incidents/{id}`.
4. Backend-ul verifica drepturile de administrator si valideaza valoarea noului status.
5. Inregistrarea este actualizata in baza de date: `status`, `admin_notes`, `updated_at`.
6. Interfata afiseaza confirmarea actualizarii.  
**Rezultat:** sesizarea reflecta noul status atat in panoul administrativ, cat si in profilul cetateanului care a raportat-o.  
**Verificare:** dupa salvare, un apel `GET /api/incidents/{id}` returneaza noul status, ex `IN_LUCRU`, si textul adaugat in `admin_notes`.  
**Alternative/erori:** identificator inexistent: `404 Not Found`.
Valoare nepermisa pentru status: `400 Bad Request`.
Cerere trimisa de un utilizator cu rol `citizen`: `403 Forbidden`, fara modificari in baza de date.

---

### 4. Exportul sesizarilor in format CSV
**Actor:** Administrator.  
**Preconditii:** utilizator autentificat cu rol de `admin`; exista sesizari in baza de date.  
**Pasi:**
1. Administratorul aplica un filtru in panoul de control (ex: doar sesizarile cu statusul `REZOLVAT`).
2. Apasa butonul "Exporta CSV".
3. Frontend-ul initiaza cererea `GET /api/admin/incidents/export?status=REZOLVAT`.
4. Backend-ul verifica rolul de `admin`, extrage inregistrarile filtrate din PostgreSQL si genereaza dinamic continutul CSV folosind pachetul `encoding/csv`.
5. Serverul trimite raspunsul cu headerele `Content-Type: text/csv` si `Content-Disposition: attachment; filename="sesizari_export.csv"`.
6. Browserul descarca automat fisierul pe statia locala.  
**Rezultat:** un fisier CSV generat la cerere, salvat local pe statia de lucru a administratorului.  
**Verificare:** la un numar de 15 sesizari cu statusul `REZOLVAT`, fisierul CSV descarcat are exact 16 linii (1 linia de antet + 15 linii de date).  
**Alternative/erori:** utilizator simplu `citizen` sau neautentificat: `403 Forbidden` / `401 Unauthorized`, fara generare de fisier.
Nu exista inregistrari pentru filtrul cerut: se genereaza un fisier CSV valid ce contine doar linia de antet cu denumirile coloanelor.
# ImpliCivic

## 1. Descrierea proiectului si problema abordata
**ImpliCivic** este o platforma web prin care aducem datele publice locale intr-un singur loc si oferim un mod simplu prin care locuitorii din Timisoara pot raporta probleme intalnite in oras.

In momentul de fata, datele puse la dispozitie la nivel local, cum ar fi locatiile parcarilor, punctele de colectare a deseurilor sau diverse lucrari, se gasesc in tabele sau fisiere CSV greu de parcurs. Totodata, sesizarile oamenilor legate de gropi, iluminat stradal defect sau masini parcate ilegal ajung de obicei pe retelele sociale, unde se pierd rapid si nu au un flux de rezolvare.

## 2. Obiectivele proiectului
- Parsarea si afisarea datelor din fisiere CSV intr-o interfata web usor de citit.
- Oferirea unui sistem prin care cetatenii logati pot trimite sesizari si pot urmari stadiul acestora.
- Crearea unui panou de administrare pentru actualizarea statusurilor si exportul sesizarilor inapoi in format CSV.

## 3. Utilizatori vizati
- **Vizitatori anonimi:** Oricine vrea doar sa consulte datele sau sa vada incidentele deja raportate, fara sa aiba cont.
- **Cetateni autentificati:** Locuitorii care vor sa trimita sesizari concrete si sa vada istoricul tichetelor proprii.
- **Administratori:** Cei care modereaza sesizarile, le schimba stadiul pe masura ce se intervine si descarca rapoarte centralizate.

## 4. Echipa si responsabilitati
- **Lechintan Joshua-Leonardo-Augustin** – *Arhitectura Backend & Persistenta Datelor*: Configurarea serverului in Go, structura modulara, schema si migratiile in PostgreSQL, procesarea concurenta a fisierelor CSV, securitatea server-side (validari, prevenire SQLi/XSS, roluri) si administrarea repository-ului de GitHub.
- **Losnariu Marius-Madalin** – *Frontend & UI/UX Design*: Structura paginilor in HTML5, stilizarea responsive cu CSS pentru mobil si desktop, designul interfetei, interactivitatea acesteia si respectarea cerintelor de accesibilitate.
- **Teodorovits Andras-Karoly** – *Integrare Full-Stack & Logica Client*: Conectarea frontend-ului la API prin cereri asincrone, validarea formularelor pe client, implementarea endpoint-urilor pentru filtrarea datelor si logica de export CSV pe backend, alaturi de testarea fluxurilor aplicatiei.

## 5. Functionalitati principale
1. **Consultarea si filtrarea datelor publice:** Citirea fisierelor CSV, procesarea lor pe server si afisarea datelor intr-un tabel cu optiuni de cautare si filtrare pe categorii sau strazi.
2. **Raportarea incidentelor:** Formular prin care utilizatorii autentificati trimit sesizari specificand titlu, categorie, descriere si locatie.
3. **Urmarirea stadiului:** Pagina de profil unde utilizatorul vede starea curenta a fiecarei sesizari trimise (`Inregistrat`, `In lucru`, `Rezolvat`, `Respins`).
4. **Panou de administrare si Export CSV:** Modul protejat pentru admin unde se pot actualiza tichetele si se poate genera o descarcare directa a raportului in format CSV.

## 6. Tehnologii alese
- **Backend:** Go.
- **Baza de date:** PostgreSQL.
- **Frontend:** HTML5, CSS3, TypeScript, React.
- **Unelte:** Git, GitHub, Postman.

## 7. Documentatia proiectului
Toate detaliile tehnice se gasesc in directorul `docs/`:
- [Cerinte Functionale si Nefunctionale](docs/requirements.md)
- [Scenarii de Utilizare (Use Cases)](docs/use-cases.md)
- [Wireframe-uri si Navigare](docs/wireframes/README.md)
- [Arhitectura Aplicatiei](docs/architecture.md)
- [Modelul de Date si Endpoint-uri API](docs/data-model.md)
- [Planul de Implementare](docs/implementation-plan.md)
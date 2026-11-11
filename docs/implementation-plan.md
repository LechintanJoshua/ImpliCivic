# Plan de Implementare – ImpliCivic

## 1. Tehnologii Alese

- **Backend (Go):** Biblioteca standard (`net/http`, `encoding/csv`), sintaxa simpla fara framework-uri greoaie si suport nativ pentru goroutines la procesarea paralela a fisierelor CSV.
- **Baza de date (PostgreSQL):** Model relational stabil, integritate prin chei externe (`users` 1 : N `incidents`) si interogari parametrizate ($1, $2) impotriva atacurilor de tip SQL Injection.
- **Frontend (HTML5, CSS3, TypeScript, React):** Aplicatie Single Page (SPA) responsive, complet decuplata de server si conectata prin cereri asincrone `fetch()` catre API.



## 2. Organizare Git

Lucram pe GitHub prin feature branches simple:
- Nu facem commit-uri direct pe `main`.
- Cream branch-uri scurte pe functionalitati.
- Modificarile ajung in `main` prin Pull Request cu verificare intre noi.



## 3. Planificare pe Etape

- **Etapa 2 (Saptamana 7 – Curenta):** Documentatia tehnica.
- **Etapa 3 (Saptamanile 8 – 12):** Prototip functional:
  - Initializarea bazei de date PostgreSQL si crearea tabelelor.
  - Serverul cu rutele de baza (`/api/auth`, `/api/incidents`, `/api/datasets`).
  - Worker-ul de citire si parsare streaming pentru fisierele CSV.
  - Interfata React conectata la rutele de raportare si vizualizare.
- **Etapa 4 (Saptamanile 12 – 14):** Finalizare si predare:
  - Panoul de administrare si functia de export a rapoartelor in format CSV.
  - Validari de securitate anti-XSS pe server, filtrare parametri.
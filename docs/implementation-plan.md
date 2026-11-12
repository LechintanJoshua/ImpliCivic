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
- **Etapa 3 (Saptamanile 8 – 10):** Prototip functional End-to-End:
  - Initializarea bazei de date PostgreSQL si crearea tabelelor.
  - Serverul Go cu rutele de baza pentru operatii CRUD (`/api/auth`, `/api/incidents`, `/api/datasets`).
  - Algoritmul de citire si parsare pentru fisierele CSV.
  - Interfata React conectata, permitand efectuarea unor operatii de salvare in baza de date.
- **Etapa 4 (Saptamanile 10 – 14):** Finalizare si predare:
  - Panoul de administrare cu functionalitati restrictionate si mecanismul de export al rapoartelor in format CSV.
  - Validari de securitate anti-XSS pe server, filtrare parametri si protejarea credentialelor.
  - Ajustari pe frontend pentru accesibilitate, design responsive si validarea datelor.
# Cerinte Tehnice

## 1. Cerinte Functionale

### Obligatorii
- **Consultare date deschise:** Orice vizitator poate vedea intr-un tabel datele publice din fisierele CSV (parcari, colectare deseuri etc.) si le poate filtra dupa categorie sau strada.
- **Raportare incident:** Un utilizator logat poate trimite o sesizare completand titlu, categorie, descriere si adresa din teren.
- **Autentificare si roluri:** Sistem de conturi care separa utilizatorii simpli (cetateni) de administratori.
- **Gestiune si actualizare status:** Administratorul poate schimba starea unui tichet (Inregistrat, In lucru, Rezolvat, Respins).
- **Export rapoarte in CSV:** Administratorul poate descarca sesizarile filtrate direct intr-un fisier CSV generat de server.

### Optionale
- **Profil cetatean:** Pagina unde utilizatorul isi poate vedea istoricul propriilor sesizari si statusul lor.
- **Upload CSV:** Buton in panoul de admin pentru a incarca un fisier CSV nou direct din browser.



## 2. Cerinte Nefunctionale

### Obligatorii
- **Securitate server-side:** Folosirea de interogari SQL parametrizate pentru a opri SQL Injection si curatarea textului impotriva XSS.
- **Protectia credentialelor:** Parolele salvate doar ca hash, iar parola de DB tinuta in variabile de mediu locale, nu pe GitHub.
- **Design responsive:** Interfata usor de folosit atat pe mobil cat si pe desktop.
- **Accesibilitate de baza:** Paginile au structura clara, etichete pe fiecare camp de formular si merg navigate din tastatura (Tab).
- **Parsare defensiva CSV:** Daca un rand din fisierul CSV e stricat sau lipsesc coloane, serverul il sare si continua sa citeasca restul fisierului fara sa pice.

### Optionale
- **Performanta la parsare:** Folosirea de goroutines pentru a procesa fisiere CSV mari in paralel fara sa blocheze alte cereri.
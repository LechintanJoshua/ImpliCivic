CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'citizen' CHECK (role IN ('citizen', 'admin')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW() -- include fusul orar, l-am pus pentru a nu avea probleme la timp intre platforme
);

CREATE TABLE IF NOT EXISTS incidents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE, -- relatie de 1:N. O data sters un cont, se sterg automat si sesizarile lui
    title VARCHAR(200) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    address VARCHAR(255) NOT NULL,
    latitude NUMERIC(9, 6) NULL, -- precizie zecimala pentru GPS, am pus asta pentru a evita erorile de rotunjire din FLOAT
    longitude NUMERIC(9, 6) NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'INREGISTRAT' CHECK (status IN ('INREGISTRAT', 'IN_LUCRU', 'REZOLVAT', 'RESPINS')),
    admin_notes TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- date publice accesibile anonim
CREATE TABLE IF NOT EXISTS public_datasets ( 
    id BIGSERIAL PRIMARY KEY,
    dataset_name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    address VARCHAR(255) NULL,
    latitude NUMERIC(9, 6) NULL,
    longitude NUMERIC(9, 6) NULL,
    imported_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS public."residentials" (
    id serial not null, 
    name text UNIQUE,
    lat float,
    long float,
    CONSTRAINT "res_pk" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."types" (
    id serial not null, 
    name text UNIQUE,
    CONSTRAINT "types_pk" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."statuses" (
    id serial not null, 
    name text UNIQUE,
    CONSTRAINT "statuses_pk" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."users" (
    id serial not null, 
    score integer,
    fid text UNIQUE,
    name text,
    CONSTRAINT "users_pk" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."items" (
    id serial not null, 
    scan text UNIQUE,
    type_id integer,
    CONSTRAINT "items_pk" PRIMARY KEY (id),
    CONSTRAINT "type_fk" FOREIGN KEY (type_id) REFERENCES types(id)
);

CREATE TABLE IF NOT EXISTS public."history" (
    id serial not null, 
    uid integer,
    rid integer, 
    type_id integer,
    amount integer,
    date_start TIMESTAMP,
    status_id integer,
    CONSTRAINT "user_fk" FOREIGN KEY (uid) REFERENCES public."user"(id),
    CONSTRAINT "res_fk" FOREIGN KEY (rid) REFERENCES public."residentials"(id),
    CONSTRAINT "type_fk" FOREIGN KEY (type_id) REFERENCES public."types"(id),
    CONSTRAINT "status_fk" FOREIGN KEY (status_id) REFERENCES public."statuses"(id),
    CONSTRAINT "history_pk" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."coordinate" (
    id integer, 
    lat float, 
    long float, 
    CONSTRAINT "coordinate_pk" PRIMARY KEY (id)
);

INSERT INTO public."types"(name) VALUES('plastic');
INSERT INTO public."types"(name) VALUES('paper');
INSERT INTO public."types"(name) VALUES('metal');

INSERT INTO public."statuses"(name) VALUES('sorted');
INSERT INTO public."statuses"(name) VALUES('inbin');
INSERT INTO public."statuses"(name) VALUES('processing');
INSERT INTO public."statuses"(name) VALUES('completed');

INSERT INTO public."residentials"(name, lat, long) VALUES('factory', 0.0, 0.0)
CREATE TYPE user_status AS ENUM ('INACTIVE', 'ACTIVE', 'BANNED');
CREATE TABLE IF NOT EXISTS public.player
(
    email character varying(50) NOT NULL PRIMARY KEY,
    name character varying(50) NOT NULL,
    status user_status NOT NULL DEFAULT 'INACTIVE'::user_status
);

CREATE TABLE IF NOT EXISTS public.wallet
(
    wallet_id character varying(50) NOT NULL PRIMARY KEY,
    user_email character varying(50) NOT NULL REFERENCES player(email),
    balance bigint NOT NULL DEFAULT 0,
    create_time bigint NOT NULL DEFAULT floor((EXTRACT(epoch FROM now()) * (1000)::numeric)),
    update_time bigint NOT NULL DEFAULT floor((EXTRACT(epoch FROM now()) * (1000)::numeric))
);

CREATE TABLE IF NOT EXISTS public.record
(
	record_id character varying(50) NOT NULL PRIMARY KEY,
    operation_trace_id character varying(50) NOT NULL,
    user_email character varying(50) NOT NULL REFERENCES player(email),
	wallet_id character varying(50) NOT NULL REFERENCES wallet(wallet_id),
    amount bigint NOT NULL,
	balance bigint NOT NULL,
	create_time bigint NOT NULL DEFAULT floor((EXTRACT(epoch FROM now()) * (1000)::numeric))
);
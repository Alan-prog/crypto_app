-- crate a table to save the user data
create table user_data
(
	id serial not null,
	name varchar(256) not null,
	last_name varchar(256) not null,
	email varchar(256) not null,
	pass_hash varchar(256) not null,
	access_token varchar(256),
	refresh_token varchar(256)
);

create unique index user_data_id_uindex
	on user_data (id);

alter table user_data
	add constraint user_data_pk
		primary key (id);

create unique index user_data_email_uindex
	on user_data (email);

alter table user_data drop column access_token;

alter table user_data drop column refresh_token;

-- create table with salary
create table salary
(
	id serial not null,
	name varchar(10) not null,
	cost float not null
);

create unique index salary_id_uindex
	on salary (id);

alter table salary
	add constraint salary_pk
		primary key (id);

insert into salary (name, cost) values ('BTC', 32853.856);
insert into salary (name, cost) values ('ETH', 2022.65);

-- create table addresses to store the unique address of the currency

create table addresses
(
	id serial not null,
	address varchar(256) not null
);

create unique index addresses_id_uindex
	on addresses (id);

alter table addresses
	add constraint addresses_pk
		primary key (id);

alter table addresses
	add user_id integer not null;

alter table addresses
	add salary varchar(10) not null;

alter table addresses
	add balance float not null;

alter table addresses
	add constraint addresses_user_data_id_fk
		foreign key (user_id) references user_data;

alter table addresses rename column salary to salary_id;

alter table addresses alter column salary_id type integer using salary_id::integer;

alter table addresses
	add constraint addresses_salary_id_fk
		foreign key (salary_id) references salary;

-- create table for transactions
create table transactions
(
	from_address integer not null
		constraint transactions_addresses_id_fk
			references addresses,
	to_address integer not null
		constraint transactions_addresses_id_fk_2
			references addresses,
	from_currency varchar(10) not null,
	to_currency varchar(10) not null,
	amount_dollars float not null,
	create_at timestamp default current_timestamp not null,
	commission float not null
);

alter table transactions drop column from_currency;

alter table transactions drop column to_currency;

alter table transactions
	add successful bool default false not null;

-- create a transaction for send money from one wallet to other with stores procedures
create or replace function make_transaction (
    first_address_id integer,
    last_address_id integer ,
    amount float ,
    commission float
)
returns table (
	response bool
)
language plpgsql
as $$
declare
    first_update integer;
    last_update integer;
    firstCost float;
    lastCost float;
begin
    select s.cost from addresses as a
        left join salary s on a.salary_id = s.id
    where a.id = first_address_id into firstCost;
    select s.cost from addresses as a
        left join salary s on a.salary_id = s.id
    where a.id = last_address_id into lastCost;

    PERFORM balance from addresses where id = first_address_id OR id = last_address_id for update;
    UPDATE addresses SET balance = balance - (amount/firstCost)/(1 - commission) WHERE id = first_address_id
            and balance >= (amount / firstCost)/(1 - commission)
    RETURNING id into first_update;
    UPDATE addresses SET balance = balance + (amount/lastCost)  WHERE id = last_address_id and
            first_update is not null
    returning id into last_update;

    INSERT INTO transactions (from_address, to_address, amount_dollars, commission, successful)
        values(first_address_id,last_address_id,amount,commission, last_update is not null) returning successful
            into response;
    return query (select response as response);
end; $$
create table "user"
(
    id           serial
        primary key,
    name         varchar(20) not null,
    surname      varchar(30),
    patronymic   varchar(20),
    password     varchar(30) not null,
    nick         varchar(20),
    phone        varchar(10) not null,
    vk           varchar(50),
    tg           varchar(50),
    email        varchar(50),
    private_data varchar(400),
    apply        timestamp,
    fired        timestamp
);

alter table "user"
    owner to postgres;

create table event
(
    id                serial
        primary key,
    name              varchar(45),
    discription       varchar(100),
    order_discription varchar(400),
    start_date        date,
    start_time        time,
    end_date          date,
    ned_time          time
);

alter table event
    owner to postgres;

create table event_roster
(
    id      integer default nextval('event_roaster_id_seq'::regclass) not null
        constraint event_roaster_pkey
            primary key,
    userid  integer
        constraint event_roaster_userid_fkey
            references "user",
    eventid integer
        constraint event_roaster_eventid_fkey
            references event,
    leader  boolean
);

alter table event_roster
    owner to postgres;

create table inventory
(
    id          serial
        primary key,
    name        varchar(50),
    type        varchar(20),
    description varchar(200),
    "uniqNum"   integer
);

alter table inventory
    owner to postgres;

create table crew_inventory
(
    id      serial
        primary key,
    item_id integer
        references inventory
);

alter table crew_inventory
    owner to postgres;

create table car
(
    id        serial
        primary key,
    name      varchar(40),
    comment   varchar(400),
    inventory integer
        references crew_inventory
);

alter table car
    owner to postgres;

create table day_crew
(
    id         serial
        primary key,
    comment    varchar(400),
    leader     integer not null
        constraint leader
            references "user",
    time_start timestamp,
    time_end   timestamp
);

alter table day_crew
    owner to postgres;

create table day_crew_roster
(
    user_id integer
        constraint date_crew_roaster_user_fkey
            references "user",
    crew_id integer
        constraint date_crew_roaster_day_crew_fkey
            references day_crew
);

alter table day_crew_roster
    owner to postgres;

create table calls
(
    id          serial
        primary key,
    description varchar(200),
    summing_up  varchar(200),
    address     varchar(45),
    time_start  timestamp,
    time_finish timestamp,
    title       varchar(50)
);

alter table calls
    owner to postgres;

create table crew_calls
(
    call_id integer
        references calls,
    crew_id integer
        references day_crew
);

alter table crew_calls
    owner to postgres;

create table permission
(
    id          serial
        primary key,
    action      boolean,
    action_name varchar(50)
);

alter table permission
    owner to postgres;

create table ranks
(
    id   serial
        primary key,
    name varchar(20)
);

alter table ranks
    owner to postgres;

create table rank_permissions
(
    permission integer
        references permission,
    rank       integer
        references ranks
);

alter table rank_permissions
    owner to postgres;

create table codes
(
    id     serial
        primary key,
    code   varchar(10),
    userid integer
        references "user"
);

alter table codes
    owner to postgres;

create table event_inventory
(
    eventid integer
        references event,
    itemid  integer
        references inventory
);

alter table event_inventory
    owner to postgres;

create table users_ranks
(
    user_id integer
        references "user",
    rank_id integer
        references ranks
);

alter table users_ranks
    owner to postgres;


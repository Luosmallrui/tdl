create table tasks
(
    id          bigint unsigned auto_increment
        primary key,
    title       varchar(255)                  not null,
    description text                          null,
    status      varchar(20) default 'pending' null,
    user_id     bigint unsigned               null,
    due_date    datetime(3)                   null,
    reminder_at datetime(3)                   null,
    created_at  datetime(3)                   null,
    updated_at  datetime(3)                   null,
    tags        longtext                      null
);

create index idx_tasks_user_id
    on tasks (user_id);

create table users
(
    id         bigint unsigned auto_increment
        primary key,
    username   longtext    null,
    nickname   longtext    null,
    email      longtext    null,
    password   longtext    null,
    created_at datetime(3) null,
    updated_at datetime(3) null
);


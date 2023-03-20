-- +migrate Up
create table action
(
    id            text  not null
                  primary key ,
    type          text not null,
    status        text not null,
    data   text not null
);

-- +migrate Down
drop table action;
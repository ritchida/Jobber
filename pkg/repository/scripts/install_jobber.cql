create keyspace if not exists jobber with replication = {'class':'SimpleStrategy', 'replication_factor' : 1};

use jobber;

create table if not exists jobs
(
    job_id timeuuid primary key,
    created timestamp,
    last_updated timestamp,
    completed timestamp,
    status text,
    type text,
    tags set<text>,
    owner text
);

create table if not exists latest_jobs
(
    bucket int,
    job_id timeuuid,
    created timestamp,
    last_updated timestamp,
    completed timestamp,
    status text,
    type text,
    tags set<text>,
    owner text,
    primary key ((bucket), job_id)
) with clustering order by (job_id desc);

create table if not exists job_messages_by_job_id
(
    job_id timeuuid,
    message_created timestamp,
    message text,
    primary key ((job_id), message_created)
);

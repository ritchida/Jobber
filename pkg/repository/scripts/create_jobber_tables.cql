use jobber;

create table jobs
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

create table latest_jobs
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

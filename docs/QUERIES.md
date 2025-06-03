# Exporters Queries

- ### parts_log:
```sql
select 
    database, 
    table, 
    sum(bytes) as bytes, 
    count() as parts, 
    sum(rows) as rows 
from system.parts
{FILTER_CLAUSE} 
group by database, table
```

- ### event_log:
```sql
select 
    event, 
    value 
from system.events
{FILTER_CLAUSE}
```

- ### disk_log:
```sql
select 
    name, 
    sum(free_space) as free_space_in_bytes, 
    sum(total_space) as total_space_in_bytes 
from system.disks 
{FILTER_CLAUSE}
group by name
```

- ### basic_log:
```sql
select 
    metric, 
    value 
from system.metrics
{FILTER_CLAUSE}
```

- ### async_log:
```sql
select 
    replaceRegexpAll(toString(metric), '-', '_') AS metric,
    value 
from system.asynchronous_metrics
{FILTER_CLAUSE}
```


- ### query_log:
```sql
SELECT 
    user, 
    type as status,
    query_kind,
    arrayJoin(tables) AS table, 
    sum(memory_usage) as memory_usage,
    count(*) AS query_num,
    sum(query_duration_ms) as query_duration_ms,
    sum(read_bytes) as read_bytes,
    sum(read_rows) as read_rows,
    sum(written_bytes) as written_bytes,
    sum(written_rows) as written_rows,
    sum(result_bytes) as result_bytes,
    sum(result_rows) as result_rows,
    sum(peak_threads_usage) as peak_threads_usage
FROM system.query_log
{FILTER_CLAUSE}
GROUP BY user, table, type,query_kind
```
CREATE USER 'exporter'@'%' IDENTIFIED BY '1f1te123fq';
GRANT PROCESS, REPLICATION CLIENT ON *.* TO 'exporter'@'%';
GRANT SELECT ON performance_schema.* TO 'exporter'@'%';
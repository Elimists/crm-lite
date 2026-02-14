Notes

Windows

Get-Content ./.local/migrations_pg.sql | docker exec -i postgres-dev psql -U pran -d crm_lite

Get-Content ./.local/seed_pg_dev.sql | docker exec -i postgres-dev psql -U pran -d crm_lite
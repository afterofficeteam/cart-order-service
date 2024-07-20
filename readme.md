
## How to run
Rename sample.config yaml to config.yaml

Clone my project in local, running on terminal / cli in project directery

```bash
  go run main.go
```

### Run Migration
Don't forget to adjust the database config, and create a db_order database in config.yaml
```bash
  cd migration/
```

```bash
go run migration.go ./sql "host=localhost port=5432 user=root dbname=db_order password=fatannajuda sslmode=disable" up
```

### Down Migration
```bash
go run migration.go ./sql "host=localhost port=5432 user=root dbname=db_order password=fatannajuda sslmode=disable" down
```

### Create new SQL
```bash
go run migration.go ./sql "host=localhost port=5432 user=root dbname=db_order sslmode=disable" create add_orders_table sql
```
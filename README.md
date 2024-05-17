# pushsense
pushsense is small server that listens for changes in **postgres** database and pushes them to typesense. In order to use this tool you have to follow some rules in your db. Here i assume you are using uuid7 for primary key and you have columns **updated_at**, **deleted_at** and **is_deleted**. For every collection you provide in conf.json program will call sql query every second and new/updated/deleted rows will be replicated in typesense

### create triggers 

trigger function
```
CREATE OR REPLACE FUNCTION update_log()
  RETURNS trigger AS
$$
BEGIN
NEW.updated_at = now()
IF NEW.is_deleted AND NOT OLD.is_deleted
THEN 
	NEW.deleted_at = NEW.updated_at
END IF

RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';
```
define trigger

```
CREATE TRIGGER user_trigger
  BEFORE UPDATE OR INSERT
  ON user
  FOR EACH ROW
  EXECUTE PROCEDURE update_log();
```
### conf.json
db and typesense objects are connection settings. Important part is that every connection has a **schema** that will be used as **api.CollectionSchema** from [typesense/typesense-go](https://github.com/typesense/typesense-go)
```
{
  "db": {
    "url": "user=admin password=xuuoH8FXDSTQkWFA7QRg host=localhost port=5432 dbname=dev sslmode=disable"
  },
  "typesense": {
    "url": "http://localhost:8108",
    "apiKey": "xyz",
    "connectionTimeout": 5,
    "circuitBreakerMaxRequest": 50,
    "circuitBreakerMaxInterval": 120,
    "circuitBreakerMaxTimeout": 60
  },
  "collections": {
    "users": {
      "tableName": "users",
      "sql": "select * from users where updated_at >= $1 LIMIT 100",
      "schema" : {
        "name": "users",  
        "fields": [
          {"name": "is_active", "type": "bool"},
          {"name": ".*", "type": "auto" }
        ]
      }
    }
  }
}
```
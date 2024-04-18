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
# pushsense
pushsense is small server that listens for changes in postgres database and pushes them to typesense

# create triggers 

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
# Optimistic concurrency control for gorm


Optimistic concurrency control:
wiki: https://en.wikipedia.org/wiki/Optimistic_concurrency_control

Gorm
link: https://github.com/go-gorm/gorm


Install

```
go get github.com/goxiaoy/go-concurrency
```

Add `version` to your entity

``` go
type TestEntity struct {
    ...
	concurrency.Version
}
```

create

``` go
	e := TestEntity{
		ID:   1,
		Name: "1",
	}
	err := DB.Create(&e).Error
```

This field will be auto set before create, or you can set it manually to prevent reflection

update

``` go
	err = concurrency.ConcurrentUpdate(DB.Model(&ec), "name", "3").Error
    // To check concurrency error
    // assert.ErrorIs(t, err, ErrConcurrent)
	
```

or use gorm update
``` go
	affected := DB.Model(&ec).Update("name", "3").RowsAffected
    // check affected == 0
```

The generate sql will be 
``` sql
UPDATE `test_entities` SET `name`="3",`version`="be2e5998-c809-482f-9c9c-64c709ba6ea3" WHERE `test_entities`.`version` = "92658491-bba8-4eba-84ce-a6ea72dcfa4a" AND `id` = 1
```

Note：
Do not use `Save` method, it will automatically create entities if not found





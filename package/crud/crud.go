package crud

import (
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
)

func Create(db orm.DB, v interface{}) error {
	_, err := db.Model(v).Insert()

	if err != nil {
		return errors.Wrap(err, "error inserting record")
	}

	return nil
}

func Delete(db orm.DB, v interface{}) error {
	_, err := db.Model(v).Delete()

	if err != nil {
		return errors.Wrap(err, "error deleting record")
	}

	return nil
}

func Find(db orm.DB, v interface{}) error {
	if err := db.Model(v).Select(v); err != nil {
		return errors.Wrap(err, "error finding records")
	}

	return nil
}

func FindAll(db orm.DB, v interface{}) error {
	if err := db.Model(v).Select(); err != nil {
		return errors.Wrap(err, "error finding all records")
	}

	return nil
}

func Update(db orm.DB, v interface{}, columns ...string) error {
	if _, err := db.Model(v).Column(columns...).Update(); err != nil {
		return errors.Wrap(err, "error updating record")
	}

	return nil
}
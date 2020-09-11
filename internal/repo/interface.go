package repo

type Getter interface {
	GetMany(dest interface{}, query string, args ...interface{}) (ok bool, err error)
	GetOne(dest interface{}, query string, args ...interface{}) (ok bool, err error)
}

type Deleter interface {
	Delete(query string, args ...interface{}) (rowsAffected int64, err error)
}

type Inserter interface {
	Insert(query string, args ...interface{}) (lastInsertID int, err error)
	InsertNamed(namedQuery string, arg interface{}) (lastInsertID int, err error)
}

type Updater interface {
	Update(query string, args ...interface{}) (rowsAffected int64, err error)
	UpdateNamed(query string, arg interface{}) (rowsAffected int64, err error)
}

type Closer interface {
	Close() error
}

type DB interface {
	Closer
	Deleter
	Getter
	Inserter
	Updater
}

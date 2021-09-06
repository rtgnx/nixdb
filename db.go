package nixdb

import (
	"os"
	"path"
)

type Database struct {
	MinUID  uint
	MinGID  uint
	Users   Passwd
	Groups  Groups
	BaseDir string
}

func NewDB(baseDir string, minUID, minGID uint) Database {
	return Database{
		MinUID:  minUID,
		MinGID:  minGID,
		Users:   make(Passwd, 0),
		BaseDir: baseDir,
	}
}

func (db *Database) Update(dbName string) error {
	switch dbName {
	case "passwd":
		return db.updatePasswd()
	case "group":
		return db.updateGroup()
	}

	return nil
}

func (db *Database) updatePasswd() error {
	passwdPath := path.Join(db.BaseDir, "passwd")
	fd, err := os.Open(passwdPath)

	if err != nil {
		return err
	}

	db.Users.Read(fd)

	for i, user := range db.Users {
		if user.UID <= db.MinUID {
			db.Users = append(db.Users[:i], db.Users[i+1:]...)
		}
	}

	return nil
}

func (db *Database) updateGroup() error {
	groupPath := path.Join(db.BaseDir, "group")
	fd, err := os.Open(groupPath)

	if err != nil {
		return err
	}

	db.Groups.Read(fd)

	for i, group := range db.Groups {
		if group.GID <= db.MinGID {
			db.Groups = append(db.Groups[:i], db.Groups[i+1:]...)
		}
	}

	return nil
}

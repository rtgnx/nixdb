package nixdb

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type GroupEntry struct {
	Name     string
	Password string
	GID      uint
	Users    []string
}

type Groups []GroupEntry

func (g GroupEntry) Members(passwd Passwd) Passwd {
	userList := make(Passwd, 0)

	for _, u := range g.Users {
		for _, match := range passwd {
			if strings.Compare(match.Name, u) == 0 {
				userList = append(userList, match)
			}
		}
	}

	return userList
}

func (g GroupEntry) Encode() []byte {
	return []byte(fmt.Sprintf(
		"%s:%s:%d:%s", g.Name, g.Password, g.GID, strings.Join(g.Users, ","),
	))
}

func (g *GroupEntry) Decode(b []byte) error {
	items := strings.Split(string(b), ":")

	if len(items) != 4 {
		return fmt.Errorf("invalid group string, expected 4 values got: %d", len(items))
	}

	g.Name, g.Password, g.Users = items[0], items[1], strings.Split(items[3], ",")

	gid, err := strconv.ParseUint(items[2], 10, 32)

	g.GID = uint(gid)

	return err
}

func (g *Groups) Read(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	for _, entry := range strings.Split(string(buf), "\n") {
		v := new(GroupEntry)
		if err := v.Decode([]byte(entry)); err != nil {
			return err
		}

		*g = append(*g, *v)
	}

	return nil
}

func (g Groups) Write(w io.Writer) error {
	for _, entry := range g {
		if _, err := w.Write(entry.Encode()); err != nil {
			return err
		}
	}

	return nil
}

func (g Groups) FindByName(name string) (GroupEntry, bool) {
	for _, group := range g {
		if strings.Compare(group.Name, name) == 0 {
			return group, true
		}
	}

	return GroupEntry{}, false
}

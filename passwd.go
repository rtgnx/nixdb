package nixdb

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type PasswdEntry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	UID      uint   `json:"uid"`
	GID      uint   `json:"gid"`
	Fullname string `json:"fullname"`
	Home     string `json:"home"`
	Shell    string `json:"shell"`
}

type Passwd []PasswdEntry

func (p *PasswdEntry) Decode(b []byte) error {

	items := strings.Split(string(b), ":")

	if len(items) < 7 {
		return fmt.Errorf("invalid passwd entry string, expected 7 values colon delimited")
	}

	p.Name, p.Password, p.Fullname, p.Home, p.Shell = items[0], items[1], items[4], items[5], items[6]

	var uid, gid uint64
	var uidErr, gidErr error

	uid, uidErr = strconv.ParseUint(items[2], 10, 32)

	gid, gidErr = strconv.ParseUint(items[2], 10, 32)

	if uidErr != nil {
		return uidErr
	} else if gidErr != nil {
		return gidErr
	}

	p.UID, p.GID = uint(uid), uint(gid)

	return nil
}

func (p PasswdEntry) Encode() []byte {
	return []byte(
		fmt.Sprintf(
			"%s:%s:%d:%d:%s:%s:%s",
			p.Name, p.Password, p.UID, p.GID, p.Fullname, p.Home, p.Shell,
		),
	)
}

func (p *Passwd) Read(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	for _, entry := range strings.Split(string(buf), "\n") {
		v := new(PasswdEntry)
		if err := v.Decode([]byte(entry)); err != nil {
			return err
		}

		*p = append(*p, *v)
	}

	return nil
}

func (p Passwd) Write(w io.Writer) error {
	for _, entry := range p {
		if _, err := w.Write(entry.Encode()); err != nil {
			return err
		}
	}

	return nil
}

func (p *Passwd) Deduplicate() {
	for i, u := range *p {
		str := string(u.Encode())

		for j, d := range *p {
			if j != i {
				str2 := string(d.Encode())

				if strings.Compare(str, str2) == 0 {
					ptr := []PasswdEntry(*p)
					ptr = append(ptr[:j], ptr[j+1:]...)
					*p = ptr
				}
			}
		}
	}
}

func (p Passwd) FindByName(name string) (PasswdEntry, bool) {
	for _, entry := range p {
		if strings.Compare(entry.Name, name) == 0 {
			return entry, true
		}
	}

	return PasswdEntry{}, false
}

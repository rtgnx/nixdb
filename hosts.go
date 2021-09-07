package nixdb

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type HostEntry struct {
	IPAddress string
	FQDN      string
	Aliases   []string
}

type Hosts []HostEntry

func (h HostEntry) Encode() []byte {
	return []byte(fmt.Sprintf("%s\t%s\t%s", h.IPAddress, h.FQDN, strings.Join(h.Aliases, " ")))
}

func (h *HostEntry) Decode(b []byte) error {
	fields := strings.Fields(string(b))

	if len(fields) < 2 {
		return fmt.Errorf("invalid hosts file entry")
	}

	h.IPAddress, h.FQDN = fields[0], fields[1]

	if len(fields) > 2 {
		h.Aliases = fields[3:]
	}

	return nil
}

func (h Hosts) Write(w io.Writer) error {
	for _, entry := range h {
		if _, err := w.Write(entry.Encode()); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hosts) Read(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	for _, entry := range strings.Split(string(buf), "\n") {
		v := new(HostEntry)
		if err := v.Decode([]byte(entry)); err != nil {
			return err
		}

		*h = append(*h, *v)
	}

	return nil
}

func (h Hosts) Lookup(phrase string) (HostEntry, bool) {
	for _, host := range h {
		switch 0 {
		case strings.Compare(phrase, host.IPAddress):
			return host, true
		case strings.Compare(phrase, host.FQDN):
			return host, true
		default:
			for _, alias := range host.Aliases {
				if strings.Compare(alias, phrase) == 0 {
					return host, true
				}
			}
		}
	}
	return HostEntry{}, false
}

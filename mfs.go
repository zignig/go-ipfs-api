package shell

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// mfs object for 0.4
type MFS struct {
	base       string
	currentRef string
	s          *Shell
}

func (s *Shell) NewMfs(base string) (m *MFS) {
	m = &MFS{}
	if strings.HasSuffix(base, "/") == false {
		base = base + "/"
	}
	if strings.HasPrefix(base, "/") == false {
		base = "/" + base
	}
	m.base = base
	m.s = s
	return
}

func (m *MFS) Mkdir(path string) (e error) {
	req := NewRequest(m.s.url, "files/mkdir", m.base+path)
	req.Opts["p"] = "true"
	resp, err := req.Send(m.s.httpcli)
	if err != nil {
		return err
	}
	defer resp.Close()
	return
}

type List struct {
	Entries []LsLink
}

func (m *MFS) Ls(path string) (list []string, e error) {
	req := NewRequest(m.s.url, "files/ls", m.base+path)
	req.Opts["l"] = "true"
	resp, err := req.Send(m.s.httpcli)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	defer resp.Close()
	var out List
	err = json.NewDecoder(resp.Output).Decode(&out)
	b, _ := json.MarshalIndent(out, " ", " ")
	fmt.Println(string(b))
	if err != nil {
		fmt.Println("decode error ", err)
	}
	list = make([]string, 0)
	for _, j := range out.Entries {
		list = append(list, j.Name)
	}
	fmt.Println(list)
	return
}

func (m *MFS) Stat(path string) (e error) {
	resp, err := NewRequest(m.s.url, "files/stat", m.base+path).Send(m.s.httpcli)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}
	defer resp.Close()
	data, err := ioutil.ReadAll(resp.Output)
	fmt.Println("print data")
	fmt.Println(string(data))
	return
}

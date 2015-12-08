// package shell implements a remote API interface for a running ipfs daemon
package shell

import (
	"encoding/json"
	"fmt"
)

type Pingee struct {
	Success bool
	Time    int
	Text    string
}

func (s *Shell) Ping(id string) (active bool, e error) {
	req := NewRequest(s.url, "ping", id)
	req.Opts["n"] = "2"
	resp, err := req.Send(s.httpcli)
	decoder := json.NewDecoder(resp.Output)
	out := new(Pingee)
	var counter = 0
	for {
		err = decoder.Decode(out)
		if out.Success == true {
			return true, err
		} else {
			counter = counter + 1
		}
		if counter > 2 {
			return false, err
		}
	}
	return false, err
}

type Addr []string

type Addrs struct {
	Addrs map[string]Addr
}

func (s *Shell) Swarm() (out *Addrs, e error) {
	resp, err := NewRequest(s.url, "swarm/addrs", "").Send(s.httpcli)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	decoder := json.NewDecoder(resp.Output)
	out = new(Addrs)
	err = decoder.Decode(out)
	if err != nil {
		return nil, err
	}
	return out, err

}

type Resp struct {
	Addrs []string
	ID    string
}

type Provs struct {
	Extra     string
	ID        string
	Responses []Resp
	Type      int
}

func (s *Shell) FindProvs(hash string) (list []string, err error) {
	resp, err := NewRequest(s.url, "dht/findprovs", hash).Send(s.httpcli)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return
	}
	decoder := json.NewDecoder(resp.Output)
	out := new(Provs)
	for {
		err = decoder.Decode(out)
		if err != nil {
			fmt.Println(err)
			break
		}
		if out.Type == 4 {
			for _, j := range out.Responses {
				list = append(list, j.ID)
				//fmt.Println("List : ", list)
			}
		}
	}
	return list, nil
}

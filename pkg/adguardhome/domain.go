package adguardhome

import (
	"sort"

	"github.com/mr-karan/doggo/pkg/resolvers"
)

// ------------------------------------------------------------------
//
// ------------------------------------------------------------------

type HostRecord struct {
	Host   string
	IPs    []string
	oldIPs []string
}

func (h *HostRecord) IPFromAnswer(answers []resolvers.Answer) {
	if len(answers) == 0 {
		return
	}
	h.oldIPs = h.IPs
	h.IPs = make([]string, 0, len(answers))
	for _, answer := range answers {
		h.IPs = append(h.IPs, answer.Address)
	}
	h.RemoveDuplicate()
}

func (h *HostRecord) DeepCopy() *HostRecord {
	r := new(HostRecord)
	r.Host = h.Host
	if h.IPs == nil {
		return r
	}

	r.IPs = make([]string, len(h.IPs))
	copy(r.IPs, h.IPs)
	return r
}

func (h *HostRecord) RemoveDuplicate() {
	if len(h.IPs) == 0 {
		return
	}
	sort.Strings(h.IPs)
	i := 0
	for j := 1; j < len(h.IPs); j++ {
		if h.IPs[j] != h.IPs[i] {
			i++
			h.IPs[i] = h.IPs[j]
		}
	}
	h.IPs = h.IPs[:i+1]
}

// IpDiff 对比两个IP数组 OldIPs、NewIPs，这两个切片都是经过排序的。
//
// OldIPs中存在并且NewIPs中不存在的需要被删除；
// NewIPs中存在并且OldIPs中不存在的需要被添加。
func (h *HostRecord) IpDiff() (toBeDelete, toBeAdd []string) {
	i, j := 0, 0
	for i < len(h.oldIPs) && j < len(h.IPs) {
		switch {
		case h.oldIPs[i] < h.IPs[j]:
			toBeDelete = append(toBeDelete, h.oldIPs[i])
			i++
		case h.oldIPs[i] > h.IPs[j]:
			toBeAdd = append(toBeAdd, h.IPs[j])
			j++
		default:
			i++
			j++
		}
	}

	for ; i < len(h.oldIPs); i++ {
		toBeDelete = append(toBeDelete, h.oldIPs[i])
	}
	for ; j < len(h.IPs); j++ {
		toBeAdd = append(toBeAdd, h.IPs[j])
	}
	return
}

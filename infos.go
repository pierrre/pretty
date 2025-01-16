package pretty

import (
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal/write"
)

type infos struct {
	showLen  bool
	len      int
	showCap  bool
	cap      int
	showAddr bool
	addr     uintptr
}

func (i infos) write(st *State) bool {
	if !st.ShowInfos {
		return false
	}
	if !i.showLen && !i.showCap && !i.showAddr {
		return false
	}
	w := st.Writer
	write.MustString(w, "(")
	wrote := false
	if i.showLen {
		write.MustString(w, "len=")
		write.Must(strconvio.WriteInt(w, int64(i.len), 10))
		wrote = true
	}
	if i.showCap {
		if wrote {
			write.MustString(w, " ")
		}
		write.MustString(w, "cap=")
		write.Must(strconvio.WriteInt(w, int64(i.cap), 10))
		wrote = true
	}
	if i.showAddr {
		if wrote {
			write.MustString(w, " ")
		}
		write.MustString(w, "addr=")
		writeUintptr(w, i.addr)
	}
	write.MustString(w, ")")
	return true
}

func (i infos) writeWithTrailingSpace(st *State) {
	if i.write(st) {
		write.MustString(st.Writer, " ")
	}
}

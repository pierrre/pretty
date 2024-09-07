package pretty

import (
	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/pretty/internal"
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
	writeString(w, "(")
	wrote := false
	if i.showLen {
		writeString(w, "len=")
		internal.MustWrite(strconvio.WriteInt(w, int64(i.len), 10))
		wrote = true
	}
	if i.showCap {
		if wrote {
			writeString(w, " ")
		}
		writeString(w, "cap=")
		internal.MustWrite(strconvio.WriteInt(w, int64(i.cap), 10))
		wrote = true
	}
	if i.showAddr {
		if wrote {
			writeString(w, " ")
		}
		writeString(w, "addr=")
		writeUintptr(w, i.addr)
	}
	writeString(w, ")")
	return true
}

func (i infos) writeWithTrailingSpace(st *State) {
	if i.write(st) {
		writeString(st.Writer, " ")
	}
}

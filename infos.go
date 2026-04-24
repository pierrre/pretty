package pretty

import (
	"strconv"
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
	st.Writer.AppendByte('(')
	wrote := false
	if i.showLen {
		st.Writer.AppendString("len=")
		st.Writer = strconv.AppendInt(st.Writer, int64(i.len), 10)
		wrote = true
	}
	if i.showCap {
		if wrote {
			st.Writer.AppendByte(' ')
		}
		st.Writer.AppendString("cap=")
		st.Writer = strconv.AppendInt(st.Writer, int64(i.cap), 10)
		wrote = true
	}
	if i.showAddr {
		if wrote {
			st.Writer.AppendByte(' ')
		}
		st.Writer.AppendString("addr=")
		writeUintptr(st, i.addr)
	}
	st.Writer.AppendByte(')')
	return true
}

func (i infos) writeWithTrailingSpace(st *State) {
	if i.write(st) {
		st.Writer.AppendByte(' ')
	}
}

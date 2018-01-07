package collect

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "K_=T_,String,Sha V_=T_,Boolean,Integer"

type K_V_Map map[K_]V_

func (mp *K_V_Map) Init() (m K_V_Map) {
	if m = *mp; m == nil {
		m = make(K_V_Map)
		*mp = m
	}
	return
}

func (mp *K_V_Map) Put(key K_, val V_) bool {
	m := mp.Init()
	if v, ok := m[key]; ok && v == val {
		return false
	}
	m[key] = val
	return true
}

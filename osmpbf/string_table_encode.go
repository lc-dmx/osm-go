package osmpbf

type StrTable struct {
	tableSize int
	idxMap    map[string]int
	table     []string
}

func NewStringTable() *StrTable {
	return &StrTable{
		idxMap: make(map[string]int, 0),
		table:  []string{""}, // Add an unused string at offset 0 which is used as a delimiter.
	}
}

func (s *StrTable) addStrToTable(str string) int {
	s.tableSize += len(str)
	s.table = append(s.table, str)

	return len(s.table) - 1
}

func (s *StrTable) getStrIndex(str string) int {
	if idx, ok := s.idxMap[str]; ok {
		return idx
	} else {
		idx = s.addStrToTable(str)
		s.idxMap[str] = idx
		return idx
	}
}

func (s *StrTable) getStrTableSize() int {
	return s.tableSize
}

func (s *StrTable) getStrTable() []string {
	return s.table
}

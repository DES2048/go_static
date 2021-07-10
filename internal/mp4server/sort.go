package mp4server

type Mp4ByModTime []FileInfo

func (m Mp4ByModTime) Len() int           { return len(m) }
func (m Mp4ByModTime) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Mp4ByModTime) Less(i, j int) bool { return m[i].ModTime().Unix() < m[j].ModTime().Unix() }

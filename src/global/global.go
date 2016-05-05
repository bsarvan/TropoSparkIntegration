package global

var GlobalData = make(map[string]GlobalDS)

type GlobalDS struct {
    Mobile []string
    Sparkid []string
}    

# go-pipes


'
// Example Usage

type LineReader struct {
	lines []string
}

func (l LineReader) ProcessData(
	msg firestorm.Message,
	outChan chan firestorm.Message,
	errChan chan error,
) {

	for _, line := range l.lines {
		var msg firestorm.Message
		msg = []byte(line)
		outChan <- msg
	}

}

type LineParser struct {
}

func (l LineParser) ProcessData(
	msg firestorm.Message,
	outChan chan firestorm.Message,
	errChan chan error,
) {
	newMsg := strings.ToUpper(string(msg))
	outChan <- []byte(newMsg)

}

type LineOutput struct {
}

func (l LineOutput) ProcessData(
	msg firestorm.Message,
	outChan chan firestorm.Message,
	errChan chan error,
) {

	log.Println("load:", string(msg))

}

func main() {
	extract := LineReader{}
	extract.lines = []string{"save", "our", "souls"}

	transform := LineParser{}

	load := LineOutput{}

	pl := firestorm.NewPipeline(extract, transform, load)
	pl.Run()
}

'

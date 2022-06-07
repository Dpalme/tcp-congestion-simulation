package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"tcp-congestion/pkg/client"
	"tcp-congestion/pkg/connection"
	"tcp-congestion/pkg/server"
)

type Config struct {
	v int
	o string
	r int
}

type LogFile struct {
	file *os.File
}

func TwoWayShake(c *client.Client, s *server.Server) {
	con := connection.New(c.IP)
	c.Connect(con)
	s.AddConnection(con)
}

func NewLogFile(filename string) *LogFile {
	dirname, _ := path.Split(filename)

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.MkdirAll(dirname, 0755)
	}
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	return &LogFile{file: f}
}

func (l *LogFile) Write(p []byte) (n int, err error) {
	return l.file.Write(p)
}

func (l *LogFile) WriteString(s string) (n int, err error) {
	return l.file.Write([]byte(s))
}

func (l *LogFile) Writef(format string, args ...interface{}) (n int, err error) {
	return l.file.WriteString(fmt.Sprintf(format, args...))
}

func (l *LogFile) Close() error {
	return l.file.Close()
}

func CreateLogger(tag string) *log.Logger {
	return log.New(log.Writer(), tag, log.LstdFlags)
}

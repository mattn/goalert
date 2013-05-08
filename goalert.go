package main

import (
	"bytes"
	"flag"
	"github.com/mattn/go-gntp"
	"io"
	"os"
	"os/exec"
)

var server = flag.String("s", "127.0.0.1:23053", "GNTP server")
var action = flag.String("a", "", "Click action")

func main() {
	flag.Parse()

	var buf bytes.Buffer
	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
	err := cmd.Run()
	growl := gntp.NewClient()
	growl.Server = *server
	growl.AppName = "goalert"
	growl.Register([]gntp.Notification {
		gntp.Notification{
			Event: "success",
			Enabled: false,
		}, gntp.Notification{
			Event: "failed",
			Enabled: true,
		},
	})
	event := "success"
	title, _ := os.Getwd()
	text := string(buf.Bytes())
	if err != nil {
		event = "failed"
	}
	callback := *action
	if callback == "" {
		callback = title
	}
	growl.Notify(&gntp.Message{
		Event: event,
		Title: title,
		Text: text,
		Callback: callback,
	})
}

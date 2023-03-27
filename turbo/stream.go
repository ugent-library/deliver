package turbo

import (
	"bytes"
	"net/http"
	"strings"
	"sync"
)

type StreamAction string

const (
	AppendAction  StreamAction = "append"
	PrependAction StreamAction = "prepend"
	ReplaceAction StreamAction = "replace"
	UpdateAction  StreamAction = "update"
	RemoveAction  StreamAction = "remove"
	BeforeAction  StreamAction = "before"
	AfterAction   StreamAction = "after"

	StreamContentType = "text/vnd.turbo-stream.html; charset=utf-8"
)

var bufPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

func StreamRequest(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Accept"), "text/vnd.turbo-stream.html")
}

func StreamSourceTag(src string) string {
	return `<turbo-stream-source src="` + src + `"></turbo-stream-source>`
}

func Encode(streams []StreamMessage) ([]byte, error) {
	b := bufPool.Get().(*bytes.Buffer)
	defer func() {
		b.Reset()
		bufPool.Put(b)
	}()

	for _, s := range streams {
		b.WriteString(`<turbo-stream action="`)
		b.WriteString(string(s.Action))
		b.WriteString(`" `)
		if s.Target != "" {
			b.WriteString(`target="`)
			b.WriteString(s.Target)
		} else {
			b.WriteString(`targets="`)
			b.WriteString(s.TargetSelector)
		}
		b.WriteString(`">`)
		if s.Action != RemoveAction {
			b.WriteString(`<template>`)
			b.WriteString(s.Template)
			b.WriteString(`</template>`)
		}
		b.WriteString(`</turbo-stream>`)
	}
	return b.Bytes(), nil
}

type StreamMessage struct {
	Action         StreamAction
	Target         string
	TargetSelector string
	Template       string
}

func Append(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   AppendAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func AppendMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         AppendAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Prepend(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   PrependAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func PrependMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         PrependAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Replace(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   ReplaceAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func ReplaceMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         ReplaceAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Update(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   UpdateAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func UpdateMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         UpdateAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Remove(target string) StreamMessage {
	return StreamMessage{
		Action: RemoveAction,
		Target: target,
	}
}

func RemoveMatch(target string) StreamMessage {
	return StreamMessage{
		Action:         RemoveAction,
		TargetSelector: target,
	}
}

func Before(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   BeforeAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func BeforeMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         BeforeAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func After(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:   AfterAction,
		Target:   target,
		Template: strings.Join(tmpls, ""),
	}
}

func AfterMatch(target string, tmpls ...string) StreamMessage {
	return StreamMessage{
		Action:         AfterAction,
		TargetSelector: target,
		Template:       strings.Join(tmpls, ""),
	}
}

func Render(w http.ResponseWriter, r *http.Request, code int, streams ...StreamMessage) error {
	if hdr := w.Header(); hdr.Get("Content-Type") == "" {
		hdr.Set("Content-Type", StreamContentType)
	}
	w.WriteHeader(code)
	b, err := Encode(streams)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

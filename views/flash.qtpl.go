// Code generated by qtc from "flash.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/flash.qtpl:1
package views

//line views/flash.qtpl:1
import "github.com/ugent-library/deliver/ctx"

//line views/flash.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/flash.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/flash.qtpl:3
func StreamAddFlash(qw422016 *qt422016.Writer, f ctx.Flash) {
//line views/flash.qtpl:3
	qw422016.N().S(`
    <div hx-swap-oob="beforeend:#flash-messages">
        `)
//line views/flash.qtpl:5
	StreamFlash(qw422016, f)
//line views/flash.qtpl:5
	qw422016.N().S(`
    </div>
`)
//line views/flash.qtpl:7
}

//line views/flash.qtpl:7
func WriteAddFlash(qq422016 qtio422016.Writer, f ctx.Flash) {
//line views/flash.qtpl:7
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/flash.qtpl:7
	StreamAddFlash(qw422016, f)
//line views/flash.qtpl:7
	qt422016.ReleaseWriter(qw422016)
//line views/flash.qtpl:7
}

//line views/flash.qtpl:7
func AddFlash(f ctx.Flash) string {
//line views/flash.qtpl:7
	qb422016 := qt422016.AcquireByteBuffer()
//line views/flash.qtpl:7
	WriteAddFlash(qb422016, f)
//line views/flash.qtpl:7
	qs422016 := string(qb422016.B)
//line views/flash.qtpl:7
	qt422016.ReleaseByteBuffer(qb422016)
//line views/flash.qtpl:7
	return qs422016
//line views/flash.qtpl:7
}

//line views/flash.qtpl:9
func StreamFlash(qw422016 *qt422016.Writer, f ctx.Flash) {
//line views/flash.qtpl:9
	qw422016.N().S(`
    <div class="toast" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="toast-body">
            `)
//line views/flash.qtpl:12
	switch f.Type {
//line views/flash.qtpl:13
	case "success":
//line views/flash.qtpl:13
		qw422016.N().S(`
            <i class="if if--success if-check-circle-fill"></i>
            `)
//line views/flash.qtpl:15
	case "info":
//line views/flash.qtpl:15
		qw422016.N().S(`
            <i class="if if--primary if-info-circle-filled"></i>
            `)
//line views/flash.qtpl:17
	case "warning":
//line views/flash.qtpl:17
		qw422016.N().S(`
            <i class="if if--warning if-alert-fill"></i>
            `)
//line views/flash.qtpl:19
	case "error":
//line views/flash.qtpl:19
		qw422016.N().S(`
            <i class="if if--error if-error-circle-fill"></i>
            `)
//line views/flash.qtpl:21
	}
//line views/flash.qtpl:21
	qw422016.N().S(`
            <div class="toast-content">
                `)
//line views/flash.qtpl:23
	if f.Title != "" {
//line views/flash.qtpl:23
		qw422016.N().S(`
                <h3 class="alert-title">`)
//line views/flash.qtpl:24
		qw422016.E().S(f.Title)
//line views/flash.qtpl:24
		qw422016.N().S(`</h3>
                `)
//line views/flash.qtpl:25
	}
//line views/flash.qtpl:25
	qw422016.N().S(`
                `)
//line views/flash.qtpl:26
	qw422016.E().S(f.Body)
//line views/flash.qtpl:26
	qw422016.N().S(`
            </div>
            <button class="btn-close" type="button" aria-label="Close"
                data-bs-dismiss="toast"
                `)
//line views/flash.qtpl:30
	if f.DismissAfter != 0 {
//line views/flash.qtpl:30
		qw422016.N().S(`
                data-delay="`)
//line views/flash.qtpl:31
		qw422016.N().DL(f.DismissAfter.Milliseconds())
//line views/flash.qtpl:31
		qw422016.N().S(`"
                `)
//line views/flash.qtpl:32
	} else {
//line views/flash.qtpl:32
		qw422016.N().S(`
                data-autohide="false"
                `)
//line views/flash.qtpl:34
	}
//line views/flash.qtpl:34
	qw422016.N().S(`
            >
                <i class="if if-close"></i>
            </button>
        </div>
    </div>
`)
//line views/flash.qtpl:40
}

//line views/flash.qtpl:40
func WriteFlash(qq422016 qtio422016.Writer, f ctx.Flash) {
//line views/flash.qtpl:40
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/flash.qtpl:40
	StreamFlash(qw422016, f)
//line views/flash.qtpl:40
	qt422016.ReleaseWriter(qw422016)
//line views/flash.qtpl:40
}

//line views/flash.qtpl:40
func Flash(f ctx.Flash) string {
//line views/flash.qtpl:40
	qb422016 := qt422016.AcquireByteBuffer()
//line views/flash.qtpl:40
	WriteFlash(qb422016, f)
//line views/flash.qtpl:40
	qs422016 := string(qb422016.B)
//line views/flash.qtpl:40
	qt422016.ReleaseByteBuffer(qb422016)
//line views/flash.qtpl:40
	return qs422016
//line views/flash.qtpl:40
}

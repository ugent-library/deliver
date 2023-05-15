// Code generated by qtc from "edit_space.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/edit_space.qtpl:1
package views

//line views/edit_space.qtpl:1
import "strings"

//line views/edit_space.qtpl:2
import "github.com/ugent-library/deliver/ctx"

//line views/edit_space.qtpl:3
import "github.com/ugent-library/deliver/models"

//line views/edit_space.qtpl:4
import "github.com/ugent-library/deliver/validate"

//line views/edit_space.qtpl:6
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/edit_space.qtpl:6
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/edit_space.qtpl:7
type EditSpace struct {
	Space            *models.Space
	ValidationErrors *validate.Errors
}

//line views/edit_space.qtpl:13
func (v *EditSpace) StreamTitle(qw422016 *qt422016.Writer, c *ctx.Ctx) {
//line views/edit_space.qtpl:13
	qw422016.N().S(`Edit `)
//line views/edit_space.qtpl:13
	qw422016.E().S(v.Space.Name)
//line views/edit_space.qtpl:13
}

//line views/edit_space.qtpl:13
func (v *EditSpace) WriteTitle(qq422016 qtio422016.Writer, c *ctx.Ctx) {
//line views/edit_space.qtpl:13
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/edit_space.qtpl:13
	v.StreamTitle(qw422016, c)
//line views/edit_space.qtpl:13
	qt422016.ReleaseWriter(qw422016)
//line views/edit_space.qtpl:13
}

//line views/edit_space.qtpl:13
func (v *EditSpace) Title(c *ctx.Ctx) string {
//line views/edit_space.qtpl:13
	qb422016 := qt422016.AcquireByteBuffer()
//line views/edit_space.qtpl:13
	v.WriteTitle(qb422016, c)
//line views/edit_space.qtpl:13
	qs422016 := string(qb422016.B)
//line views/edit_space.qtpl:13
	qt422016.ReleaseByteBuffer(qb422016)
//line views/edit_space.qtpl:13
	return qs422016
//line views/edit_space.qtpl:13
}

//line views/edit_space.qtpl:15
func (v *EditSpace) StreamContent(qw422016 *qt422016.Writer, c *ctx.Ctx) {
//line views/edit_space.qtpl:15
	qw422016.N().S(`
    <div class="w-100 u-scroll-wrapper">
        <div class="bg-white">
            <div class="bc-navbar bc-navbar--xlarge bc-navbar--white bc-navbar--bordered-bottom">
                <div class="bc-toolbar">
                    <div class="bc-toolbar-left">
                        <div class="bc-toolbar-item">
                            <h4 class="bc-toolbar-title">Edit `)
//line views/edit_space.qtpl:22
	qw422016.E().S(v.Space.Name)
//line views/edit_space.qtpl:22
	qw422016.N().S(`</h4>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="u-scroll-wrapper__body p-6">
            <div class="card w-100 mb-6">
                <div class="card-header">
                    <div class="bc-toolbar">
                        <div class="bc-toolbar-left">
                            <div class="bc-toolbar-item">Edit `)
//line views/edit_space.qtpl:33
	qw422016.E().S(v.Space.Name)
//line views/edit_space.qtpl:33
	qw422016.N().S(`</div>
                        </div>
                        <div class="bc-toolbar-right">
                            <div class="bc-toolbar-item">
                                <a class="btn btn-link btn-link-muted" href="`)
//line views/edit_space.qtpl:37
	qw422016.E().S(c.PathTo("space", "spaceName", v.Space.Name).String())
//line views/edit_space.qtpl:37
	qw422016.N().S(`">
                                    <i class="if if-close"></i>
                                    <span class="btn-text">Discard changes</span>
                                </a>
                            </div>
                            <div class="bc-toolbar-item">
                                <button class="btn btn-primary" data-submit-target="#update-space">
                                    <i class="if if-check"></i>
                                    <span class="btn-text">Save changes</span>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card-body">
                    <form action="`)
//line views/edit_space.qtpl:52
	qw422016.E().S(c.PathTo("updateSpace", "spaceName", v.Space.Name).String())
//line views/edit_space.qtpl:52
	qw422016.N().S(`" method="POST" id="update-space">
                        <input type="hidden" name="_method" value="PUT">
                        `)
//line views/edit_space.qtpl:54
	qw422016.N().S(c.CSRFTag)
//line views/edit_space.qtpl:54
	qw422016.N().S(`
                        <div class="form-group">
                            <div class="form-row form-group">
                                <label class="col-lg-3 col-xl-2 col-form-label" for="space-admins">Space admins</label>
                                <div class="col-lg-5 col-xl-4">
                                    <input class="form-control" type="text" value="`)
//line views/edit_space.qtpl:59
	qw422016.E().S(strings.Join(v.Space.Admins, ","))
//line views/edit_space.qtpl:59
	qw422016.N().S(`" id="space-admins" name="admins">
                                    <p class="small form-text text-muted">Separate usernames with a comma.</p>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
`)
//line views/edit_space.qtpl:69
}

//line views/edit_space.qtpl:69
func (v *EditSpace) WriteContent(qq422016 qtio422016.Writer, c *ctx.Ctx) {
//line views/edit_space.qtpl:69
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/edit_space.qtpl:69
	v.StreamContent(qw422016, c)
//line views/edit_space.qtpl:69
	qt422016.ReleaseWriter(qw422016)
//line views/edit_space.qtpl:69
}

//line views/edit_space.qtpl:69
func (v *EditSpace) Content(c *ctx.Ctx) string {
//line views/edit_space.qtpl:69
	qb422016 := qt422016.AcquireByteBuffer()
//line views/edit_space.qtpl:69
	v.WriteContent(qb422016, c)
//line views/edit_space.qtpl:69
	qs422016 := string(qb422016.B)
//line views/edit_space.qtpl:69
	qt422016.ReleaseByteBuffer(qb422016)
//line views/edit_space.qtpl:69
	return qs422016
//line views/edit_space.qtpl:69
}

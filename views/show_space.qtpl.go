// Code generated by qtc from "show_space.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/show_space.qtpl:1
package views

//line views/show_space.qtpl:1
import "time"

//line views/show_space.qtpl:2
import "github.com/ugent-library/friendly"

//line views/show_space.qtpl:3
import "github.com/ugent-library/deliver/ctx"

//line views/show_space.qtpl:4
import "github.com/ugent-library/deliver/models"

//line views/show_space.qtpl:5
import "github.com/ugent-library/deliver/validate"

//line views/show_space.qtpl:7
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/show_space.qtpl:7
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/show_space.qtpl:8
type ShowSpace struct {
	Space            *models.Space
	UserSpaces       []*models.Space
	Folder           *models.Folder
	ValidationErrors *validate.Errors
}

//line views/show_space.qtpl:16
func (v *ShowSpace) StreamTitle(qw422016 *qt422016.Writer, c *ctx.Ctx) {
//line views/show_space.qtpl:16
	qw422016.E().S(v.Space.Name)
//line views/show_space.qtpl:16
}

//line views/show_space.qtpl:16
func (v *ShowSpace) WriteTitle(qq422016 qtio422016.Writer, c *ctx.Ctx) {
//line views/show_space.qtpl:16
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/show_space.qtpl:16
	v.StreamTitle(qw422016, c)
//line views/show_space.qtpl:16
	qt422016.ReleaseWriter(qw422016)
//line views/show_space.qtpl:16
}

//line views/show_space.qtpl:16
func (v *ShowSpace) Title(c *ctx.Ctx) string {
//line views/show_space.qtpl:16
	qb422016 := qt422016.AcquireByteBuffer()
//line views/show_space.qtpl:16
	v.WriteTitle(qb422016, c)
//line views/show_space.qtpl:16
	qs422016 := string(qb422016.B)
//line views/show_space.qtpl:16
	qt422016.ReleaseByteBuffer(qb422016)
//line views/show_space.qtpl:16
	return qs422016
//line views/show_space.qtpl:16
}

//line views/show_space.qtpl:18
func (v *ShowSpace) StreamContent(qw422016 *qt422016.Writer, c *ctx.Ctx) {
//line views/show_space.qtpl:18
	qw422016.N().S(`
    <div class="c-sub-sidebar c-sidebar--bordered"
        hx-ext="ws"
        ws-connect="`)
//line views/show_space.qtpl:21
	qw422016.E().S(c.WebSocketPath("space." + v.Space.ID))
//line views/show_space.qtpl:21
	qw422016.N().S(`"
    >
        <div class="bc-navbar bc-navbar--xlarge bc-navbar--bordered-bottom">
            <div class="bc-toolbar">
                <div class="bc-toolbar-left">
                    <div class="bc-toolbar-item">
                        <h1 class="bc-toolbar-title">Your deliver spaces</h1>
                    </div>
                </div>
            </div>
        </div>
        <div class="c-sub-sidebar__menu my-6">
            <nav>
                <ul class="c-sub-sidebar-menu">
                    `)
//line views/show_space.qtpl:35
	for _, s := range v.UserSpaces {
//line views/show_space.qtpl:35
		qw422016.N().S(`
                    <li class="c-sub-sidebar__item`)
//line views/show_space.qtpl:36
		if s.ID == v.Space.ID {
//line views/show_space.qtpl:36
			qw422016.N().S(` c-sub-sidebar__item--active`)
//line views/show_space.qtpl:36
		}
//line views/show_space.qtpl:36
		qw422016.N().S(`">
                        <a href="`)
//line views/show_space.qtpl:37
		qw422016.E().S(c.PathTo("space", "spaceName", s.Name).String())
//line views/show_space.qtpl:37
		qw422016.N().S(`">
                            <span class="c-sidebar__label">`)
//line views/show_space.qtpl:38
		qw422016.E().S(s.Name)
//line views/show_space.qtpl:38
		qw422016.N().S(`</span>
                        </a>
                    </li>
                    `)
//line views/show_space.qtpl:41
	}
//line views/show_space.qtpl:41
	qw422016.N().S(`
                    `)
//line views/show_space.qtpl:42
	if c.Permissions.IsAdmin(c.User) {
//line views/show_space.qtpl:42
		qw422016.N().S(`
                    <li class="c-sub-sidebar__item">
                        <a href="`)
//line views/show_space.qtpl:44
		qw422016.E().S(c.PathTo("newSpace").String())
//line views/show_space.qtpl:44
		qw422016.N().S(`">
                            <span class="c-sidebar__label">
                                <i class="if if-add"></i>
                                Make a new space
                            </span>
                        </a>
                    </li>
                    `)
//line views/show_space.qtpl:51
	}
//line views/show_space.qtpl:51
	qw422016.N().S(`
                </ul>
            </nav>
        </div>
    </div>

    <div class="w-100 u-scroll-wrapper">
        <div class="bg-white">
            <div class="bc-navbar bc-navbar--xlarge bc-navbar--white bc-navbar--bordered-bottom">
                <div class="bc-toolbar">
                    <div class="bc-toolbar-left">
                        <div class="bc-toolbar-item">
                            <h1 class="bc-toolbar-title">`)
//line views/show_space.qtpl:63
	qw422016.E().S(v.Space.Name)
//line views/show_space.qtpl:63
	qw422016.N().S(` folders</h1>
                        </div>
                    </div>
                    `)
//line views/show_space.qtpl:66
	if c.Permissions.IsAdmin(c.User) {
//line views/show_space.qtpl:66
		qw422016.N().S(`
                    <div class="bc-toolbar-right">
                        <div class="bc-toolbar-item">
                            <a class="btn btn-link btn-link-muted" href="`)
//line views/show_space.qtpl:69
		qw422016.E().S(c.PathTo("editSpace", "spaceName", v.Space.Name).String())
//line views/show_space.qtpl:69
		qw422016.N().S(`">
                                <i class="if if-edit"></i>
                                <span class="btn-text">Edit space</span>
                            </a>
                        </div>
                    </div>
                    `)
//line views/show_space.qtpl:75
	}
//line views/show_space.qtpl:75
	qw422016.N().S(`
                </div>
            </div>
        </div>
        <div class="u-scroll-wrapper__body p-6">
            <div class="card w-100 mb-6">
                <div class="card-header">
                    <div class="bc-toolbar">
                        <div class="bc-toolbar-left">
                            <div class="bc-toolbar-item">
                                <h2>Make a folder</h2>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card-body">
                    <form action="`)
//line views/show_space.qtpl:91
	qw422016.E().S(c.PathTo("createFolder", "spaceName", v.Space.Name).String())
//line views/show_space.qtpl:91
	qw422016.N().S(`" method="POST">
                        `)
//line views/show_space.qtpl:92
	qw422016.N().S(c.CSRFTag)
//line views/show_space.qtpl:92
	qw422016.N().S(`
                        <div class="mb-6`)
//line views/show_space.qtpl:93
	if v.ValidationErrors.Get("name") != nil {
//line views/show_space.qtpl:93
		qw422016.N().S(` is-invalid`)
//line views/show_space.qtpl:93
	}
//line views/show_space.qtpl:93
	qw422016.N().S(`">
                            <label class="c-label" for="folder-name">Folder name</label>
                            <div class="row">
                                <div class="col-md-6">
                                    `)
//line views/show_space.qtpl:97
	if err := v.ValidationErrors.Get("name"); err != nil {
//line views/show_space.qtpl:97
		qw422016.N().S(`
                                    <input class="form-control is-invalid" type="text" value="`)
//line views/show_space.qtpl:98
		qw422016.E().S(v.Folder.Name)
//line views/show_space.qtpl:98
		qw422016.N().S(`" id="folder-name" name="name" aria-invalid="true" aria-describedby="folder-name-invalid">
                                    <small class="invalid-feedback" id="folder-name-invalid">`)
//line views/show_space.qtpl:99
		qw422016.E().S(err.Error())
//line views/show_space.qtpl:99
		qw422016.N().S(`</small>
                                    `)
//line views/show_space.qtpl:100
	} else {
//line views/show_space.qtpl:100
		qw422016.N().S(`
                                    <input class="form-control" type="text" value="`)
//line views/show_space.qtpl:101
		qw422016.E().S(v.Folder.Name)
//line views/show_space.qtpl:101
		qw422016.N().S(`" id="folder-name" name="name">
                                    `)
//line views/show_space.qtpl:102
	}
//line views/show_space.qtpl:102
	qw422016.N().S(`
                                    <small class="form-text text-muted">
                                        We will generate a shareable public link for you.
                                        <br>
                                        Each folder will expire one month after creation date.
                                    </small>
                                </div>
                                <div class="col-md-3">
                                    <button class="btn btn-primary ms-4" type="submit">
                                        <i class="if if-check"></i>
                                        <span class="btn-text">Make folder</span>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="card w-100 mb-6">
                <div class="card-header">
                    <div class="bc-toolbar">
                        <div class="bc-toolbar-left">
                            <div class="bc-toolbar-item">
                                <h2>Folders</h2>
                            </div>
                        </div>
                        <div class="bc-toolbar-right">
                            <div class="bc-toolbar-item">
                                <p>Showing `)
//line views/show_space.qtpl:130
	qw422016.N().D(len(v.Space.Folders))
//line views/show_space.qtpl:130
	qw422016.N().S(` of `)
//line views/show_space.qtpl:130
	qw422016.N().D(len(v.Space.Folders))
//line views/show_space.qtpl:130
	qw422016.N().S(` folders</p>
                            </div>
                        </div>
                    </div>
                </div>
                `)
//line views/show_space.qtpl:135
	if len(v.Space.Folders) > 0 {
//line views/show_space.qtpl:135
		qw422016.N().S(`
                <div class="table-responsive">
                    <table class="table table-sm table-bordered">
                        <thead>
                            <tr>
                                <th class="table-col-lg-fixed table-col-sm-fixed-left text-nowrap">Folder</th>
                                <th class="text-nowrap">Public shareable link</th>
                                <th class="text-nowrap">Expires on</th>
                                <th class="text-nowrap">Documents</th>
                                <th class="text-nowrap">Created at</th>
                                <th class="text-nowrap">Updated at</th>
                                <th class="table-col-sm-fixed table-col-sm-fixed-right"></th>
                            </tr>
                        </thead>
                        <tbody>
                            `)
//line views/show_space.qtpl:150
		for _, f := range v.Space.Folders {
//line views/show_space.qtpl:150
			qw422016.N().S(`
                            <tr>
                                <td class="text-nowrap table-col-lg-fixed table-col-sm-fixed-left">
                                    <a href="`)
//line views/show_space.qtpl:153
			qw422016.E().S(c.PathTo("folder", "folderID", f.ID).String())
//line views/show_space.qtpl:153
			qw422016.N().S(`">`)
//line views/show_space.qtpl:153
			qw422016.E().S(f.Name)
//line views/show_space.qtpl:153
			qw422016.N().S(`</a>
                                </td>
                                <td class="text-nowrap">
                                    <div class="input-group" style="min-width: 375px;">
                                        <button type="button" class="btn btn-outline-secondary"
                                            data-clipboard="`)
//line views/show_space.qtpl:158
			qw422016.E().S(c.URLTo("shareFolder", "folderID", f.ID, "folderSlug", f.Slug()).String())
//line views/show_space.qtpl:158
			qw422016.N().S(`"
                                        >
                                            <i class="if if-copy text-primary"></i>
                                            <span class="btn-text">Copy link</span>
                                        </button>
                                        <input id="" type="text" class="form-control input-select-text" style="min-width: 250px;" readonly
                                            value="`)
//line views/show_space.qtpl:164
			qw422016.E().S(c.URLTo("shareFolder", "folderID", f.ID, "folderSlug", f.Slug()).String())
//line views/show_space.qtpl:164
			qw422016.N().S(`"
                                            data-select-value
                                        >
                                    </div>
                                </td>
                                <td class="text-nowrap">
                                    <p>`)
//line views/show_space.qtpl:170
			qw422016.E().S(f.ExpiresAt.In(c.Timezone).Format("2006-01-02 15:04"))
//line views/show_space.qtpl:170
			qw422016.N().S(`</p>
                                    `)
//line views/show_space.qtpl:171
			if time.Until(f.ExpiresAt) < time.Hour*24*7 {
//line views/show_space.qtpl:171
				qw422016.N().S(`
                                    <p class="badge rounded-pill badge-default mt-2">
                                        <i class="if if-info-circle"></I>
                                        <span class="badge-text">Expires in `)
//line views/show_space.qtpl:174
				qw422016.E().S(friendly.TimeRemaining(time.Until(f.ExpiresAt), friendly.EnglishTimeUnits))
//line views/show_space.qtpl:174
				qw422016.N().S(`.</span>
                                    </p>
                                    `)
//line views/show_space.qtpl:176
			}
//line views/show_space.qtpl:176
			qw422016.N().S(`
                                </td>
                                <td class="text-nowrap">
                                    <p>`)
//line views/show_space.qtpl:179
			qw422016.N().D(len(f.Files))
//line views/show_space.qtpl:179
			qw422016.N().S(` files</p>
                                    <ul class="c-meta-list c-meta-list-horizontal">
                                        <li class="c-meta-item">
                                            <span>`)
//line views/show_space.qtpl:182
			qw422016.E().S(friendly.Bytes(f.TotalSize()))
//line views/show_space.qtpl:182
			qw422016.N().S(`</span>
                                        </li>
                                        <li class="c-meta-item">
                                            <span>`)
//line views/show_space.qtpl:185
			qw422016.N().DL(f.TotalDownloads())
//line views/show_space.qtpl:185
			qw422016.N().S(` downloads</span>
                                        </li>
                                    </ul>
                                </td>
                                <td class="text-nowrap">
                                    <p>`)
//line views/show_space.qtpl:190
			qw422016.E().S(f.CreatedAt.In(c.Timezone).Format("2006-01-02 15:04"))
//line views/show_space.qtpl:190
			qw422016.N().S(`</p>
                                </td>
                                <td class="text-nowrap">
                                    <p>`)
//line views/show_space.qtpl:193
			qw422016.E().S(f.UpdatedAt.In(c.Timezone).Format("2006-01-02 15:04"))
//line views/show_space.qtpl:193
			qw422016.N().S(`</p>
                                </td>
                                <td class="table-col-sm-fixed table-col-sm-fixed-right">
                                    <div class="c-button-toolbar">
                                        <a class="btn btn-link" href="`)
//line views/show_space.qtpl:197
			qw422016.E().S(c.PathTo("folder", "folderID", f.ID).String())
//line views/show_space.qtpl:197
			qw422016.N().S(`">
                                            <i class="if if-draft"></i>
                                            <span class="btn-text">Open</span>
                                        </a>
                                    </div>
                                </td>
                            </tr>
                            `)
//line views/show_space.qtpl:204
		}
//line views/show_space.qtpl:204
		qw422016.N().S(`
                        </tbody>
                    </table>
                </div>
                `)
//line views/show_space.qtpl:208
	} else {
//line views/show_space.qtpl:208
		qw422016.N().S(`
                <div class="c-blank-slate c-blank-slate-muted">
                    <div class="bc-avatar">
                        <i class="if if-info-circle"></i>
                    </div>
                    <p>Make a folder to get started</p>
                </div>
                `)
//line views/show_space.qtpl:215
	}
//line views/show_space.qtpl:215
	qw422016.N().S(`
            </div>
        </div>
    </div>
`)
//line views/show_space.qtpl:219
}

//line views/show_space.qtpl:219
func (v *ShowSpace) WriteContent(qq422016 qtio422016.Writer, c *ctx.Ctx) {
//line views/show_space.qtpl:219
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/show_space.qtpl:219
	v.StreamContent(qw422016, c)
//line views/show_space.qtpl:219
	qt422016.ReleaseWriter(qw422016)
//line views/show_space.qtpl:219
}

//line views/show_space.qtpl:219
func (v *ShowSpace) Content(c *ctx.Ctx) string {
//line views/show_space.qtpl:219
	qb422016 := qt422016.AcquireByteBuffer()
//line views/show_space.qtpl:219
	v.WriteContent(qb422016, c)
//line views/show_space.qtpl:219
	qs422016 := string(qb422016.B)
//line views/show_space.qtpl:219
	qt422016.ReleaseByteBuffer(qb422016)
//line views/show_space.qtpl:219
	return qs422016
//line views/show_space.qtpl:219
}

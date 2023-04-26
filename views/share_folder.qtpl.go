// Code generated by qtc from "share_folder.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/share_folder.qtpl:1
package views

//line views/share_folder.qtpl:1
import "github.com/ugent-library/friendly"

//line views/share_folder.qtpl:2
import "github.com/ugent-library/deliver/ctx"

//line views/share_folder.qtpl:3
import "github.com/ugent-library/deliver/models"

//line views/share_folder.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/share_folder.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/share_folder.qtpl:6
type ShareFolder struct {
	Folder *models.Folder
}

//line views/share_folder.qtpl:11
func (v *ShareFolder) StreamTitle(qw422016 *qt422016.Writer, c *ctx.Ctx) {
//line views/share_folder.qtpl:11
	qw422016.E().S(v.Folder.Name)
//line views/share_folder.qtpl:11
}

//line views/share_folder.qtpl:11
func (v *ShareFolder) WriteTitle(qq422016 qtio422016.Writer, c *ctx.Ctx) {
//line views/share_folder.qtpl:11
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/share_folder.qtpl:11
	v.StreamTitle(qw422016, c)
//line views/share_folder.qtpl:11
	qt422016.ReleaseWriter(qw422016)
//line views/share_folder.qtpl:11
}

//line views/share_folder.qtpl:11
func (v *ShareFolder) Title(c *ctx.Ctx) string {
//line views/share_folder.qtpl:11
	qb422016 := qt422016.AcquireByteBuffer()
//line views/share_folder.qtpl:11
	v.WriteTitle(qb422016, c)
//line views/share_folder.qtpl:11
	qs422016 := string(qb422016.B)
//line views/share_folder.qtpl:11
	qt422016.ReleaseByteBuffer(qb422016)
//line views/share_folder.qtpl:11
	return qs422016
//line views/share_folder.qtpl:11
}

//line views/share_folder.qtpl:13
func (v *ShareFolder) StreamContent(qw422016 *qt422016.Writer, c *ctx.Ctx) {
//line views/share_folder.qtpl:13
	qw422016.N().S(`
    <div class="d-flex u-maximize-height">
        <div class="w-100 u-scroll-wrapper">
            <div class="bg-white">
                <div class="bc-navbar bc-navbar--xlarge bc-navbar--white bc-navbar--bordered-bottom">
                    <div class="w-100">
                        <div class="bc-toolbar bc-toolbar--auto">
                            <div class="bc-toolbar-left">
                                <div class="bc-toolbar-item">
                                    <h4 class="bc-toolbar-title">Library delivery from `)
//line views/share_folder.qtpl:22
	qw422016.E().S(v.Folder.Space.Name)
//line views/share_folder.qtpl:22
	qw422016.N().S(`: `)
//line views/share_folder.qtpl:22
	qw422016.E().S(v.Folder.Name)
//line views/share_folder.qtpl:22
	qw422016.N().S(`</h4>
                                </div>
                                <div class="bc-toolbar-item">
                                    <p>Expires on `)
//line views/share_folder.qtpl:25
	qw422016.E().S(v.Folder.ExpiresAt.Format("2006-01-02 15:04"))
//line views/share_folder.qtpl:25
	qw422016.N().S(`</p>
                                </div>
                            </div>
                        </div>
                        <div class="bc-toolbar bc-toolbar--auto mt-2">
                            <div class="bc-toolbar-left">
                                <div class="bc-toolbar-item">
                                    <p class="text-muted">
                                        Public shareable link: `)
//line views/share_folder.qtpl:33
	qw422016.E().S(c.URLTo("share_folder", "folderID", v.Folder.ID, "folderSlug", v.Folder.Slug()).String())
//line views/share_folder.qtpl:33
	qw422016.N().S(`
                                    </p>
                                </div>
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
                                <div class="bc-toolbar-item">
                                    <h5>Download your files before `)
//line views/share_folder.qtpl:47
	qw422016.E().S(v.Folder.ExpiresAt.Format("2006-01-02 15:04"))
//line views/share_folder.qtpl:47
	qw422016.N().S(`</h5>
                                </div>
                            </div>
                            <div class="bc-toolbar-right">
                                <div class="bc-toolbar-item">
                                    <p>`)
//line views/share_folder.qtpl:52
	qw422016.N().D(len(v.Folder.Files))
//line views/share_folder.qtpl:52
	qw422016.N().S(` items</p>
                                </div>
                            </div>
                        </div>
                    </div>
                    `)
//line views/share_folder.qtpl:57
	if len(v.Folder.Files) > 0 {
//line views/share_folder.qtpl:57
		qw422016.N().S(`
                    <div class="table-responsive">
                        <table class="table table-sm table-bordered">
                            <thead>
                                <tr>
                                    <th class="table-col-lg-fixed table-col-sm-fixed-left text-nowrap align-middle">File name</th>
                                    <th class="text-nowrap align-middle">Size</th>
                                    <th class="text-nowrap align-middle">Type</th>
                                    <th class="text-nowrap align-middle">Downloads</th>
                                    <th class="text-nowrap align-middle">Created at</th>
                                    <th class="table-col-sm-fixed table-col-sm-fixed-right text-right align-middle">
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                `)
//line views/share_folder.qtpl:72
		for _, f := range v.Folder.Files {
//line views/share_folder.qtpl:72
			qw422016.N().S(`
                                <tr class="clickable-table-row">
                                    <td class="text-nowrap table-col-lg-fixed table-col-sm-fixed-left">
                                        <a href="`)
//line views/share_folder.qtpl:75
			qw422016.E().S(c.PathTo("download_file", "fileID", f.ID).String())
//line views/share_folder.qtpl:75
			qw422016.N().S(`">
                                            <span>`)
//line views/share_folder.qtpl:76
			qw422016.E().S(f.Name)
//line views/share_folder.qtpl:76
			qw422016.N().S(`</span>
                                        </a>
                                        <br><small class="text-muted">md5 checksum: `)
//line views/share_folder.qtpl:78
			qw422016.E().S(f.MD5)
//line views/share_folder.qtpl:78
			qw422016.N().S(`</small>
                                    </td>
                                    <td class="text-nowrap">
                                        <p>`)
//line views/share_folder.qtpl:81
			qw422016.E().S(friendly.Bytes(f.Size))
//line views/share_folder.qtpl:81
			qw422016.N().S(`</p>
                                    </td>
                                    <td class="text-nowrap">
                                        <p>`)
//line views/share_folder.qtpl:84
			qw422016.E().S(f.ContentType)
//line views/share_folder.qtpl:84
			qw422016.N().S(`</p>
                                    </td>
                                    <td class="text-nowrap">
                                        <p>`)
//line views/share_folder.qtpl:87
			qw422016.N().DL(f.Downloads)
//line views/share_folder.qtpl:87
			qw422016.N().S(`</p>
                                    </td>
                                    <td class="text-nowrap">
                                        <p>`)
//line views/share_folder.qtpl:90
			qw422016.E().S(f.CreatedAt.Format("2006-01-02 15:04"))
//line views/share_folder.qtpl:90
			qw422016.N().S(`</p>
                                    </td>
                                    <td class="table-col-sm-fixed table-col-sm-fixed-right">
                                        <div class="c-button-toolbar">
                                            <a class="btn btn-link" href="`)
//line views/share_folder.qtpl:94
			qw422016.E().S(c.PathTo("download_file", "fileID", f.ID).String())
//line views/share_folder.qtpl:94
			qw422016.N().S(`">
                                                <i class="if if-download"></i>
                                                <span class="btn-txt">Download</span>
                                            </a>
                                        </div>
                                    </td>
                                </tr>
                                `)
//line views/share_folder.qtpl:101
		}
//line views/share_folder.qtpl:101
		qw422016.N().S(`
                            </tbody>
                        </table>
                    </div>
                    `)
//line views/share_folder.qtpl:105
	} else {
//line views/share_folder.qtpl:105
		qw422016.N().S(`
                    <div class="c-blank-slate c-blank-slate-muted">
                        <div class="bc-avatar">
                            <i class="if if-info-circle"></i>
                        </div>
                        <p>
                            No files to display.
                            <br>
                            Please get in touch with the person that sent you this link.
                        </p>
                    </div>
                    `)
//line views/share_folder.qtpl:116
	}
//line views/share_folder.qtpl:116
	qw422016.N().S(`
                </div>
            </div>
        </div>
    </div>
`)
//line views/share_folder.qtpl:121
}

//line views/share_folder.qtpl:121
func (v *ShareFolder) WriteContent(qq422016 qtio422016.Writer, c *ctx.Ctx) {
//line views/share_folder.qtpl:121
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/share_folder.qtpl:121
	v.StreamContent(qw422016, c)
//line views/share_folder.qtpl:121
	qt422016.ReleaseWriter(qw422016)
//line views/share_folder.qtpl:121
}

//line views/share_folder.qtpl:121
func (v *ShareFolder) Content(c *ctx.Ctx) string {
//line views/share_folder.qtpl:121
	qb422016 := qt422016.AcquireByteBuffer()
//line views/share_folder.qtpl:121
	v.WriteContent(qb422016, c)
//line views/share_folder.qtpl:121
	qs422016 := string(qb422016.B)
//line views/share_folder.qtpl:121
	qt422016.ReleaseByteBuffer(qb422016)
//line views/share_folder.qtpl:121
	return qs422016
//line views/share_folder.qtpl:121
}
// Code generated by qtc from "show_folder.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/show_folder.qtpl:1
package views

//line views/show_folder.qtpl:1
import "github.com/ugent-library/friendly"

//line views/show_folder.qtpl:2
import "github.com/ugent-library/deliver/ctx"

//line views/show_folder.qtpl:3
import "github.com/ugent-library/deliver/models"

//line views/show_folder.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/show_folder.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/show_folder.qtpl:6
type ShowFolder struct {
	Folder      *models.Folder
	MaxFileSize int64
}

//line views/show_folder.qtpl:12
func (v *ShowFolder) StreamTitle(qw422016 *qt422016.Writer, c *ctx.Ctx) {
//line views/show_folder.qtpl:12
	qw422016.E().S(v.Folder.Name)
//line views/show_folder.qtpl:12
}

//line views/show_folder.qtpl:12
func (v *ShowFolder) WriteTitle(qq422016 qtio422016.Writer, c *ctx.Ctx) {
//line views/show_folder.qtpl:12
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/show_folder.qtpl:12
	v.StreamTitle(qw422016, c)
//line views/show_folder.qtpl:12
	qt422016.ReleaseWriter(qw422016)
//line views/show_folder.qtpl:12
}

//line views/show_folder.qtpl:12
func (v *ShowFolder) Title(c *ctx.Ctx) string {
//line views/show_folder.qtpl:12
	qb422016 := qt422016.AcquireByteBuffer()
//line views/show_folder.qtpl:12
	v.WriteTitle(qb422016, c)
//line views/show_folder.qtpl:12
	qs422016 := string(qb422016.B)
//line views/show_folder.qtpl:12
	qt422016.ReleaseByteBuffer(qb422016)
//line views/show_folder.qtpl:12
	return qs422016
//line views/show_folder.qtpl:12
}

//line views/show_folder.qtpl:14
func (v *ShowFolder) StreamContent(qw422016 *qt422016.Writer, c *ctx.Ctx) {
//line views/show_folder.qtpl:14
	qw422016.N().S(`
    `)
//line views/show_folder.qtpl:15
	qw422016.N().S(c.TurboStreamTag("folder." + v.Folder.ID))
//line views/show_folder.qtpl:15
	qw422016.N().S(`

    <div class="w-100 u-scroll-wrapper">
        <div class="bg-white">
            <div class="border-bottom py-5 px-6">
                <div class="bc-toolbar bc-toolbar--auto">
                    <div class="bc-toolbar-left">
                        <div class="bc-toolbar-item">
                            <h4 class="bc-toolbar-title">
                                <a href="`)
//line views/show_folder.qtpl:24
	qw422016.E().S(c.PathTo("space", "spaceName", v.Folder.Space.Name).String())
//line views/show_folder.qtpl:24
	qw422016.N().S(`">
                                    <i class="if if-arrow-left"></i>
                                    <span>`)
//line views/show_folder.qtpl:26
	qw422016.E().S(v.Folder.Space.Name)
//line views/show_folder.qtpl:26
	qw422016.N().S(`</span>
                                </a>
                                &nbsp;&mdash; `)
//line views/show_folder.qtpl:28
	qw422016.E().S(v.Folder.Name)
//line views/show_folder.qtpl:28
	qw422016.N().S(`
                            </h4>
                        </div>
                        <div class="bc-toolbar-item">
                            <p>expires on `)
//line views/show_folder.qtpl:32
	qw422016.E().S(v.Folder.ExpiresAt.Format("2006-01-02 15:04"))
//line views/show_folder.qtpl:32
	qw422016.N().S(`</p>
                        </div>
                    </div>
                    <div class="bc-toolbar-right">
                        <div class="bc-toolbar-item">
                            <a class="btn btn-link btn-link-muted" href="`)
//line views/show_folder.qtpl:37
	qw422016.E().S(c.PathTo("space", "spaceName", v.Folder.Space.Name).String())
//line views/show_folder.qtpl:37
	qw422016.N().S(`">
                                <i class="if if-add"></i>
                                <span class="btn-text">Make new folder</span>
                            </a>
                        </div>
                        <div class="bc-toolbar-item">
                            <a class="btn btn-link btn-link-muted" href="`)
//line views/show_folder.qtpl:43
	qw422016.E().S(c.PathTo("edit_folder", "folderID", v.Folder.ID).String())
//line views/show_folder.qtpl:43
	qw422016.N().S(`">
                                <i class="if if-edit"></i>
                                <span class="btn-text">Edit</span>
                            </a>
                        </div>
                    </div>
                </div>
                <div class="bc-toolbar bc-toolbar--auto mt-3 mb-2">
                    <div class="bc-toolbar-left">
                        <div class="bc-toolbar-item col-lg-8 col-sm-12 pl-0">
                            <div class="input-group" data-controller="clipboard">
                                <div class="input-group-prepend">
                                    <button type="button" class="btn btn-outline-secondary"
                                        data-action="clipboard#copy"
                                        data-clipboard-target="button"
                                    >
                                        <i class="if if-copy text-primary"></i>
                                        <span class="btn-text text-primary">Copy public shareable link</span>
                                    </button>
                                </div>
                                <input type="text" class="form-control input-select-text" readonly
                                    value="`)
//line views/show_folder.qtpl:64
	qw422016.E().S(c.URLTo("share_folder", "folderID", v.Folder.ID, "folderSlug", v.Folder.Slug()).String())
//line views/show_folder.qtpl:64
	qw422016.N().S(`"
                                    data-action="click->clipboard#select"
                                    data-clipboard-target="source"
                                >
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="u-scroll-wrapper__body p-6">
            <div class="card w-100 mb-6">
                <div class="card-header sr-only">
                    <div class="bc-toolbar">
                        <div class="bc-toolbar-left">
                            <h2 class="bc-toolbar-item">Upload files</h2>
                        </div>
                    </div>
                </div>
                <div class="card-body">
                    <form
                        action="`)
//line views/show_folder.qtpl:86
	qw422016.E().S(c.PathTo("upload_file", "folderID", v.Folder.ID).String())
//line views/show_folder.qtpl:86
	qw422016.N().S(`"
                        enctype="multipart/form-data"
                        method="POST"
                    >
                        `)
//line views/show_folder.qtpl:90
	qw422016.N().S(c.CSRFTag)
//line views/show_folder.qtpl:90
	qw422016.N().S(`
                        <div class="c-file-upload c-file-upload--small">
                            <input type="file" name="file" multiple
                                data-upload-progress-target="file-upload-progress"
                                data-upload-max-file-size="`)
//line views/show_folder.qtpl:94
	qw422016.N().DL(v.MaxFileSize)
//line views/show_folder.qtpl:94
	qw422016.N().S(`"
                                data-upload-msg-file-aborted="File upload aborted by you"
                                data-upload-msg-file-uploading="Uploading your file. Hold on, do not refresh the page."
                                data-upload-msg-file-processing="Processing your file. Hold on, do not refresh the page."
                                data-upload-msg-file-too-large="File is too large. Maximum file size is `)
//line views/show_folder.qtpl:98
	qw422016.E().S(friendly.Bytes(v.MaxFileSize))
//line views/show_folder.qtpl:98
	qw422016.N().S(`"
                                data-upload-msg-dir-not-found="File upload failed: directory has been removed. Please reload"
                                data-upload-msg-unexpected="File upload failed: unexpected server error"
                            >
                            <div class="c-file-upload__content">
                                <p class="pt-2">Drag and drop or</p>
                                <button class="btn btn-outline-primary" data-loading="Uploading...">
                                    Upload files
                                </button>
                                <p class="small pt-2 mb-0">Maximum file size: `)
//line views/show_folder.qtpl:107
	qw422016.E().S(friendly.Bytes(v.MaxFileSize))
//line views/show_folder.qtpl:107
	qw422016.N().S(`</p>
                            </div>
                        </div>
                    </form>
                    <ul class="list-group list-group-flush" id="file-upload-progress"></ul>
                </div>
            </div>

            `)
//line views/show_folder.qtpl:115
	StreamFiles(qw422016, c, v.Folder.Files)
//line views/show_folder.qtpl:115
	qw422016.N().S(`
        </div>
    </div>

    <template id="tmpl-upload-progress">
        <li class="list-group-item">
            <div class="list-group-item-inner">
                <div class="list-group-item-main u-min-w-0">
                    <div class="bc-toolbar bc-toolbar--auto">
                        <div class="bc-toolbar-left">
                            <div class="bc-toolbar-item">
                                <h4 class="list-group-item-title mb-0 upload-name"></h4>
                            </div>
                        </div>
                        <div class="bc-toolbar-right">
                            <div class="bc-toolbar-item ml-auto ml-lg-0">
                                <button class="btn btn-sm btn-cancel-upload" type="button">
                                    <i class="if if-close"></i>
                                    <span class="btn-text">Cancel upload</span>
                                </button>
                                <button class="btn btn-sm btn-remove-upload d-none" type="button">
                                    <i class="if if-close"></i>
                                    <span class="btn-text">Remove</span>
                                </button>
                            </div>
                        </div>
                    </div>
                    <div class="progress w-100 mt-4">
                        <div class="progress-bar progress-bar-striped" role="progressbar" style="width:0%" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100"></div>
                    </div>
                    <div class="bc-toolbar bc-toolbar--auto mt-2">
                        <div class="bc-toolbar-left">
                            <span class="small text-muted upload-msg"></span>
                        </div>
                        <div class="bc-toolbar-right">
                            <span class="small text-muted"><span class="upload-size"></span> &mdash; <span class="upload-percent">0</span>%</span>
                        </div>
                    </div>
                </div>
            </div>
        </li>
    </template>
`)
//line views/show_folder.qtpl:157
}

//line views/show_folder.qtpl:157
func (v *ShowFolder) WriteContent(qq422016 qtio422016.Writer, c *ctx.Ctx) {
//line views/show_folder.qtpl:157
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/show_folder.qtpl:157
	v.StreamContent(qw422016, c)
//line views/show_folder.qtpl:157
	qt422016.ReleaseWriter(qw422016)
//line views/show_folder.qtpl:157
}

//line views/show_folder.qtpl:157
func (v *ShowFolder) Content(c *ctx.Ctx) string {
//line views/show_folder.qtpl:157
	qb422016 := qt422016.AcquireByteBuffer()
//line views/show_folder.qtpl:157
	v.WriteContent(qb422016, c)
//line views/show_folder.qtpl:157
	qs422016 := string(qb422016.B)
//line views/show_folder.qtpl:157
	qt422016.ReleaseByteBuffer(qb422016)
//line views/show_folder.qtpl:157
	return qs422016
//line views/show_folder.qtpl:157
}

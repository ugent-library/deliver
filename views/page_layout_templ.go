// Code generated by templ@v0.2.334 DO NOT EDIT.

package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/ugent-library/deliver/ctx"

func pageLayout(c *ctx.Ctx, title string, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<!doctype html><html class=\"u-maximize-height\" dir=\"ltr\" lang=\"en\"><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><meta name=\"robots\" content=\"noindex\"><meta name=\"csrf-token\" content=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(c.CSRFToken))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><link rel=\"stylesheet\" href=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(c.AssetPath("/css/app.css")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><link rel=\"icon\" href=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(c.AssetPath("/favicon.ico")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><script type=\"application/javascript\" src=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(c.AssetPath("/js/app.js")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		var_2 := ``
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><title>")
		if err != nil {
			return err
		}
		var var_3 string = title
		_, err = templBuffer.WriteString(templ.EscapeString(var_3))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</title></head><body class=\"u-maximize-height overflow-hidden u-scroll-wrapper\" hx-swap=\"none\"><header>")
		if err != nil {
			return err
		}
		err = pageBanner(c).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div class=\"bc-navbar bc-navbar--small bc-navbar--bordered-bottom bc-navbar--white bc-navbar--fixed bc-navbar--scrollable shadow-sm px-4\"><div class=\"bc-toolbar bc-toolbar-sm\"><div class=\"bc-toolbar-left\"><div class=\"bc-toolbar-item\"><nav aria-label=\"breadcrumb\"><ol class=\"breadcrumb\"><li class=\"breadcrumb-item\"><a href=\"")
		if err != nil {
			return err
		}
		var var_4 templ.SafeURL = templ.URL(c.PathTo("home").String())
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_4)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><img class=\"d-none d-lg-inline-block\" src=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(c.AssetPath("/images/ghent-university-library-logo.svg")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" alt=\"Ghent University Library\"><img class=\"d-inline-block d-lg-none\" src=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(c.AssetPath("/images/ghent-university-library-mark.svg")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" alt=\"Ghent University Library\"></a></li><li class=\"breadcrumb-item\" aria-current=\"page\"><a href=\"")
		if err != nil {
			return err
		}
		var var_5 templ.SafeURL = templ.URL(c.PathTo("home").String())
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_5)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" class=\"text-muted\">")
		if err != nil {
			return err
		}
		var_6 := `Home`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></li></ol></nav></div></div><div class=\"bc-toolbar-right\"><div class=\"bc-toolbar-item\"><div id=\"side-panels\"><ul class=\"nav nav-main\">")
		if err != nil {
			return err
		}
		if c.User != nil {
			_, err = templBuffer.WriteString("<li class=\"nav-item\"><a class=\"nav-link\" href=\"mailto:libservice@ugent.be\" target=\"_blank\"><i class=\"if if-service\"></i><span class=\"btn-text\">")
			if err != nil {
				return err
			}
			var_7 := `Ask help`
			_, err = templBuffer.WriteString(var_7)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span></a></li>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<li class=\"nav-item\"><a class=\"nav-link\" href=\"https://www.ugent.be/intranet/nl/op-het-werk/bibliotheek/publieksdiensten/deliverhandleiding\" target=\"_blank\"><i class=\"if if-book\"></i><span class=\"btn-text\">")
		if err != nil {
			return err
		}
		var_8 := `Manual &ndash; NL`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></a></li><li class=\"nav-item\">")
		if err != nil {
			return err
		}
		if c.User != nil {
			_, err = templBuffer.WriteString("<div class=\"dropdown position-static\"><button class=\"nav-link dropdown-toggle\" role=\"button\" data-bs-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\"><div class=\"bc-avatar-and-text\"><div class=\"bc-avatar bc-avatar--sm bc-avatar--small me-2\"><i class=\"if if-user\"></i></div><div class=\"bc-avatar-text\">")
			if err != nil {
				return err
			}
			var var_9 string = c.User.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div></div></button><div class=\"dropdown-menu\"><div class=\"bc-avatar-and-text m-4\"><div class=\"bc-avatar bc-avatar--sm\"><i class=\"if if-user\"></i></div><div class=\"bc-avatar-text\"><h4>")
			if err != nil {
				return err
			}
			var var_10 string = c.User.Name
			_, err = templBuffer.WriteString(templ.EscapeString(var_10))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h4><p class=\"text-muted c-body-small\">")
			if err != nil {
				return err
			}
			var var_11 string = c.User.Email
			_, err = templBuffer.WriteString(templ.EscapeString(var_11))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p></div></div><hr class=\"dropdown-divider\"><a class=\"dropdown-item\" href=\"")
			if err != nil {
				return err
			}
			var var_12 templ.SafeURL = templ.URL(c.PathTo("logout").String())
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_12)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"><i class=\"if if-log-out\"></i><span>")
			if err != nil {
				return err
			}
			var_13 := `Log out`
			_, err = templBuffer.WriteString(var_13)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span></a></div></div>")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString("<a class=\"btn btn-link btn-sm\" href=\"")
			if err != nil {
				return err
			}
			var var_14 templ.SafeURL = templ.URL(c.PathTo("login").String())
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_14)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"><i class=\"if if-arrow-right mt-0 ms-2\"></i><span class=\"btn-text me-2\">")
			if err != nil {
				return err
			}
			var_15 := `Log in`
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</span></a>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</li></ul></div></div></div></div></div></header><main><div class=\"d-flex u-maximize-height\">")
		if err != nil {
			return err
		}
		var var_16 = []any{"c-sidebar", templ.KV("c-sidebar--dark-gray", c.User != nil), "d-none", "d-lg-flex"}
		err = templ.RenderCSSItems(ctx, templBuffer, var_16...)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div class=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(templ.CSSClasses(var_16).String()))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><div class=\"c-sidebar__menu\"><nav><ul class=\"c-sidebar-menu\"><li class=\"c-sidebar__item c-sidebar__item--active\"><a href=\"")
		if err != nil {
			return err
		}
		var var_17 templ.SafeURL = templ.URL(c.PathTo("home").String())
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_17)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"><span class=\"c-sidebar__icon\"><i class=\"if if-file\"></i></span><span class=\"c-sidebar__label\">")
		if err != nil {
			return err
		}
		var_18 := `Deliver`
		_, err = templBuffer.WriteString(var_18)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span></a></li></ul></nav></div><div class=\"c-sidebar__bottom\"><img src=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(c.AssetPath("/images/logo-ugent-white.svg")))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" alt=\"Logo UGent\" height=\"48px\" width=\"auto\"></div></div>")
		if err != nil {
			return err
		}
		err = content.Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></main><div id=\"flash-messages\">")
		if err != nil {
			return err
		}
		for _, f := range c.Flash {
			err = flash(f).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div><template id=\"modal-confirm\"><div class=\"modal fade\" aria-modal=\"true\" aria-hidden=\"true\" role=\"dialog\"><div class=\"modal-dialog modal-dialog-centered\" role=\"document\"><div class=\"modal-content\"><div class=\"modal-body\"><div class=\"c-blank-slate c-blank-slate-muted\"><div class=\"bc-avatar\"><i class=\"if if-alert\"></i></div><h4 class=\"confirm-header\">")
		if err != nil {
			return err
		}
		var_19 := `Are you sure?`
		_, err = templBuffer.WriteString(var_19)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4><p>")
		if err != nil {
			return err
		}
		var_20 := `You cannot undo this action.`
		_, err = templBuffer.WriteString(var_20)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p></div></div><div class=\"modal-footer\"><button class=\"btn btn-link\" data-bs-dismiss=\"modal\">")
		if err != nil {
			return err
		}
		var_21 := `No, cancel`
		_, err = templBuffer.WriteString(var_21)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button><button class=\"btn btn-danger confirm-proceed\" data-bs-dismiss=\"modal\">")
		if err != nil {
			return err
		}
		var_22 := `Yes, proceed`
		_, err = templBuffer.WriteString(var_22)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></div></div></div></div></template></body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

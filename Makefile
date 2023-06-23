.PHONY: install-dev dev

install-dev:
	npm install
	go install github.com/valyala/quicktemplate/qtc@latest
	go install github.com/cespare/reflex@latest

dev:
	env $$(cat .env|xargs) reflex -d none -c reflex.conf

TAILWIND_SRC := src/css/tailwind.css
TAILWIND_OUT := public/assets/css/global.css

.PHONY: generate
generate:
	go run github.com/a-h/templ/cmd/templ@v0.3.1070 generate

.PHONY: frontend
frontend:
	npx @tailwindcss/cli -i $(TAILWIND_SRC) -o $(TAILWIND_OUT) --minify

.PHONY: watch
watch:
	$(MAKE) -j2 watch-frontend watch-backend

.PHONY: watch-frontend
watch-frontend:
	npx @tailwindcss/cli -i $(TAILWIND_SRC) -o $(TAILWIND_OUT) --watch

.PHONY: watch-backend
watch-backend:
	go run github.com/a-h/templ/cmd/templ@v0.3.1070 generate --watch --proxy=http://localhost:8080 --open-browser=false --cmd="go run -tags dev ./cmd/server"

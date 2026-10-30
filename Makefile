TAILWIND_SRC := src/css/tailwind.css
TAILWIND_OUT := public/assets/css/global.css

.PHONY: generate
generate:
	go run github.com/a-h/templ/cmd/templ@v0.3.1070 generate

.PHONY: frontend
frontend:
	npx @tailwindcss/cli -i $(TAILWIND_SRC) -o $(TAILWIND_OUT) --minify

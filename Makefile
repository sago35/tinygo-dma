.PHONY: smoketest fmt-check

smoketest:
	tinygo build -o app.hex -target feather-m4 -size short ./examples/adc
	tinygo build -o app.hex -target feather-m4 -size short ./examples/adctc
	tinygo build -o app.hex -target feather-m4 -size short ./examples/dac
	tinygo build -o app.hex -target feather-m4 -size short ./examples/dactc
	tinygo build -o app.hex -target feather-m4 -size short ./examples/memcpy
	tinygo build -o app.hex -target feather-m4 -size short ./examples/spitx
	tinygo build -o app.hex -target feather-m4 -size short ./examples/spitxrx
	tinygo build -o app.hex -target feather-m4 -size short ./examples/uarttx
	tinygo build -o app.hex -target feather-m4 -size short ./examples/uarttx_with_linked_descriptor
	tinygo build -o app.hex -target feather-m4 -size short ./examples/uarttxrx

fmt-check:
	@unformatted=$$(gofmt -l `find . -name "*.go"`); [ -z "$$unformatted" ] && exit 0; echo "Unformatted:"; for fn in $$unformatted; do echo "  $$fn"; done; exit 1

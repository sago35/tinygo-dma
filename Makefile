.PHONY: smoketest

smoketest:
	tinygo build -o app.hex -target feather-m4 -size short ./examples/adc
	tinygo build -o app.hex -target feather-m4 -size short ./examples/adctc
	tinygo build -o app.hex -target feather-m4 -size short ./examples/dac
	tinygo build -o app.hex -target feather-m4 -size short ./examples/dactc
	tinygo build -o app.hex -target feather-m4 -size short ./examples/memcpy
	tinygo build -o app.hex -target feather-m4 -size short ./examples/spitx
	tinygo build -o app.hex -target feather-m4 -size short ./examples/spitxrx

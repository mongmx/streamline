dev:
	go run main.go
doc:
	swag i
tpl:
	hero -source="./templates/" -dest="./templates/t/" -pkgname="t"
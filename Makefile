doc:
	swag i
tpl:
	hero -source="./templates/" -dest="./templates/t/" -pkgname="t"
dev:
	hero -source="./templates/" -dest="./templates/t/" -pkgname="t"
	go run main.go
build-dev:
	swag i
	hero -source="./templates/" -dest="./templates/t/" -pkgname="t"
	go run main.go

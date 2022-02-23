build:
	go build -o asy ./cmd/asy

install.cli:
	@ echo "Installing asy cli from local"
	go install ./cmd/asy

install.gui:
	@ echo "Installing GUI"
	cd cmd/gui && wails build

down:
	./asy \
		--down \
		--debug --force --overwrite \
		--path=./debugdata \
		--apollo.portaladdr=http://127.0.0.1:8070 \
		--apollo.appid=demo \
		--apollo.secret=82a95a5722ae8649f64ca5859a13032acab4b2a3

up:
	./asy \
		--up \
		--debug --force --overwrite \
		--path=./debugdata \
		--apollo.portaladdr=http://127.0.0.1:8070 \
		--apollo.appid=demo \
		--apollo.secret=82a95a5722ae8649f64ca5859a13032acab4b2a3
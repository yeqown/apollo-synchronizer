build:
	go build -o asy ./cmd/asy


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
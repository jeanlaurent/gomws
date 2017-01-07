depends: clean
	go get .

all: clean
	mkdir _build
	env GOOS=linux GOARCH=386 go build -o _build/mws-bridge .
	docker build -t jeanlaurent/mws-bridge .

run:
	docker run -e MWSSellerID=${MWSSellerID} -e MWSAccessKey=${MWSAccessKey} -e MWSSecretKey=${MWSSecretKey} -p8080:8080 -t jeanlaurent/mws-bridge

clean:
	rm -rf _build

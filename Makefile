all: clean
	mkdir _build
	env GOOS=linux GOARCH=386 go build -o _build/mws-bridge .
	docker build -t jeanlaurent/mws-bridge .

clean:
	rm -rf _build

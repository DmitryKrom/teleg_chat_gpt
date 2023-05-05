build:
	docker build -t botchat .

run:
	docker run -d -p 3000:3000 --name botchat --rm botchat
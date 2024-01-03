UP_DB:
	docker run --name mysql -p "3306:3306" -e MYSQL_ROOT_PASSWORD=qwerty -e MYSQL_USER=qwerty -e MYSQL_DATABASE=db -d mysql:latest
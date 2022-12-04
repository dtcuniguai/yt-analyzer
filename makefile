##############################
# Develope
##############################

.PHONY: app

app.init:
	echo "sqlite 資料庫analyzer建立...."
	sqlite3 ./data/analyzer.db
	echo "資料表`youtuber`, `video`建立..."
	sqlite3 ./data/analyzer.db < ./data/init.sql
	echo "環境檔建立....."
	cp ./.env.example ./.env
	go build -o $(PWD)/$(shell basename $(CURDIR))

##############################
# Develope
##############################
app.build:
	go build -o $(PWD)/$(shell basename $(CURDIR))

app.dev:
	go run main.go



app.prod PLATFORM:
	echo 'run platform at $(PLATFORM)'
##############################
# Develope
##############################

.PHONY: app

app.init:
	echo "sqlite 資料庫analyzer初始化建立...."
	sqlite3 -batch ./data/analyzer.db < ./data/init.sql
	# echo "資料表`youtuber`, `video`建立..."
	sqlite3 -batch ./data/analyzer.db < ./data/init.sql
	echo "環境檔建立....."
	# cp ./.env.example ./.env
	# go build -o $(CURDIR)/cmd/ytanalyzer

##############################
# Develope
##############################
app.build:
	go build -o $(CURDIR)/cmd/ytanalyzer

app.dev:
	go run $(CURDIR)/cmd/ytanalyzer



app.prod PLATFORM:
	echo 'run platform at $(PLATFORM)'
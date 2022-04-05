#!/bin/sh

case "$1" in
	apm) # apm ui
		npm install --prefix ./HCCTV-apm
		npm run serve --prefix ./HCCTV-apm
	;;
	dev) # dev env up
		if [ "$2" == "" ] ; then
		  npm install --prefix ./HCCTV-apm
			npm run build --prefix ./HCCTV-apm
			docker compose up --build dev_db logger_db -d
			./check-db-ready.sh
			docker compose up --build echo-dev nginx 
 		else
			echo "'$2' is unknwon option"
		fi
	;;
	deploy)
		#release env up
	;;
	down) # env down
		if [ "$2" == "" ] ; then
			docker compose down
		else
			echo "-down : No option"
		fi
	;;
	*) # exception
	echo "'$1' is unknown command + '$2'"
	;;
esac

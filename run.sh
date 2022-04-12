#!/bin/sh

case "$1" in
	dev) # dev env up
<<<<<<< Updated upstream
		if [ "$2" == "" ] ; then
			docker compose up --build
		elif [ "$2" == "-d" ]; then
			docker compose up --build -d
 		else
=======
		if [ "$2" == "apm" ] ; then
		  	npm install --prefix ./HCCTV-apm
			npm run build --prefix ./HCCTV-apm
			docker compose up --build dev_db logger_db -d
			./check-db-ready.sh
			docker compose up --build echo-dev nginx 
 		elif [ "$2" == "" ] ; then
			docker compose up --build dev_db logger_db -d
			./check-db-ready.sh
			docker compose up --build echo-dev nginx 
		else
>>>>>>> Stashed changes
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

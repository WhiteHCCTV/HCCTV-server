#!/bin/sh

case "$1" in
	dev) # dev env up
		if [ "$2" == "" ] ; then
			docker compose up --build
		elif [ "$2" == "-d" ]; then
			docker compose up --build -d
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

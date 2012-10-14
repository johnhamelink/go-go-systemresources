`dirname $0`/wmfsSystem | while true; read line; do wmfs -c status "default $line"; done

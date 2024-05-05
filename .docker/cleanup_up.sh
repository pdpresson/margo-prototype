#keep ./gogs/app/gogs/config/app.ini
sudo rm -R ./gogs
mkdir -p ./gogs/app/gogs/conf
mkdir ./gogs/db_data
mkdir ./gogs/logs

cp ../config/gogs/app.ini ./gogs/app/gogs/conf/app.ini
sudo chmod -R 777 ./gogs

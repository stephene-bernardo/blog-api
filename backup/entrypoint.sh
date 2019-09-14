./docker-entrypoint.sh postgres > database.log 2>&1 &
isBlogDbExist=$(psql -U postgres -c "SELECT 1 from pg_database WHERE datname='blog'" | grep 1)
while [ -z "$isBlogDbExist" ]
do
  echo "Waiting for database..."
  sleep 5
  isBlogDbExist=$(psql -U postgres -c "SELECT 1 from pg_database WHERE datname='blog'" | grep 1)
done
echo "Initialization of database is completed. To check db logs access /database.log"
./blog-api
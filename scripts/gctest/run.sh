echo "Cleaning logs at AWS"
python ./scripts/gctest/gctest.py --command clean
echo "Compiling"
make clean
make gctest

echo "Deploying no gc"
sls deploy -c gctestnogc.yml
read -p "Please Input HTTP gateway url for beldi-dev-gctest: " gcp
echo "Initializing Database"
go run ./internal/gctest/init/init.go
echo "Running"
ENDPOINT="$gcp" ./tools/wrk -t4 -c10 -d1860s -R10 -s ./benchmark/gctest/workload.lua --timeout 10s "$gcp"
echo "Collecting metrics"
python ./scripts/gctest/gctest.py --command run --config nogc

echo "Reset Database"
go run ./internal/gctest/init/init.go
echo "Deploying 1 min gc"
sls deploy -c gctest1min.yml >/dev/null
echo "Running"
ENDPOINT="$gcp" ./tools/wrk -t4 -c10 -d1860s -R10 -s ./benchmark/gctest/workload.lua --timeout 10s "$gcp" >/dev/null
echo "Collecting metrics"
python ./scripts/gctest/gctest.py --command run --config gc1min

echo "Reset Database"
go run ./internal/gctest/init/init.go
echo "Deploying 10 min gc"
sls deploy -c gctest10min.yml >/dev/null
echo "Running"
ENDPOINT="$gcp" ./tools/wrk -t4 -c10 -d1860s -R10 -s ./benchmark/gctest/workload.lua --timeout 10s "$gcp" >/dev/null
echo "Collecting metrics"
python ./scripts/gctest/gctest.py --command run --config gc10min
go run ./internal/gctest/init/init.go clean

make clean
make gctesttxn
echo "Deploying beldi-txn"
sls deploy -c gctestnogc.yml
echo "Reset Database"
go run ./internal/gctest/init/init.go txn
echo "Running"
ENDPOINT="$gcp" ./tools/wrk -t4 -c10 -d1860s -R10 -s ./benchmark/gctest/workload.lua --timeout 10s "$gcp"
echo "Collecting metrics"
python ./scripts/gctest/gctest.py --command run --config txn

python ./scripts/gctest/generate.py
read -p "Choose mode (fast or full, default: full): " mode
mode=${mode:-"full"}
if [ "$mode" == "fast" ]; then
  duration=300
  durationmin=5
else
  duration=1860
  durationmin=30
fi
echo "Compiling"
make clean >/dev/null
make gctest >/dev/null
echo "Deploying no gc"
sls deploy -c gctestnogc.yml >/dev/null
read -p "Please Input HTTP gateway url for beldi-dev-gctest: " gcp
echo "Initializing Database"
go run ./internal/gctest/init/init.go
echo "Running"
ENDPOINT="$gcp" ./tools/wrk -t4 -c10 -d"$duration"s -R10 -s ./benchmark/gctest/workload.lua --timeout 10s "$gcp" >/dev/null
echo "Collecting metrics"
python ./scripts/gctest/gctest.py --command run --config nogc --duration "$durationmin"

echo "Reset Database"
go run ./internal/gctest/init/init.go
echo "Deploying 1 min gc"
sls deploy -c gctest1min.yml >/dev/null
echo "Running"
ENDPOINT="$gcp" ./tools/wrk -t4 -c10 -d"$duration"s -R10 -s ./benchmark/gctest/workload.lua --timeout 10s "$gcp" >/dev/null
echo "Collecting metrics"
python ./scripts/gctest/gctest.py --command run --config gc1min --duration "$durationmin"

echo "Reset Database"
go run ./internal/gctest/init/init.go
echo "Deploying 10 min gc"
sls deploy -c gctest10min.yml >/dev/null
echo "Running"
ENDPOINT="$gcp" ./tools/wrk -t4 -c10 -d"$duration"s -R10 -s ./benchmark/gctest/workload.lua --timeout 10s "$gcp" >/dev/null
echo "Collecting metrics"
python ./scripts/gctest/gctest.py --command run --config gc10min --duration "$durationmin"

make clean >/dev/null
make gctest >/dev/null
echo "Deploying beldi-txn"
sls deploy -c gctestnogc.yml >/dev/null
echo "Reset Database"
go run ./internal/gctest/init/init.go txn
echo "Running"
ENDPOINT="$gcp" ./tools/wrk -t4 -c10 -d"$duration"s -R10 -s ./benchmark/gctest/workload.lua --timeout 10s "$gcp" >/dev/null
echo "Collecting metrics"
python ./scripts/gctest/gctest.py --command run --config txn --duration "$durationmin"
echo "Cleanup"
go run ./internal/gctest/init/init.go clean

python ./scripts/gctest/generate.py

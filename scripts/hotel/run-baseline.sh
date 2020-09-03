read -p "Please Input request rate (default: 100): " rate
rate=${rate:-100}
echo "Compiling"
make clean >/dev/null
make hotel-baseline >/dev/null
echo "Deploying"
sls deploy -c hotel-baseline.yml >/dev/null
read -p "Please Input HTTP gateway url for beldi-dev-bgateway: " bp
echo "Initializing Database"
go run ./internal/hotel/init/init.go clean baseline
go run ./internal/hotel/init/init.go create baseline
go run ./internal/hotel/init/init.go populate baseline
echo "Running baseline"
ENDPOINT="$bp" ./tools/wrk -t4 -c$rate -d420s -R$rate -s ./benchmark/hotel/workload.lua --timeout 10s "$bp" >/dev/null
echo "Collecting metrics"
python ./scripts/hotel/hotel.py --command run --config baseline --duration 7
echo "Cleanup"
go run ./internal/hotel/init/init.go clean baseline

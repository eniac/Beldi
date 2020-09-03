read -p "Please Input request rate (default: 100): " rate
rate=${rate:-100}
echo "Compiling"
make clean >/dev/null
make media >/dev/null
echo "Deploying"
sls deploy -c media.yml >/dev/null
read -p "Please Input HTTP gateway url for beldi-dev-Frontend: " bp
echo "Initializing Database"
go run ./internal/media/init/init.go clean beldi
go run ./internal/media/init/init.go create beldi
go run ./internal/media/init/init.go populate beldi $(pwd)/internal/media/data/compressed.json
echo "Running beldi"
ENDPOINT="$bp" ./tools/wrk -t4 -c"$rate" -d420s -R"$rate" -s ./benchmark/media/workload.lua --timeout 10s "$bp" >/dev/null
echo "Collecting metrics"
python ./scripts/media/media.py --command run --config beldi --duration 7
echo "Cleanup"
go run ./internal/media/init/init.go clean beldi

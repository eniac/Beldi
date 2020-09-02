read -p "Please Input request rate (default: 100)" rate
rate=${rate:-100}
echo "Compiling"
make clean
make media-baseline
echo "Deploying"
sls deploy -c media-baseline.yml
read -p "Please Input HTTP gateway url for beldi-dev-bFrontend: " bp
echo "Initializing Database"
go run ./internal/media/init/init.go clean baseline
sleep 60
go run ./internal/media/init/init.go create baseline
sleep 60
go run ./internal/media/init/init.go populate baseline $(pwd)/internal/media/data/compressed.json
echo "Running baseline"
ENDPOINT="$bp" ./tools/wrk -t4 -c$rate -d420s -R$rate -s ./benchmark/media/workload.lua --timeout 10s "$bp"
echo "Collecting metrics"
python ./scripts/media/media.py --command run --config baseline --duration 7
echo "Cleanup"
go run ./internal/media/init/init.go clean baseline

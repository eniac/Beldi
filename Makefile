.PHONY: build clean deploy

build:
# single operation

singleop:
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI -X github.com/eniac/Beldi/pkg/beldilib.DLOGSIZE=1000" -o bin/singleop/singleop internal/singleop/main/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI -X github.com/eniac/Beldi/pkg/beldilib.DLOGSIZE=1000" -o bin/singleop/nop internal/singleop/nop/nop.go

	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI -X main.TXN=ENABLE -X github.com/eniac/Beldi/pkg/beldilib.DLOGSIZE=1000" -o bin/tsingleop/tsingleop internal/singleop/main/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI -X main.TXN=ENABLE -X github.com/eniac/Beldi/pkg/beldilib.DLOGSIZE=1000" -o bin/tsingleop/tnop internal/singleop/nop/nop.go

	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bsingleop/bsingleop internal/singleop/main/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bsingleop/bnop internal/singleop/nop/nop.go

hotel-baseline:
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/geo internal/hotel/main/handlers/geo/geo.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/profile internal/hotel/main/handlers/profile/profile.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/rate internal/hotel/main/handlers/rate/rate.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/recommendation internal/hotel/main/handlers/recommendation/recommendation.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/user internal/hotel/main/handlers/user/user.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/search internal/hotel/main/handlers/search/search.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/hotel internal/hotel/main/handlers/hotel/hotel.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/flight internal/hotel/main/handlers/flight/flight.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/order internal/hotel/main/handlers/order/order.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/frontend internal/hotel/main/handlers/frontend/frontend.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bhotel/gateway internal/hotel/main/handlers/gateway/gateway.go

hotel:
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/geo internal/hotel/main/handlers/geo/geo.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/profile internal/hotel/main/handlers/profile/profile.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/rate internal/hotel/main/handlers/rate/rate.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/recommendation internal/hotel/main/handlers/recommendation/recommendation.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/user internal/hotel/main/handlers/user/user.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/search internal/hotel/main/handlers/search/search.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/hotel internal/hotel/main/handlers/hotel/hotel.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/flight internal/hotel/main/handlers/flight/flight.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/order internal/hotel/main/handlers/order/order.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/frontend internal/hotel/main/handlers/frontend/frontend.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/gateway internal/hotel/main/handlers/gateway/gateway.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/gc internal/hotel/main/gc/gc.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/hotel/collector internal/hotel/main/collector/collector.go

media-baseline:
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/CastInfo internal/media/core/handlers/castInfo/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/ComposeReview internal/media/core/handlers/composeReview/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/Frontend internal/media/core/handlers/frontend/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/MovieId internal/media/core/handlers/movieId/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/MovieInfo internal/media/core/handlers/movieInfo/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/MovieReview internal/media/core/handlers/movieReview/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/Page internal/media/core/handlers/page/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/Plot internal/media/core/handlers/plot/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/Rating internal/media/core/handlers/rating/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/ReviewStorage internal/media/core/handlers/reviewStorage/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/Text internal/media/core/handlers/text/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/UniqueId internal/media/core/handlers/uniqueId/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/User internal/media/core/handlers/user/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BASELINE" -o bin/bmedia/UserReview internal/media/core/handlers/userReview/main.go

media:
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/CastInfo internal/media/core/handlers/castInfo/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/ComposeReview internal/media/core/handlers/composeReview/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/Frontend internal/media/core/handlers/frontend/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/MovieId internal/media/core/handlers/movieId/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/MovieInfo internal/media/core/handlers/movieInfo/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/MovieReview internal/media/core/handlers/movieReview/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/Page internal/media/core/handlers/page/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/Plot internal/media/core/handlers/plot/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/Rating internal/media/core/handlers/rating/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/ReviewStorage internal/media/core/handlers/reviewStorage/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/Text internal/media/core/handlers/text/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/UniqueId internal/media/core/handlers/uniqueId/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/User internal/media/core/handlers/user/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/UserReview internal/media/core/handlers/userReview/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/gc internal/media/core/gc/gc.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/media/collector internal/media/core/collector/collector.go

gctest:
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/gctest/gctest internal/gctest/core/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI" -o bin/gctest/gc internal/gctest/core/gc/gc.go

gctesttxn:
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI -X github.com/eniac/Beldi/pkg/beldilib.DLOGSIZE=101" -o bin/gctest/gctest internal/gctest/core/main.go
	env GOOS=linux go build -ldflags="-s -w -X github.com/eniac/Beldi/pkg/beldilib.TYPE=BELDI -X github.com/eniac/Beldi/pkg/beldilib.DLOGSIZE=101" -o bin/gctest/gc internal/gctest/core/gc/gc.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

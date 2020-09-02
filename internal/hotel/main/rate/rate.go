package rate

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
	"sort"
)

func GetRates(env *beldilib.Env, req Request) Result {
	var plans RatePlans
	for _, i := range req.HotelIds {
		plan := data.RatePlan{}
		res := beldilib.Read(env, data.Trate(), i)
		err := mapstructure.Decode(res, &plan)
		beldilib.CHECK(err)
		if plan.HotelId != "" {
			plans = append(plans, plan)
		}
	}
	sort.Sort(plans)
	return Result{RatePlans: plans}
}

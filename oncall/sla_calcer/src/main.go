package main

func main() {
	cfg, err := ParseConfig("../configs/config.yml")
	if err != nil {
		panic(err.Error())
	}
	prom := CreateServicer(cfg)
	var influx *InfluxClient
	for {
		influx, err = CreateInfluxClient()
		if err != nil {
			//fmt.Println(err.Error())
		} else {
			break
		}
	}
	slaC := SLACreater{}
	for {
		t, err := prom.GetTimeToResponce()
		if err != nil {
			//fmt.Println(err.Error())
		}
		form, err := FormatTimeResp(t)
		if err != nil {
			//fmt.Println(err.Error())
		}
		sla := slaC.GenMetric(form)
		for key, val := range sla {
			tags := make(map[string]string)
			tags["dur"] = val.name
			vals := make(map[string]interface{})
			vals["actual"] = val.actualValue
			vals["isBad"] = val.isBad
			vals["limit"] = val.limitValue
			influx.WriteToDb(key, tags, vals, "iot", "users_business_events")
		}

	}

}

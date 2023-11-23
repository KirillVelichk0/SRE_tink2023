package main

type SLAMetric struct {
	name        string
	limitValue  float64
	actualValue float64
	isBad       bool
}

type SLACreater struct {
	last FormattedCollectionData
}

func Compare(name string, limit float64, goodCur uint64, goodPrev uint64, allCur uint64, allPrev uint64) SLAMetric {
	var result SLAMetric
	result.name = name
	result.limitValue = limit
	result.actualValue = float64(goodCur-goodPrev) / float64(allCur-allPrev) * 100
	result.isBad = result.actualValue > result.limitValue
	return result
}

func (creater *SLACreater) GenMetric(collection FormattedCollectionData) map[string]SLAMetric {
	result := make(map[string]SLAMetric)
	result["100"] = Compare("100", 90.0, collection.f100.val, creater.last.f100.val,
		collection.fBig.val, creater.last.fBig.val)
	result["500"] = Compare("500", 95.0, collection.f500.val, creater.last.f500.val,
		collection.fBig.val, creater.last.fBig.val)
	result["1000"] = Compare("1000", 99.0, collection.f1000.val, creater.last.f1000.val,
		collection.fBig.val, creater.last.fBig.val)
	creater.last = collection
	return result
}

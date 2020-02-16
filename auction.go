package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	Sites   []site
	Bidders []bidder
}

type site struct {
	Name    string   `json:"name"`
	Bidders []string `json:"bidders"`
	Floor   float64  `json:"floor"`
}

type bidder struct {
	Name       string  `json:"name"`
	Adjustment float64 `json:"adjustment"`
}

type auction struct {
	Site  string   `json:"site"`
	Units []string `json:"units"`
	Bids  []bid    `json:"bids"`
}

type bid struct {
	Bidder        string  `json:"bidder"`
	Unit          string  `json:"unit"`
	Bid           float64 `json:"bid"`
	AdjustedValue float64 `json:"-"`
}

type siteData struct {
	Bidders []string
	Floor   float64
}

func mapBidders(bidders []bidder) map[string]float64 {
	bidderMap := make(map[string]float64)
	for _, bidder := range bidders {
		bidderMap[bidder.Name] = bidder.Adjustment
	}
	return bidderMap
}

func mapSiteInfo(sites []site) map[string]siteData {
	siteMap := make(map[string]siteData)
	for _, site := range sites {
		siteMap[site.Name] = siteData{
			Bidders: site.Bidders,
			Floor:   site.Floor,
		}
	}
	return siteMap
}

func contains(value string, list []string) bool {
	for _, s := range list {
		if value == s {
			return true
		}
	}
	return false
}

func validBidsPerUnit(auction auction, bidderMap map[string]float64, siteData siteData) map[string][]bid {
	mapUnitsToBids := make(map[string][]bid)
	for _, bid := range auction.Bids {
		// Check if unit appears on the site
		if ok := contains(bid.Unit, auction.Units); ok {
			// Check if bidder is known
			if adjustment, ok := bidderMap[bid.Bidder]; ok {
				// Check if bidder is permitted to bid on the site
				if ok := contains(bid.Bidder, siteData.Bidders); ok {
					// Check if adjusted value is greater or equal to floor
					adjustedValue := adjustment*bid.Bid + bid.Bid
					if adjustedValue >= siteData.Floor {
						bid.AdjustedValue = adjustedValue
						mapUnitsToBids[bid.Unit] = append(mapUnitsToBids[bid.Unit], bid)
					}
				}
			}
		}

	}
	return mapUnitsToBids
}

func maxBidsPerSite(mapUnitsToBids map[string][]bid, units []string) []bid {
	var winnerBids []bid
	for _, unit := range units {
		var maxVal float64
		var winnerBid bid
		if validBidsForSite, ok := mapUnitsToBids[unit]; ok {
			for _, validBid := range validBidsForSite {
				if validBid.AdjustedValue >= maxVal {
					maxVal = validBid.AdjustedValue
					winnerBid = validBid
				}
			}
		}
		if maxVal != 0 {
			winnerBids = append(winnerBids, winnerBid)
		}
	}
	return winnerBids
}

func main() {

	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	bidderMap := mapBidders(config.Bidders)
	siteMap := mapSiteInfo(config.Sites)

	input, err := ioutil.ReadFile("input.json")
	if err != nil {
		panic(err)
	}

	var auctions []auction
	err = json.Unmarshal(input, &auctions)
	if err != nil {
		panic(err)
	}

	var solution [][]bid
	for _, auction := range auctions {
		if siteData, ok := siteMap[auction.Site]; ok {
			validBidsPerUnit := validBidsPerUnit(auction, bidderMap, siteData)
			maxBidsPerSite := maxBidsPerSite(validBidsPerUnit, auction.Units)
			solution = append(solution, maxBidsPerSite)
		}
	}

	solBytes, err := json.Marshal(solution)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(solBytes))

}

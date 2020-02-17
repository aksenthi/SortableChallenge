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

func (site *site) canBid(bidder string) bool {
	return contains(bidder, site.Bidders)
}

type bidder struct {
	Name       string  `json:"name"`
	Adjustment float64 `json:"adjustment"`
}

func (bidder *bidder) adjustedBid(bid float64) float64 {
	return bidder.Adjustment*bid + bid
}

type auction struct {
	Site     string   `json:"site"`
	Units    []string `json:"units"`
	Bids     []bid    `json:"bids"`
	SiteInfo *site
}

func (auction *auction) isRegisteredUnit(unit string) bool {
	return contains(unit, auction.Units)
}

type bid struct {
	Bidder        string  `json:"bidder"`
	Unit          string  `json:"unit"`
	Bid           float64 `json:"bid"`
	AdjustedValue float64 `json:"-"`
	BidderInfo    *bidder `json:"-"`
}

type siteData struct {
	Bidders []string
	Floor   float64
}

func mapBidders(bidders []bidder) map[string]bidder {
	bidderMap := make(map[string]bidder)
	for _, bidder := range bidders {
		bidderMap[bidder.Name] = bidder
	}
	return bidderMap
}

func mapSiteInfo(sites []site) map[string]site {
	siteMap := make(map[string]site)
	for _, site := range sites {
		siteMap[site.Name] = site
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

func (bid *bid) isValid(auction *auction) bool {
	// Check if unit appears on the site
	if auction.isRegisteredUnit(bid.Unit) {
		// Check if bidder is permitted to bid on the site
		if auction.SiteInfo.canBid(bid.Bidder) {
			// Check if adjusted value is greater or equal to floor
			bid.AdjustedValue = bid.BidderInfo.adjustedBid(bid.Bid)
			if bid.AdjustedValue >= auction.SiteInfo.Floor {
				return true
			}
		}
	}
	return false
}

func (auction *auction) findWinners(bidderMap map[string]bidder) []bid {
	winnersPerUnit := make(map[string]bid)
	winners := []bid{}
	for _, bid := range auction.Bids {
		// Attach bidder to bid if it exists; if it doesn't exist, ignore the bid
		if bidder, ok := bidderMap[bid.Bidder]; ok {
			bid.BidderInfo = &bidder
			// Check if bid is valid
			if bid.isValid(auction) {
				currentWinner := winnersPerUnit[bid.Unit]
				if bid.AdjustedValue > currentWinner.Bid {
					winnersPerUnit[bid.Unit] = bid
				}
			}
		}
	}
	for _, winningBid := range winnersPerUnit {
		if winningBid.Bid > 0 {
			winners = append(winners, winningBid)
		}
	}
	return winners
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
			auction.SiteInfo = &siteData
			winners := auction.findWinners(bidderMap)
			solution = append(solution, winners)
		} else {
			solution = append(solution, []bid{})
		}
	}

	solBytes, err := json.Marshal(solution)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(solBytes))

}

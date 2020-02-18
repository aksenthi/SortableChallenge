package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	SiteName string   `json:"site"`
	Units    []string `json:"units"`
	Bids     []bid    `json:"bids"`
	Site     *site
}

func (auction *auction) isRegisteredUnit(unit string) bool {
	return contains(unit, auction.Units)
}

type bid struct {
	Bidder string  `json:"bidder"`
	Unit   string  `json:"unit"`
	Bid    float64 `json:"bid"`
}

func createBidderIndex(bidders []bidder) map[string]*bidder {
	bidderIndex := make(map[string]*bidder)
	for i, bidder := range bidders {
		bidderIndex[bidder.Name] = &bidders[i]

	}
	return bidderIndex
}

func createSiteIndex(sites []site) map[string]*site {
	siteIndex := make(map[string]*site)
	for i, site := range sites {
		siteIndex[site.Name] = &sites[i]
	}
	return siteIndex
}

func contains(value string, list []string) bool {
	for _, s := range list {
		if value == s {
			return true
		}
	}
	return false
}

func (auction *auction) isValid(bid *bid) bool {
	// Check if unit appears on the site and if bidder is permitted to bid on the site
	return auction.isRegisteredUnit(bid.Unit) && auction.Site.canBid(bid.Bidder)
}

func (auction *auction) findWinners(bidderIndex map[string]*bidder) []bid {
	winnersPerUnit := make(map[string]bid)
	winningBidsPerUnit := make(map[string]float64)
	for _, unit := range auction.Units {
		winningBidsPerUnit[unit] = auction.Site.Floor
	}
	winners := []bid{}
	for _, bid := range auction.Bids {
		if bidder, ok := bidderIndex[bid.Bidder]; ok {
			// Check if bid is valid
			if auction.isValid(&bid) {
				adjustedValue := bidder.adjustedBid(bid.Bid)
				if adjustedValue >= winningBidsPerUnit[bid.Unit] {
					winningBidsPerUnit[bid.Unit] = adjustedValue
					winnersPerUnit[bid.Unit] = bid
				}
			}
		}
	}
	for _, winningBid := range winnersPerUnit {
		winners = append(winners, winningBid)
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

	bidderIndex := createBidderIndex(config.Bidders)
	siteIndex := createSiteIndex(config.Sites)

	var auctions []auction

	err = json.NewDecoder(os.Stdin).Decode(&auctions)
	if err != nil {
		panic(err)
	}

	var solution [][]bid
	for _, auction := range auctions {
		if siteData, ok := siteIndex[auction.SiteName]; ok {
			auction.Site = siteData
			winners := auction.findWinners(bidderIndex)
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

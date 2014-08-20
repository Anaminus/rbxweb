// Deals with services related to ROBLOX currency.
package currency

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anaminus/rbxweb"
	"io"
	"net/url"
	"os"
	"time"
)

func TradeTickets(client *rbxweb.Client, tickets int64, robux int64, limit bool, split bool) (err error) {
	page := client.GetURL(`www`, `/My/Money.aspx`, nil)
	query := url.Values{
		"__EVENTTARGET":                                                           {"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$SubmitTradeButton"},
		"__VIEWSTATE":                                                             {},
		"__EVENTVALIDATION":                                                       {},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$HaveCurrencyDropDownList": {"Tickets"},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$HaveAmountTextBoxRestyle": {strconv.Itoa(tickets)},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$WantCurrencyDropDownList": {"Robux"},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$WantAmountTextBox":        {strconv.Itoa(robux)},
	}
	if limit {
		query.Set("ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$OrderType", "LimitOrderRadioButton")
	} else {
		query.Set("ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$OrderType", "MarketOrderRadioButton")
	}
	if split {
		query.Set("ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$AllowSplitTradesCheckBox", "on")
	} else {
		query.Set("ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$AllowSplitTradesCheckBox", "off")
	}
	err = client.DoRawPost(page, query)
	return
}

func TradeRobux(client *rbxweb.Client, robux int64, tickets int64, limit bool, split bool) (err error) {
	page := client.GetURL(`www`, `/My/Money.aspx`, nil)
	query := url.Values{
		"__EVENTTARGET":                                                           {"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$SubmitTradeButton"},
		"__VIEWSTATE":                                                             {},
		"__EVENTVALIDATION":                                                       {},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$HaveCurrencyDropDownList": {"Robux"},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$HaveAmountTextBoxRestyle": {strconv.Itoa(robux)},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$WantCurrencyDropDownList": {"Tickets"},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$WantAmountTextBox":        {strconv.Itoa(tickets)},
	}
	if limit {
		query.Set("ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$OrderType", "LimitOrderRadioButton")
	} else {
		query.Set("ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$OrderType", "MarketOrderRadioButton")
	}
	if split {
		query.Set("ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$AllowSplitTradesCheckBox", "on")
	} else {
		query.Set("ctl00$ctl00$cphRoblox$cphMyRobloxContent$ctl00$AllowSplitTradesCheckBox", "off")
	}
	err = client.DoRawPost(page, query)
	return
}

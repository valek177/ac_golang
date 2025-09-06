package main

import urllib "net/url"

func isUrlValid(url string) bool {
	u, err := urllib.Parse(url)
	return err == nil && u.Scheme != "" && u.Host != ""
}

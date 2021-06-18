# GuideStar Scraper
This is a fairly simple go program that scrapes GuideStar as well as some other things.

This scraper does not require authentication tokens or keys.

Use this program **at your own risk**. At the end of the day, it is scraping and scraping can very easily get you banned off the site.

## Setup
In order for this script to work make sure you have the latest version of go and you run the `setup.sh` script. (If you're on a Windows machine run the 2 commands in the file manually. They simply install 2 packages which make the scripts run).

Make sure you run `go build` in the root folder so that the executables will be able to run on your computer. The names of the executable is `scraper`

Following this, place all the nonprofit names you want to scrape into `npnames` (Delimiter is newline). Using these names it will output a json file for each nonprofit in `./nonprofit-info`. Each JSON file will contain the name, EIN, description, and programs of the particular nonprofit.

If for some reason it couldn't find the nonprofit. It will put that nonprofit name in the `missing-eins` file.

## Excel
Once you run `scraper` you can run `json-to-excel` in the `json-to-excel` to get the information as a large excel file called `nonprofits.xslx`. Make sure only to run this after you run `scraper` as it uses the json files in `nonprofit-info`.

In order to run this executable. Go into `json-to-excel`, run `go build` to get an up-to-date executable, and then simply run `json-to-excel`.

## How It Works
This script works by taking the nonprofit name and sending it to a ProPublica API to get the EIN. Using this EIN we are able to bypass authentication on GuideStar entirely and scrape the site.

## Additional Comments
Personally, I prefer to run this script 100 nonprofits at a time in order to minimize my chances of getting blocked. But, that is up to you.

In my experience, this scraper has about a 50% hitrate. This is mainly due to the fact that ProPublica sometimes can't find a nonprofit. 

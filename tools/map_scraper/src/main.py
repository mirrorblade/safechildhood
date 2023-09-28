from argparse import ArgumentParser

from src.yandex_map_scraper import YandexMapScraper

from selenium import webdriver

def main():
    parser = ArgumentParser()
    parser.add_argument("-o", "--output", type=str, help="path to output csv file (example: \"./output.csv\")", required=True)
    parser.add_argument("-ec", "--excludedCategories", type=str, nargs="+", help="list with excluded categories (example: \"игровая комната\" \"игровая площадка\")", required=False, default=[])
    parser.add_argument("-u", "--url", type=str, help="url to yandex map with search results", required=True)

    args = parser.parse_args()

    scraper = YandexMapScraper(webdriver.Firefox())
    scraper.parse_search_results(args.output, args.excludedCategories, args.url)

if __name__ == "__main__":
    main()
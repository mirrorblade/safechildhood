from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions
from selenium.common.exceptions import *

from nanoid import generate

class YandexMapScraper():
    def __init__(self, driver):
        self.__driver = driver

    def parse_search_results(self, csvOutputPath, excludedCategories, url):
        self.__driver.get(url)

        self.__scroll_to_end(By.XPATH, "//div[@class='add-business-view']")

        search_results = self.__driver.find_element(By.XPATH, "//ul[@class='search-list-view__list']").find_elements(By.TAG_NAME, "li")

        categoryCheckExists = len(excludedCategories) > 0
        if categoryCheckExists:
            excludedCategories = list(map(lambda x: x.lower(), excludedCategories))

        print(f"search results length: {len(search_results) - 1}")

        count = 0
        
        with open(csvOutputPath, 'w', encoding='UTF-8') as file:
            file.write('COORDINATES;ADDRESS;ID\n')

            for result in search_results[1:]:
                coordinates = result\
                        .find_element(By.XPATH, ".//div[contains(concat(' ', normalize-space(@class), ' '), ' search-snippet-view__body _type_business ')]")\
                        .get_attribute('data-coordinates').split(',')

                address = result.find_element(By.XPATH, ".//*[@class='search-business-snippet-view__address']").text
            
                categories_list = result.find_elements(By.XPATH, ".//a[@class='search-business-snippet-view__category']")
                text_list = list(map(lambda x: x.text.lower(), categories_list))

                flag = False

                if categoryCheckExists:
                    for category in excludedCategories:
                        if category in text_list:
                            flag = True

                            break

                if flag:
                    continue
                    
                count += 1

                file.write(f'{coordinates[1]},{coordinates[0]};{address};{generate()}\n') 

        print(f"search results with excluded categories length: {count}")   

        print("file uploaded!")

        self.__driver.quit()

    def __scroll_to_end(self, by_end_element, value_end_element, x_offset=384, y_offset=200): #384 - это ширина на которой расположен скролл блока =). Изначально даётся пролистывание скроллбара на 200 пикселей вниз, но если возникает ошибка о том, что значение вылезает за рамки разрешения, изменяем пролистывание на /2
        chain = ActionChains(self.__driver)

        scroll_bar = WebDriverWait(self.__driver, 3000, ignored_exceptions=(MoveTargetOutOfBoundsException, NoSuchElementException,StaleElementReferenceException))\
            .until(expected_conditions.visibility_of_element_located((By.XPATH, "//div[@class='scroll__scrollbar-thumb']")))
        
        new_y_offset = y_offset

        while True:
            try:
                chain.drag_and_drop_by_offset(scroll_bar, x_offset, new_y_offset).perform()

                self.__driver.find_element(by_end_element, value_end_element) #Элемент появляется тогда, когда все элементы поиска выгрузились на сайт. Если его нет, вылезает ошибка.  

            except MoveTargetOutOfBoundsException:
                new_y_offset /= 2

            except NoSuchElementException:
                new_y_offset = y_offset
            
            except Exception as e:
                raise e
            
            else:
                print("nice scrolling!")

                break



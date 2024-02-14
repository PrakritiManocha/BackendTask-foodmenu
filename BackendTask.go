package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "github.com/tealeg/xlsx"
)

type Meal struct {
    Day   string
    Date  string
    Meal  string
    Items []string
}

func readMenuFromExcel(menuXlsx string) (map[string]map[string][]string, error) {
    menu := make(map[string]map[string][]string)
    xlFile, err := xlsx.OpenFile(menuXlsx)
    if err != nil {
        return nil, err
    }
    for _, sheet := range xlFile.Sheets {
        dayMenu := make(map[string][]string)
        for _, row := range sheet.Rows {
            day := row.Cells[0].String()
            meal := row.Cells[2].String()
            items := row.Cells[3:]
            var itemList []string
            for _, item := range items {
                itemList = append(itemList, item.String())
            }
            dayMenu[meal] = itemList
        }
        menu[sheet.Name] = dayMenu
    }
    return menu, nil
}

func getItemsForMeal(menu map[string]map[string][]string, day, meal string) []string {
    return menu[day][meal]
}

func countItemsForMeal(menu map[string]map[string][]string, day, meal string) int {
    return len(menu[day][meal])
}

func isItemAvailable(menu map[string]map[string][]string, day, meal, item string) bool {
    items := menu[day][meal]
    for _, i := range items {
        if i == item {
            return true
        }
    }
    return false
}

func saveMenuAsJSON(menu map[string]map[string][]string) error {
    file, err := os.Create("menu.json")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    err = encoder.Encode(menu)
    if err != nil {
        return err
    }
    fmt.Println("Menu successfully saved as menu.json")
    return nil
}

func (m Meal) PrintDetails() {
    fmt.Printf("Day: %s\nDate: %s\nMeal: %s\nItems: %v\n", m.Day, m.Date, m.Meal, m.Items)
}

func main() {
    menu, err := readMenuFromExcel("menu.xlsx")
    if err != nil {
        log.Fatal(err)
    }

    day := "MONDAY"
    meal := "BREAKFAST"

    items := getItemsForMeal(menu, day, meal)
    fmt.Printf("Items available for %s on %s: %v\n", meal, day, items)

    numItems := countItemsForMeal(menu, day, meal)
    fmt.Printf("Number of items for %s on %s: %d\n", meal, day, numItems)

    item := "POHA"
    available := isItemAvailable(menu, day, meal, item)
    fmt.Printf("Is %s available for %s on %s: %t\n", item, meal, day, available)

    err = saveMenuAsJSON(menu)
    if err != nil {
        log.Fatal(err)
    }

    firstMeal := Meal{
        Day:   "MONDAY",
        Date:  "05-Feb-24",
        Meal:  "BREAKFAST",
        Items: menu["MONDAY"]["BREAKFAST"],
    }
    firstMeal.PrintDetails()
}

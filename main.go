package main

import (
	"encoding/json"
	"fmt"           
	"log"           
	"os"            
	"strconv"       
	"github.com/rivo/tview"
)

type Item struct {
	Name  string `json:"name"`  
	Stock int    `json:"stock"` 
}

var (
	inventory     = []Item{}         
	inventoryFile = "inventory.json" 
)

func loadInventory() {
	if _, err := os.Stat(inventoryFile); err == nil {
		data, err := os.ReadFile(inventoryFile)
		if err != nil {
			log.Fatal("Error reading inventory file:", err)
		}
		json.Unmarshal(data, &inventory)
	}
}

func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", "  ")
	if err != nil {
		log.Fatal("Error saving inventory:", err)
	}
	os.WriteFile(inventoryFile, data, 0644)
}

func deleteItem(index int) {
	if index < 0 || index >= len(inventory) {
		fmt.Println("Invalid item index.")
		return
	}
	inventory = append(inventory[:index], inventory[index+1:]...)
	saveInventory()
}

func modifyStock(index int, newStock int) {
    if index < 0 || index >= len(inventory) {
        fmt.Println("Invalid item index.")
        return
    }
    inventory[index].Stock = newStock
    saveInventory()
}

func main() {
	app := tview.NewApplication()

	loadInventory()

	inventoryList := tview.NewTextView().
		SetDynamicColors(true). 
		SetRegions(true).       
		SetWordWrap(true)       

	inventoryList.SetBorder(true).SetTitle("Inventory") 

	refreshInventory := func() {
		inventoryList.Clear()
		if len(inventory) == 0 {
			fmt.Fprintln(inventoryList, "No items in inventory.")
		} else {
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Stock: %d)\n", i+1, item.Name, item.Stock)
			}
		}
	}

	itemNameInput := tview.NewInputField().SetLabel("Item Name: ")
	itemStockInput := tview.NewInputField().SetLabel("Stock: ")
	itemIDInput := tview.NewInputField().SetLabel("Item ID to delete: ")

	form := tview.NewForm().
		AddFormItem(itemNameInput).    
		AddFormItem(itemStockInput).   
		AddFormItem(itemIDInput).      
		AddButton("Add Item", func() { 
			name := itemNameInput.GetText()
			stock := itemStockInput.GetText()
			if name != "" && stock != "" {
				quantity, err := strconv.Atoi(stock)
				if err != nil {
					fmt.Fprintln(inventoryList, "Invalid stock value.")
					return
				}
				inventory = append(inventory, Item{Name: name, Stock: quantity})
				saveInventory()
				refreshInventory()
				itemNameInput.SetText("")
				itemStockInput.SetText("")
			}
		}).
		AddButton("Delete Item", func() { 
			idStr := itemIDInput.GetText()
			if idStr == "" {
				fmt.Fprintln(inventoryList, "Please enter an item ID to delete.")
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 || id > len(inventory) {
				fmt.Fprintln(inventoryList, "Invalid item ID.")
				return
			}
			deleteItem(id - 1)
			fmt.Fprintf(inventoryList, "Item [%d] deleted.\n", id)
			refreshInventory()
			itemIDInput.SetText("") 
		}).
		AddButton("Modify stock", func(){
			idStr := itemIDInput.GetText()
			stockStr := itemStockInput.GetText()
			if idStr == "" {
				fmt.Fprintln(inventoryList, "Please enter an item ID to delete.")
				return
			}
			if stockStr == "" {
				fmt.Fprintln(inventoryList, "Please enter a new stock number")
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 || id > len(inventory) {
				fmt.Fprintln(inventoryList, "Invalid item ID.")
				return
			}
			newStock, err := strconv.Atoi(stockStr)
			if err != nil {
				fmt.Fprintln(inventoryList, "Error parsing stock quantity.")
				return
			}
			if newStock == 0 {
				deleteItem(id - 1)
			}
			modifyStock(id - 1, newStock)	
			fmt.Fprintf(inventoryList, "Item [%d] stock updated.\n", id)
			refreshInventory()
			itemIDInput.SetText("") 
			itemStockInput.SetText("") 
		}).
		AddButton("Exit", func() { 
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(inventoryList, 0, 1, false). 
		AddItem(form, 0, 1, true)            

	manageButton := tview.NewButton("Go to Inventory").
		SetSelectedFunc(func() {
			app.SetRoot(flex, true)
		})		
	manageButton.SetBorder(true)
	exitButton := tview.NewButton("Exit").
		SetSelectedFunc(func() {
			app.Stop()
		})		
	exitButton.SetBorder(true)

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(manageButton, 3, 0, true).
		AddItem(exitButton, 3, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(nil, 0, 1, false)

	menuFlex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(buttonFlex, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(nil, 0, 1, false)

	refreshInventory()

	if err := app.SetRoot(menuFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

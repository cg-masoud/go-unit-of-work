package models

type Order struct {
    ID    int    `json:"id"`
    Items []Item `json:"items"`
}

type Item struct {
    ID        int `json:"id"`
    OrderID   int `json:"order_id"`
    ProductID int `json:"product_id"`
    Quantity  int `json:"quantity"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password     string         `json:"-" gorm:"size:255;not null"`
	Nickname     string         `json:"nickname" gorm:"size:50"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	Phone        string         `json:"phone" gorm:"size:20"`
	Email        string         `json:"email" gorm:"size:100"`
	Balance      float64        `json:"balance" gorm:"type:decimal(10,2);default:0"`
	Integral     int            `json:"integral" gorm:"default:0"`
	Status       int            `json:"status" gorm:"default:1"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Admin struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Role      string         `json:"role" gorm:"size:20;default:'admin'"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:50;not null"`
	Icon      string         `json:"icon" gorm:"size:255"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Products  []Product      `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  uint           `json:"category_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Image       string         `json:"image" gorm:"size:255"`
	Images      string         `json:"images" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	MarketPrice float64       `json:"market_price" gorm:"type:decimal(10,2)"`
	Stock       int            `json:"stock" gorm:"default:0"`
	Sales       int            `json:"sales" gorm:"default:0"`
	Unit        string         `json:"unit" gorm:"size:20"`
	Description string         `json:"description" gorm:"type:text"`
	Status      int            `json:"status" gorm:"default:1"`
	IsRecommend int            `json:"is_recommend" gorm:"default:0"`
	IsNew       int            `json:"is_new" gorm:"default:0"`
	Integral    int            `json:"integral" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Category    *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type ShoppingCart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	ProductID uint      `json:"product_id" gorm:"not null;index"`
	Quantity  int       `json:"quantity" gorm:"default:1"`
	Selected  int       `json:"selected" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Product   *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type Address struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	Name       string         `json:"name" gorm:"size:50;not null"`
	Phone      string         `json:"phone" gorm:"size:20;not null"`
	Province   string         `json:"province" gorm:"size:50"`
	City       string         `json:"city" gorm:"size:50"`
	District   string         `json:"district" gorm:"size:50"`
	Address    string         `json:"address" gorm:"size:255;not null"`
	IsDefault  int            `json:"is_default" gorm:"default:0"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type Order struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	OrderNo       string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	UserID        uint           `json:"user_id" gorm:"not null;index"`
	AddressID     uint           `json:"address_id"`
	TotalPrice    float64        `json:"total_price" gorm:"type:decimal(10,2);not null"`
	ActualPrice   float64        `json:"actual_price" gorm:"type:decimal(10,2)"`
	Freight       float64        `json:"freight" gorm:"type:decimal(10,2);default:0"`
	Status        int            `json:"status" gorm:"default:0"`
	PayType       string         `json:"pay_type" gorm:"size:20"`
	PayTime       *time.Time     `json:"pay_time"`
	Remark        string         `json:"remark" gorm:"size:255"`
	CourierNo     string         `json:"courier_no" gorm:"size:50"`
	ShippingTime  *time.Time     `json:"shipping_time"`
	ReceivedTime  *time.Time     `json:"received_time"`
	CompleteTime  *time.Time     `json:"complete_time"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	Items         []OrderItem    `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	Address       *Address       `json:"address,omitempty" gorm:"foreignKey:AddressID"`
}

type OrderItem struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	OrderID     uint     `json:"order_id" gorm:"not null;index"`
	ProductID   uint     `json:"product_id" gorm:"not null;index"`
	ProductName string   `json:"product_name" gorm:"size:100;not null"`
	ProductImage string  `json:"product_image" gorm:"size:255"`
	Price       float64  `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity    int      `json:"quantity" gorm:"default:1"`
	TotalPrice  float64  `json:"total_price" gorm:"type:decimal(10,2);not null"`
	Commented   int      `json:"commented" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Favorite struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;uniqueIndex:idx_user_product"`
	ProductID uint           `json:"product_id" gorm:"not null;uniqueIndex:idx_user_product;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Product   *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type Comment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	ProductID   uint           `json:"product_id" gorm:"not null;index"`
	OrderID     uint           `json:"order_id" gorm:"index"`
	OrderItemID uint          `json:"order_item_id" gorm:"index"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Images      string         `json:"images" gorm:"type:text"`
	Rating      int            `json:"rating" gorm:"default:5"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type Knowledge struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Cover     string         `json:"cover" gorm:"size:255"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Author    string         `json:"author" gorm:"size:50"`
	Views     int            `json:"views" gorm:"default:0"`
	Likes     int            `json:"likes" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Community struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Images    string         `json:"images" gorm:"type:text"`
	Likes     int            `json:"likes" gorm:"default:0"`
	Comments  int            `json:"comments" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type Banner struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:100"`
	Image     string         `json:"image" gorm:"size:255;not null"`
	Link      string         `json:"link" gorm:"size:255"`
	Type      int            `json:"type" gorm:"default:1"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Promotion struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Type        int            `json:"type" gorm:"default:1"`
	Discount    float64        `json:"discount" gorm:"type:decimal(10,2)"`
	MinAmount   float64        `json:"min_amount" gorm:"type:decimal(10,2);default:0"`
	ProductIDs  string         `json:"product_ids" gorm:"type:text"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     time.Time      `json:"end_time"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Recharge struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	OrderNo   string         `json:"order_no" gorm:"uniqueIndex;size:50;not null"`
	Amount    float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	Bonus     float64        `json:"bonus" gorm:"type:decimal(10,2);default:0"`
	Status    int            `json:"status" gorm:"default:0"`
	PayTime   *time.Time     `json:"pay_time"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Integral struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Type      int       `json:"type" gorm:"not null"`
	Integral  int       `json:"integral" gorm:"not null"`
	Source    string    `json:"source" gorm:"size:100"`
	Remark    string    `json:"remark" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
}

type Forum struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Images    string         `json:"images" gorm:"type:text"`
	Views     int            `json:"views" gorm:"default:0"`
	Likes     int            `json:"likes" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type News struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Cover     string         `json:"cover" gorm:"size:255"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	Author    string         `json:"author" gorm:"size:50"`
	Views     int            `json:"views" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

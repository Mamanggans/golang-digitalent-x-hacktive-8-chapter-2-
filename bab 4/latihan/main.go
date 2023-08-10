package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	user     = "postgres"
	password = "jakarta2017"
	dbname   = "challengeEight"
	port     = "5432"
)
const PORT = ":3000"

// type Model struct {
// 	ID        uint `gorm:"primaryKey"`
// 	CreatedAt time.Time
// 	updated   time.Time
// 	DeletedAt time.Time `gorm:"index"`
// }

type User struct {
	// Model
	gorm.Model
	// ID        uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null;type:varchar(191)"`
	Books []Book
	// CreatedAt time.Time
	// updated   time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Printf("before create user: %+v \n", *u)
	err = nil
	return
}

type Book struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"unique;not null;type:varchar(191) json:"title"`
	Author    string    `gorm:"not null;type:varchar(191) json:"author"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	updated   time.Time `json:"updated_at"`
}

type BookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	// UserID uint   `json:"userid"`
}

var (
	db  *gorm.DB
	err error
)

// --------------------------- dijalankan sebelum main  --------------------------------\\
func init() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
	// --------------------------- for create new table --------------------------------\\

	// err = db.Debug().AutoMigrate(User{}, Book{})

	// if err != nil {
	// 	panic(err.Error())
	// }
	fmt.Println("you have successfully connect to database")

}

// --------------------------- main --------------------------------\\
func main() {

	route := gin.Default()

	bookRoutes := route.Group("/books")
	{
		bookRoutes.POST("/", createBook)
		bookRoutes.GET("/", getbooks)
		bookRoutes.GET("/:booksId", getBooksbyId)
		bookRoutes.PUT("/:booksId", updateBook)
		bookRoutes.DELETE("/:booksId", deleteBooks)
	}

	route.Run(PORT)
	// -----------------------------------------------------------------------------------\\

	// createUser("orang2")
	// updateUser(1, "orang1")
	// saveUser(1, "orang orangan")
}

// --------------------------- for create user --------------------------------\\

func createUser(name string) {
	user := User{
		Name: name,
	}

	err = db.Create(&user).Error

	if err != nil {
		log.Panicf("error creating user %s:", err.Error())
	}
	fmt.Printf("new user has been created with detail %+v\n", user)
}

// --------------------------- for updated user --------------------------------\\

func updateUser(userId uint, newName string) {
	user := User{}

	user.ID = userId

	err = db.Model(&user).Updates(User{Name: newName}).Error

	if err != nil {
		log.Panicf("error updating user: %s \n", err.Error())
	}

	fmt.Printf("updated user %+v \n", user)

}

// --------------------------- for save user--------------------------------\\
func saveUser(userId uint, newName string) {
	user := User{}

	user.ID = userId

	user.Name = newName

	err = db.Save(&user).Error

	if err != nil {
		log.Panicf("error saving user: %s \n", err.Error())
	}
}

// --------------------------- for create books--------------------------------\\

func createBook(c *gin.Context) {
	var bookRequest BookRequest
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(422, gin.H{
			"message": "invalid request body",
		})
		return
	}

	book := Book{
		Title:  bookRequest.Title,
		Author: bookRequest.Author,
		UserID: 1,
		// uint(bookRequest.UserID),
	}

	err = db.Create(&book).Error

	if err != nil {
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(201, book)
}

// --------------------------- for get books by id--------------------------------\\

func getBooksbyId(c *gin.Context) {
	bookId := c.Param("booksId")

	parseId, err := strconv.Atoi(bookId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid book id, please check again your book id",
		})
		return
	}

	book := Book{
		ID: uint(parseId),
	}

	err = db.First(&book).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "books not found, please check again your book id",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, book)
}

// --------------------------- for get books--------------------------------\\

func getbooks(c *gin.Context) {
	books := []Book{}

	err = db.Find(&books).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

// --------------------------- for update books--------------------------------\\

func updateBook(c *gin.Context) {
	booksId := c.Param("booksId")

	parseId, err := strconv.Atoi(booksId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid book id please check your book id",
		})
		return
	}
	var bookRequest BookRequest

	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(422, gin.H{
			"message": "invalid request body",
		})
		return
	}

	book := Book{
		ID: uint(parseId),
	}

	err = db.Model(&book).Updates(Book{Title: bookRequest.Title, Author: bookRequest.Author}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong please contact the admin ",
		})
		return
	}

	c.JSON(http.StatusOK, book)
}

func deleteBooks(c *gin.Context) {
	bookId := c.Param("booksId")

	parseId, err := strconv.Atoi(bookId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "your books id is invalid please check again your id",
		})
		return
	}

	err = db.Delete(&Book{}, parseId).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong please contact the admin",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "book has been deleted successfully",
	})
}

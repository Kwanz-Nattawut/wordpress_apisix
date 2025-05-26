package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type Record struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Value string `db:"value" json:"value"`
}

// Middleware สำหรับตรวจสอบ API Key
func apiKeyAuth(c *fiber.Ctx) error {
	key := c.Get("apikey")
	API_KEY := os.Getenv("APISIX_KEY")
	if key != API_KEY {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid API key",
		})
	}

	return c.Next() // ให้ผ่านไปยัง handler ถ้าถูกต้อง
}
func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Local",
		dbUser, dbPassword, dbHost, dbName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to MariaDB!")

	app := fiber.New()
	app.Use(apiKeyAuth)
	app.Get("/data", func(c *fiber.Ctx) error {
		var records []Record

		err := db.Select(&records, "SELECT id,name,value FROM records")
		if err != nil {
			return c.Status(500).SendString("DB error: " + err.Error())
		}

		return c.JSON(records)
	})
	app.Post("/data", func(c *fiber.Ctx) error {
		// สร้าง struct ตัวใหม่เพื่อรับข้อมูลจาก JSON body
		var newRecord Record

		// Parse JSON body ลงใน struct
		if err := c.BodyParser(&newRecord); err != nil {
			return c.Status(400).SendString("Invalid request body: " + err.Error())
		}

		// ตรวจสอบข้อมูล (optional)
		if newRecord.Name == "" || newRecord.Value == "" {
			return c.Status(400).SendString("Name and Value are required")
		}

		// ทำการ Insert ข้อมูลลงฐานข้อมูล
		result, err := db.Exec("INSERT INTO records (name, value) VALUES (?, ?)", newRecord.Name, newRecord.Value)
		if err != nil {
			return c.Status(500).SendString("DB insert error: " + err.Error())
		}
		// ดึง id ที่เพิ่มล่าสุด (optional)
		id, err := result.LastInsertId()
		if err != nil {
			log.Println("Warning: cannot get last insert ID:", err)
		} else {
			newRecord.ID = int(id)
		}

		// ส่งกลับข้อมูล record ที่เพิ่มแล้วในรูปแบบ JSON
		//return c.Status(201).JSON(newRecord)
		return c.Status(201).JSON("new record successfully")

	})
	app.Put("/data/:id", func(c *fiber.Ctx) error {
		// ดึง ID จาก URL param
		id := c.Params("id")

		// รับข้อมูลจาก body
		var updatedRecord Record
		if err := c.BodyParser(&updatedRecord); err != nil {
			return c.Status(400).SendString("Invalid request body: " + err.Error())
		}

		// ตรวจสอบข้อมูล (optional)
		if updatedRecord.Name == "" || updatedRecord.Value == "" {
			return c.Status(400).SendString("Name and Value are required")
		}

		// ทำการอัปเดต record ที่ตรงกับ ID
		result, err := db.Exec("UPDATE records SET name = ?, value = ? WHERE id = ?", updatedRecord.Name, updatedRecord.Value, id)
		if err != nil {
			return c.Status(500).SendString("DB update error: " + err.Error())
		}

		// ตรวจสอบว่า row ถูกอัปเดตจริงไหม
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return c.Status(500).SendString("Cannot get affected rows: " + err.Error())
		}
		if rowsAffected == 0 {
			return c.Status(404).SendString("Record not found")
		}

		// ส่ง response กลับ
		return c.Status(200).JSON(fiber.Map{
			"message": "Record updated successfully",
			"id":      id,
		})
	})
	app.Delete("/data/:id", func(c *fiber.Ctx) error {
		// ดึง id จาก URL param
		id := c.Params("id")

		// ทำการลบ record ออกจากฐานข้อมูล
		result, err := db.Exec("DELETE FROM records WHERE id = ?", id)
		if err != nil {
			return c.Status(500).SendString("DB delete error: " + err.Error())
		}

		// ตรวจสอบว่า row ถูกลบจริงไหม
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return c.Status(500).SendString("Cannot get affected rows: " + err.Error())
		}
		if rowsAffected == 0 {
			return c.Status(404).SendString("Record not found")
		}

		return c.Status(200).SendString("Record deleted successfully")
	})
	log.Fatal(app.Listen(":3000"))
}

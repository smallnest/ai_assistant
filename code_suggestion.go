package aiassistant

import (
	"database/sql"

	"gorm.io/gorm"
)

// 生成一个Student struct
// The Student struct is used to store information about a student.
type Student struct {
	Name  string
	Age   int
	Class string
}

// CreateStudent creates a student object
func CreateStudent(name string, age int, class string) Student {
	return Student{
		Name:  name,
		Age:   age,
		Class: class,
	}
}

// 从数据库中查询所有的年龄大于20的学生
func QueryStudentsFromDB(db *sql.DB) []*Student {
	var sudtents []*Student
	rows, err := db.Query("select * from student where age > ?", 20)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var age int
		var class string
		err = rows.Scan(&name, &age, &class)
		if err != nil {
			panic(err)
		}
		student := CreateStudent(name, age, class)
		// do something with student
		sudtents = append(sudtents, &student)
	}

	return sudtents
}

// 使用gorm方式查询所有的年龄大于20的学生
// QueryStudentsFromDBWithGorm queries students from database using GORM.
func QueryStudentsFromDBWithGorm(db *gorm.DB) ([]*Student, error) {
	var students []*Student
	if err := db.Where("age > ?", 20).Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

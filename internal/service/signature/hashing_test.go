package signature

//
//import (
//	"fmt"
//	"helper/internal/data"
//	"log"
//	"testing"
//)
//
//func TestHash(t *testing.T) {
//
//	users := []data.User{
//		data.User{
//			Date:        "20.12.2022",
//			Participant: "Olena Sporova",
//			CourseTitle: "Theory of database organization and basic SQL",
//		},
//
//		data.User{
//			Date:        "10.05.2016",
//			Participant: "Nikita Magda",
//			CourseTitle: "Cryptocurrencies and Distributed Systems",
//		},
//	}
//
//	results := "d109143293b242e776ae2050f4e437347ee566b6d1f4539ade9f17c7b60be4ab"
//
//	for _, user := range users {
//		got := fmt.Sprintf("%x", Hashing(&user))
//		want := results
//
//		if got != want {
//			t.Errorf("got %q, wanted %q", got, want)
//			continue
//		}
//		log.Println("PASS")
//
//	}
//
//}

package worker

import (
	"fmt"
	"time"

	"github.com/guidogimeno/smartpay-be/db"
	"github.com/guidogimeno/smartpay-be/scrapper"
	"github.com/guidogimeno/smartpay-be/utils"
)

type Worker struct {
	db db.Storer
}

func New(db db.Storer) *Worker {
	return &Worker{
		db: db,
	}
}

func (w *Worker) Start() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		startDate1 := utils.NewClock().AddMonths(-2).Format()
		finishDate1 := utils.NewClock().AddMonths(-1).Format()
		inflationIndexes, err := scrapper.ScrapInflation(startDate1, finishDate1)
		if err != nil {
			fmt.Println(err)
		}
		w.db.Create(inflationIndexes[0])

		startDate2 := utils.NewClock().AddDays(-20).Format()
		finishDate2 := utils.NewClock().AddDays(-7).Format()
		tnas, err := scrapper.ScrapTNA(startDate2, finishDate2)
		if err != nil {
			fmt.Println(err)
		}
		w.db.Create(tnas[0])

		fmt.Println("holaaa")
		<-ticker.C
	}
}

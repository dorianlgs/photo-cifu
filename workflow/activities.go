package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/cschleiden/go-workflows/activity"
	"github.com/pocketbase/pocketbase"
)

type activities struct {
	pb *pocketbase.PocketBase
}

func (act *activities) Activity1(ctx context.Context, x, y int) (int, error) {
	logger := activity.Logger(ctx)
	logger.Info("Entering Activity1")

	recordId := "70915774024uqu1"

	record, err := act.pb.FindRecordById("galleries", recordId)

	if err != nil {
		logger.Error("err", "detail", err.Error())
		return 0, fmt.Errorf("record not found: %s", recordId)
	}

	logger.Info("OK!", "record", record.Get("name"))

	time.Sleep(10 * time.Second)

	return x + y, nil
}

func (act *activities) Activity2(ctx context.Context) (int, error) {

	logger := activity.Logger(ctx)
	logger.Info("Entering Activity2")

	time.Sleep(3 * time.Second)

	return 12, nil
}

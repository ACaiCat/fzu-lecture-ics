package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/west2-online/jwch"
)

type LectureRequest struct {
	UID      string `query:"uid,required"`
	Password string `query:"password,required"`
}

//goland:noinspection ALL
func GetLectureIcs(c context.Context, ctx *app.RequestContext) {
	var req LectureRequest

	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, utils.H{
			"msg": "invalid parameters: " + err.Error(),
		})
		return
	}

	stu, err := Login(req.UID, req.Password)
	if err != nil {
		ctx.JSON(consts.StatusForbidden, utils.H{
			"msg": err.Error(),
		})
		return
	}

	calendar, err := getCalendar(stu)

	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, utils.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.Data(consts.StatusOK, "text/calendar", calendar)
}

// 代码部分来自 https://github.com/renbaoshuo/fzu-ics
// 根据 GPL-3.0 License 的要求进行使用和分发
func getCalendar(stu *jwch.Student) ([]byte, error) {
	// 初始化
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = cstSh

	// 转换为 ics 格式
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetXWRCalName(fmt.Sprintf("福州大学讲座 [%s]", stu.ID))
	cal.SetTimezoneId("Asia/Shanghai")
	cal.SetXWRTimezone("Asia/Shanghai")

	lectures, err := stu.GetLectures()

	if err != nil {
		return nil, err
	}

	for _, lecture := range lectures {

		eventIdBase := fmt.Sprintf("%s__%d_%s_%s",
			lecture.Category, lecture.IssueNumber, lecture.Title, lecture.Speaker)
		description := fmt.Sprintf(
			"主讲人：%s\n"+
				"类别：%s\n"+
				"听取情况：%s\n",
			lecture.Speaker, lecture.Category, lecture.AttendanceStatus)

		event := cal.AddEvent(md5Str(eventIdBase))
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetSummary(lecture.Title)
		event.SetDescription(description)

		// 位置信息
		lat, lon := findGeoLocation(lecture.Location)
		if lat != 0 && lon != 0 {
			event.SetGeo(lat, lon)
		}
		event.SetLocation(lecture.Location)

		// 开始时间
		startTime := time.UnixMilli(lecture.Timestamp)
		event.SetStartAt(startTime)
		event.SetEndAt(startTime.Add(1 * time.Hour))

		// 提醒
		alarmDescription := "地点: " + lecture.Location + "\n"
		alarm := event.AddAlarm()
		alarm.SetAction(ics.ActionDisplay)
		alarm.SetSummary(lecture.Title)
		alarm.SetTrigger("-PT15M")
		alarm.SetDescription(alarmDescription)
	}

	calendarContent := cal.Serialize()

	return []byte(calendarContent), nil
}

var GEO = map[string][2]float64{
	"晋江校区A":   {118.585637992599, 24.557758583454426},
	"晋江校区B":   {118.58492906232163, 24.557395418615222},
	"大梦书屋":    {119.19733999999994, 26.05926300000002},
	"嘉锡楼":     {119.19447500000001, 26.059567},
	"阳光科技楼":   {119.20233499999995, 26.052188999999988},
	"阳光楼":     {119.20233499999995, 26.052188999999988},
	"福州大学图书馆": {119.19767194009012, 26.05896002955218},
	"晋江楼":     {119.20112561376197, 26.061205809110632},
	"铜盘科报厅A":  {26.103566850977305, 119.26282667315306},
	"铜盘科报厅B":  {26.103443904554403, 119.26236726136176},
}

func findGeoLocation(location string) (float64, float64) {
	for key, value := range GEO {
		if strings.Contains(location, key) {
			return value[0], value[1]
		}
	}
	return 0, 0
}

func md5Str(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	fullHash := hex.EncodeToString(hasher.Sum(nil)) // 32-bit (full) hash

	return fullHash
}

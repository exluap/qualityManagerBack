/**
 * Project IntelliJ IDEA created by exluap
 * Date: 08.10.2019 18:19
 * twitter: https://twitter.com/exluap
 * keybase: https://keybase.io/exluap
 * website: https://exluap.com
 */
package structs

import (
	"github.com/jinzhu/gorm"
	"time"
)

type SRSCommunication struct {
	gorm.Model
	ID          int
	Number      int
	CommCall    int
	CommMail    int
	CommMeet    int
	CommNothing int
	CommChat    int
} //Тип коммуникации

type SRS struct {
	gorm.Model
	ID       int
	Number   string
	Type     int
	Result   int  // Дата и время создания записи
	Waited   bool // Флаг ожидания (отложен)
	Overtime bool // Флаг овертайма
	Owner    int  // Создатель запроса
} //Информация по SRам

type SRSTypeDictionary struct {
	gorm.Model
	ID   int
	Name string
} //Справочник по типу КО

type SRSResultDictionary struct {
	gorm.Model
	ID   int
	Name string
} //Справочник по результатам КО

type SRSSettings struct {
	gorm.Model
	ID            int
	Number        int
	NoRecords     int
	NoRecordsOnly int
	Expenditure   int
	MoreThing     int
	ExpClaim      int
} // Параметры КО

type SRSAdditionalActions struct {
	gorm.Model
	ID                 int
	Number             int
	FinCorr            int
	CloseAccount       int
	NeedUnblock        int
	NeedCorrectLoyatly int
	DeniedPhone        int
	NeedOther          int
	Information        int
	DueDateAction      int
} // Дополнительные действия у КО

type SRSDueDateAction struct {
	gorm.Model
	ID   int
	Name string
} // Справочник действий с минимальным платежом

type Users struct {
	gorm.Model
	ID       int
	Login    string
	Password string
	UserId   int
	Token    string
} // Список пользователей

type UsersInfo struct {
	gorm.Model
	ID         int
	FirstName  string
	LastName   string
	MiddleName string
	Overtime   bool
	Email      string
} // Информация о пользователе

type UsersGroups struct {
	gorm.Model
	ID      int
	UserId  int
	GroupId int
} // Группы пользователя

type Groups struct {
	gorm.Model
	ID   int
	Name string // Наименование группы
} // Группы

type Blog struct {
	gorm.Model
	ID      int
	Owner   int
	Post    string
	Preview string
	Title   string
} // Блог

type Images struct {
	gorm.Model
	ID   int
	Path string
	Post int
}

type Likes struct {
	gorm.Model
	ID    int
	Post  int
	Owner int
} // Лайки

type Notifications struct {
	gorm.Model
	ID                  int
	NotificationText    string
	PrimaryNotification bool
	Owner               string
} // Уведомления

type NotificationsRead struct {
	gorm.Model
	ID             int
	UserId         int
	NotificationID int
} //Прочитанные Уведомления

type Tasks struct {
	gorm.Model
	ID     int
	Closed time.Time
	Author int
	Type   int
	Status int
	Owner  int
}

type TasksTypeDictionary struct {
	gorm.Model
	ID   int
	Name string
}

type TasksStatusDictionary struct {
	gorm.Model
	ID   int
	Name string
}

type TasksInfo struct {
	gorm.Model
	ID             int
	SRS            string
	TaskID         int
	LinkComm       string
	LinkEvaluation string
	Rated          string
	JiraLink       string
	BadQuestion    bool
	TaskInfo       string
	AccountId      string
	ContactId      string
}

type Comments struct {
	gorm.Model
	ID     int
	Owner  int
	TaskID int
	Hidden bool
}

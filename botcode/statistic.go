package botcode

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserStat struct {
	UserName     string `json:"user_name"`
	Id           int64
	MessageCount int       `json:"message_count"`
	LastTime     time.Time `json:"last_time"`
}

const fileName = "users.json"

func SaveStatistic(user *tgbotapi.User) {
	users := make(map[int64]UserStat)

	var statFile, err = os.OpenFile(fileName, os.O_RDWR, 0666)
	if err == nil {
		result, _ := ioutil.ReadFile(statFile.Name())
		err = json.Unmarshal(result, &users)
		if err != nil {
			log.Println(err)
		}
	} else {
		statFile, _ = os.Create(fileName)
	}

	u, is := users[user.ID]
	if !is {
		u = UserStat{
			UserName: user.UserName,
			Id:       user.ID,
		}
	}
	u.MessageCount++
	u.LastTime = time.Now()

	users[user.ID] = u
	saveUsersStatToFile(users, statFile)
	_ = statFile.Close()
}

func saveUsersStatToFile(usersMap map[int64]UserStat, file *os.File) {
	jsn, _ := json.MarshalIndent(usersMap, " ", "")
	er := ioutil.WriteFile(file.Name(), jsn, os.ModeExclusive)
	if er != nil {
		log.Println(er)
	}
}

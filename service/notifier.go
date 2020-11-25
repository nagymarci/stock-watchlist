package service

//go:generate $GOPATH/bin/mockgen -source=notifier.go -destination=mocks/mock_notifier-deps.go -package=mocks
import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	userprofileModel "github.com/nagymarci/stock-user-profile/model"
	"github.com/nagymarci/stock-watchlist/model"
)

type Notifier struct {
	recommendations   recommendationProvider
	watchlists        watchlistList
	stockClient       stockGetter
	stockService      stockRecommendator
	userprofileClient userprofileGetter
	emailClient       emailSender
}

type watchlistList interface {
	List() ([]model.Watchlist, error)
}

type recommendationProvider interface {
	Get(id primitive.ObjectID) ([]string, error)
	Update(log *logrus.Entry, id primitive.ObjectID, stocks []string) error
}

type emailSender interface {
	SendNotification(profileName string, removed, added, currentStocks []string, email string) error
}

type stockGetter interface {
	Get(symbol string) (model.StockData, error)
}

type stockRecommendator interface {
	GetAllRecommendedStock(stocks []model.StockData, numReqs int, userprofile *userprofileModel.Userprofile) []model.CalculatedStockInfo
}

type userprofileGetter interface {
	GetUserprofile(userId string) (userprofileModel.Userprofile, error)
}

func NewNotifier(r recommendationProvider, w watchlistList, sc stockGetter, ss stockRecommendator, uc userprofileGetter, ec emailSender) *Notifier {
	return &Notifier{
		recommendations:   r,
		watchlists:        w,
		stockClient:       sc,
		stockService:      ss,
		userprofileClient: uc,
		emailClient:       ec,
	}
}

func (n *Notifier) NotifyChanges() {
	watchlists, err := n.watchlists.List()

	if err != nil {
		logrus.Errorf("Failed to get watchlists [%v]", err)
		return
	}

	for _, watchlist := range watchlists {
		log := logrus.WithField("watchlistId", watchlist.ID)
		previouStocks, _ := n.recommendations.Get(watchlist.ID)

		var stockInfos []model.StockData

		for _, symbol := range watchlist.Stocks {
			result, err := n.stockClient.Get(symbol)

			if err != nil {
				log.Warnf("Failed to get stock [%s]: [%v]\n", symbol, err)
				continue
			}

			stockInfos = append(stockInfos, result)
		}

		userprofile, err := n.userprofileClient.GetUserprofile(watchlist.UserID)

		if err != nil {
			log.Errorln("Failed to get userprofile to notification ", err)
			continue
		}

		calculatedStockData := n.stockService.GetAllRecommendedStock(stockInfos, 2, &userprofile)

		currentStocks := filterGreenPrices(calculatedStockData)

		removed, added := getChanges(previouStocks, currentStocks)

		if len(removed) == 0 && len(added) == 0 {
			continue
		}

		err = n.emailClient.SendNotification(watchlist.Name, removed, added, currentStocks, userprofile.Email)

		if err != nil {
			log.Errorln("Failed to send notification ", err)
			continue
		}

		n.recommendations.Update(log, watchlist.ID, currentStocks)
	}
}

func filterGreenPrices(stockInfos []model.CalculatedStockInfo) []string {
	var result []string

	for _, calc := range stockInfos {
		if calc.PriceColor != "green" {
			continue
		}

		result = append(result, calc.Ticker)
	}
	return result
}

func getChanges(old, new []string) ([]string, []string) {
	if len(old) == 0 || len(new) == 0 {
		return old, new
	}

	var removed []string

	for _, symbol := range old {
		if !contains(new, symbol) {
			removed = append(removed, symbol)
		}
	}

	var added []string

	for _, symbol := range new {
		if !contains(old, symbol) {
			added = append(added, symbol)
		}
	}

	return removed, added
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

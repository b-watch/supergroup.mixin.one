package models

import (
	"container/heap"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/jinzhu/now"
	"github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
)

var rankHeight int = 100

var rankRanges = map[string]func(time.Time) TimeRange{
	"week": func(t time.Time) TimeRange {
		t = t.UTC()
		weekStart := now.With(t).Monday().Add(time.Minute)
		return TimeRange{weekStart, weekStart.Add(7 * 24 * time.Hour).Add(-time.Nanosecond)}
	},
	"month": func(t time.Time) TimeRange {
		t = t.UTC()
		monthStart := now.With(t).BeginningOfMonth().Add(time.Minute)
		monthEnd := now.With(t).EndOfMonth().Add(time.Minute)
		return TimeRange{monthStart, monthEnd}
	},
	"all": func(t time.Time) TimeRange {
		return TimeRange{}
	},
}

var RankManager = &rankManager{}

var rankPGChannelName = "TipSender"

var rankPriceCache = cache.New(5*time.Minute, 10*time.Minute)

var dumbTipSum = TipSum{}

type TipSum struct {
	UserID    string          `json:"user_id"`
	AvatarURL string          `json:"avatar_url"`
	FullName  string          `json:"full_name"`
	TipDetail TipDetail       `json:"tip_detail"`
	TipUSD    decimal.Decimal `json:"tip_usd"`
	TipCount  int             `json:"tip_count"`
	Index     int             `json:"-"`
}

type TipDetail map[string]decimal.Decimal

type Rank []*TipSum

type rankManager struct {
	Ranks map[string]*RankContainer
}

type RankContainer struct {
	RankRange TimeRange
	TopUsers  Rank
	Height    int
	rwMutex   sync.RWMutex
}

type TimeRange struct {
	TimeStart time.Time
	TimeEnd   time.Time
}

type Tip struct {
	SenderID    string
	RecipientID string
	Detail      TipDetail
	TraceID     string
	Time        time.Time
}

type RankResult struct {
	CurrentRank TipSum          `json:"current_rank"`
	Ranks       map[string]Rank `json:"ranks"`
}

func (tr TimeRange) IsInfinite() bool {
	if tr.TimeStart.IsZero() && tr.TimeEnd.IsZero() {
		return true
	}
	return false
}

func (current *User) ShowTiprank(ctx context.Context) (*RankResult, error) {
	ranks, err := listTipRanks(ctx)
	if err != nil {
		return nil, err
	}
	var tipSum TipSum
	tipSum, err = RankManager.pullUser(ctx, current.UserId, TimeRange{})
	tipSum.AvatarURL = current.AvatarURL
	tipSum.FullName = current.FullName
	tipSum.TipUSD = tipSum.TotalUSD()
	rankRes := RankResult{Ranks: ranks, CurrentRank: tipSum}
	return &rankRes, nil
}

func CreateTip(ctx context.Context, senderID, recipientID, assetID, amount, traceID string, tipTime time.Time) error {
	tipAmount, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}
	tipDetail := make(TipDetail)
	tipDetail[assetID] = tipAmount
	var tipDetailJson []byte
	tipDetailJson, err = json.Marshal(tipDetail)
	if err != nil {
		return err
	}
	err = session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "INSERT INTO tips (sender_id, recipient_id, detail, trace_id, time) VALUES ($1, $2, $3, $4, $5)", senderID, recipientID, string(tipDetailJson), traceID, string(pq.FormatTimestamp(tipTime)))
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "SELECT pg_notify($1, $2)", rankPGChannelName, senderID)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (rm *rankManager) Init(ctx context.Context, dsn string) error {
	rm.Ranks = map[string]*RankContainer{}
	for rankName, timeFn := range rankRanges {
		timeRange := timeFn(time.Now())
		if err := rm.newRank(ctx, rankName, timeRange); err != nil {
			return err
		}
	}
	rm.tick(ctx, dsn)
	return nil
}

func (rm *rankManager) Push(ctx context.Context, senderID string) {
	for _, rc := range rm.Ranks {
		tipSum, err := rm.pullUser(ctx, senderID, rc.RankRange)
		if err != nil {
			fmt.Println(err)
			return
		}
		rc.push(&tipSum)
	}
}

func (rm *rankManager) List() (ranks map[string]Rank) {
	ranks = make(map[string]Rank)
	for rankName, rc := range rm.Ranks {
		var tips Rank
		for _, u := range rc.list() {
			if u != &dumbTipSum {
				tips = append(tips, u)
			}
		}
		ranks[rankName] = tips
	}
	return
}

func (rm *rankManager) newRank(ctx context.Context, rankName string, timeRange TimeRange) error {
	rc := RankContainer{RankRange: timeRange, Height: rankHeight}
	tipSums, err := rm.pullUsers(ctx, timeRange)
	if err != nil {
		return err
	}
	rc.init(tipSums)
	rm.Ranks[rankName] = &rc
	return nil
}

func (rm *rankManager) tick(ctx context.Context, dsn string) {
	ticker := time.NewTicker(time.Minute)
	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	if err := listener.Listen(rankPGChannelName); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case t := <-ticker.C:
				for rankName, timeFn := range rankRanges {
					rc := rm.Ranks[rankName]
					timeRange := timeFn(t)
					if rc.RankRange != timeRange {
						tipSums, _ := rm.pullUsers(ctx, timeRange)
						rc.reset(timeRange, tipSums)
					}
				}
			case n := <-listener.Notify:
				if n != nil {
					rm.Push(ctx, n.Extra)
				}
			}
		}
	}()
}

func (rm *rankManager) pullUsers(ctx context.Context, timeRange TimeRange) (tips []TipSum, err error) {
	err = session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var rows *sql.Rows
		var err error
		if timeRange.IsInfinite() {
			rows, err = tx.QueryContext(
				ctx,
				`select sender_id as user_id, jsonb_object_agg(key, val) as tip_detail, sum(tip_count)
				from (
					select sender_id, key, sum(value::numeric) val, count(sender_id) tip_count
					from tips t, jsonb_each_text(detail)
					group by t.sender_id, key
				) s
				group by user_id;`,
			)
		} else {
			rows, err = tx.QueryContext(
				ctx,
				`select sender_id as user_id, jsonb_object_agg(key, val) as tip_detail, sum(tip_count)
				from (
					select sender_id, key, sum(value::numeric) val, count(sender_id) tip_count
					from tips t, jsonb_each_text(detail)
					where t.time between $1 and $2
					group by t.sender_id, key
				) s
				group by user_id;`,
				timeRange.TimeStart,
				timeRange.TimeEnd,
			)
		}
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var tipSum TipSum
			var tipDetailJson []byte
			if err := rows.Scan(&tipSum.UserID, &tipDetailJson, &tipSum.TipCount); err != nil {
				return err
			}
			err := json.Unmarshal(tipDetailJson, &tipSum.TipDetail)
			if err != nil {
				return err
			}
			tips = append(tips, tipSum)
		}
		return nil
	})
	return
}

func (rm *rankManager) pullUser(ctx context.Context, userID string, timeRange TimeRange) (tipSum TipSum, err error) {
	err = session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var row *sql.Row
		if timeRange.IsInfinite() {
			row = tx.QueryRowContext(
				ctx,
				`select sender_id as user_id, jsonb_object_agg(key, val) as tip_detail, sum(tip_count)
				from (
					select sender_id, key, sum(value::numeric) val, count(sender_id) tip_count
					from tips t, jsonb_each_text(detail)
					where t.sender_id = $1
					group by sender_id, key
				) s
				group by user_id;`,
				userID,
			)
		} else {
			row = tx.QueryRowContext(
				ctx,
				`select sender_id as user_id, jsonb_object_agg(key, val) as tip_detail, sum(tip_count)
				from (
					select sender_id, key, sum(value::numeric) val, count(sender_id) tip_count
					from tips t, jsonb_each_text(detail)
					where t.sender_id = $1 and t.time between $2 and $3
					group by sender_id, key
				) s
				group by user_id;`,
				userID,
				timeRange.TimeStart,
				timeRange.TimeEnd,
			)
		}
		var tipDetailJson []byte
		err := row.Scan(&tipSum.UserID, &tipDetailJson, &tipSum.TipCount)
		if err != nil {
			return err
		}
		err = json.Unmarshal(tipDetailJson, &tipSum.TipDetail)
		return err
	})
	return
}

func listTipRanks(ctx context.Context) (ranks map[string]Rank, err error) {
	ranks = RankManager.List()
	userFilter := make(map[string]bool)
	for _, users := range ranks {
		for _, u := range users {
			if _, ok := userFilter[u.UserID]; !ok {
				userFilter[u.UserID] = true
			}
		}
	}

	userIDs := make([]string, 0, len(userFilter))
	for k := range userFilter {
		userIDs = append(userIDs, k)
	}

	var groupUsers []*User
	if len(userIDs) == 0 {
		return
	}
	groupUsers, err = FindUsers(ctx, userIDs)
	if err != nil {
		return
	}

	groupUsersMap := make(map[string]*User)
	for _, u := range groupUsers {
		groupUsersMap[u.UserId] = u
	}

	for i := range ranks {
		for j, u := range ranks[i] {
			gu := groupUsersMap[u.UserID]
			ranks[i][j].AvatarURL = gu.AvatarURL
			ranks[i][j].FullName = gu.FullName
			ranks[i][j].TipUSD = ranks[i][j].TotalUSD()
		}
	}

	return
}

func (tipSum TipSum) TotalUSD() (totalUSD decimal.Decimal) {
	for assetID, amount := range tipSum.TipDetail {
		usdPrice := assetPriceUSD(assetID)
		totalUSD = totalUSD.Add(amount.Mul(usdPrice))
	}
	totalUSD = totalUSD.Round(3)
	return
}

func (rc *RankContainer) init(tips []TipSum) {
	rc.rwMutex.Lock()
	defer rc.rwMutex.Unlock()
	rc.TopUsers = make(Rank, rc.Height)

	userCount := len(tips)
	for i := range rc.TopUsers {
		if i > userCount-1 {
			rc.TopUsers[i] = &dumbTipSum
		} else {
			var t TipSum
			t, tips = tips[0], tips[1:]
			rc.TopUsers[i] = &t
		}
	}

	heap.Init(&rc.TopUsers)

	for i := range tips {
		topLast := rc.TopUsers.peak()
		if tips[i].TotalUSD().GreaterThan(topLast.TotalUSD()) {
			rc.TopUsers[0] = &tips[i]
			heap.Fix(&rc.TopUsers, 0)
		}
	}

	fmt.Printf("rank container: %s~%s \n", rc.RankRange.TimeStart, rc.RankRange.TimeEnd)
}

func (rc *RankContainer) push(t *TipSum) {
	rc.rwMutex.Lock()
	defer rc.rwMutex.Unlock()
	for i, tip := range rc.TopUsers {
		// while user is already in top
		if t.UserID == tip.UserID {
			rc.TopUsers[i] = t
			heap.Fix(&rc.TopUsers, i)
			return
		}
	}

	// while user not inside top
	topLast := rc.TopUsers.peak()
	if t.TotalUSD().GreaterThan(topLast.TotalUSD()) {
		rc.TopUsers[0] = t
		heap.Fix(&rc.TopUsers, 0)
	}
}

func (rc *RankContainer) reset(timeRange TimeRange, tips []TipSum) {
	rc.rwMutex.Lock()
	rc.RankRange = timeRange
	rc.rwMutex.Unlock()

	rc.init(tips)
}

func (rc *RankContainer) list() (temp Rank) {
	rc.rwMutex.Lock()
	temp = append(temp, rc.TopUsers...)
	rc.rwMutex.Unlock()

	sort.Sort(sort.Reverse(Rank(temp)))
	return temp
}

func (r Rank) Len() int { return len(r) }

func (r Rank) Less(i, j int) bool {
	return r[i].TotalUSD().LessThan(r[j].TotalUSD())
}

func (r Rank) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
	r[i].Index = i
	r[j].Index = j
}

func (r *Rank) Push(x interface{}) {
	n := len(*r)
	t := x.(*TipSum)
	t.Index = n
	*r = append(*r, t)
}

func (r Rank) peak() *TipSum {
	return r[0]
}

func (r *Rank) Pop() interface{} {
	old := *r
	n := len(old)
	u := old[n-1]
	old[n-1] = nil // avoid memory leak
	u.Index = -1   // for safety
	*r = old[0 : n-1]
	return u
}

func assetPriceUSD(assetID string) (usdPrice decimal.Decimal) {
	if p, found := rankPriceCache.Get(assetID); found {
		usdPrice = p.(decimal.Decimal)
	} else {
		assetUSDPrice, err := botReadAssetUSDPrice(assetID)
		if err != nil {
			fmt.Printf("currency price not found %s\n", assetID)
		}

		usdPrice, err = decimal.NewFromString(assetUSDPrice)
		if err != nil {
			fmt.Printf("price converion failed %s\n", assetID)
		}
		rankPriceCache.Set(assetID, usdPrice, cache.DefaultExpiration)
	}
	return
}

func botReadAssetUSDPrice(assetID string) (string, error) {
	accessToken, err := botAuthenticationToken("GET", "/assets/"+assetID, "")
	if err != nil {
		return "", err
	}

	asset, err := bot.AssetShow(context.TODO(), assetID, accessToken)
	if err != nil {
		return "", err
	}
	return asset.PriceUSD, nil
}

func botAuthenticationToken(httpMethod string, uri string, body string) (token string, err error) {
	token, err = bot.SignAuthenticationToken(config.AppConfig.Mixin.ClientId, config.AppConfig.Mixin.SessionId, config.AppConfig.Mixin.SessionKey, httpMethod, uri, body)
	return
}

package models

import (
	"context"
	"database/sql"
	"math"
	"testing"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"
	number "github.com/MixinNetwork/go-number"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/stretchr/testify/assert"
)

func TestPacketCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	user, err := createUser(ctx, "accessToken", bot.UuidNewV4().String(), "1000", "name", "http://localhost")
	assert.Nil(err)
	assert.NotNil(user)
	err = user.Subscribe(ctx)
	assert.Nil(err)
	sum, err := user.Prepare(ctx)
	assert.Nil(err)
	assert.Equal(int64(1), sum)

	li, err := createUser(ctx, "accessToken", bot.UuidNewV4().String(), "1001", "Li", "http://localhost")
	assert.Nil(err)
	assert.NotNil(li)
	err = li.Subscribe(ctx)
	assert.Nil(err)
	sum, err = user.Prepare(ctx)
	assert.Nil(err)
	assert.Equal(int64(2), sum)

	asset := &Asset{
		AssetId:  bot.UuidNewV4().String(),
		Symbol:   "XIN",
		Name:     "Mixin",
		IconURL:  "http://mixin.one",
		PriceBTC: "0",
		PriceUSD: "0",
		Balance:  "100",
	}
	err = upsertAssets(ctx, []*Asset{asset})
	assert.Nil(err)
	packet, err := li.createPacket(ctx, asset, number.FromString("1"), 2, "Hello Packet")
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(PacketStateInitial, packet.State)
	packet, err = PayPacket(ctx, packet.PacketId, asset.AssetId, "1")
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(PacketStatePaid, packet.State)
	packet, err = ShowPacket(ctx, packet.PacketId)
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal("1", packet.Amount)
	assert.Equal(int64(2), packet.TotalCount)
	packet, err = li.ClaimPacket(ctx, packet.PacketId)
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(PacketStatePaid, packet.State)
	packet, err = RefundPacket(ctx, packet.PacketId)
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(PacketStatePaid, packet.State)
	packet, err = user.ClaimPacket(ctx, packet.PacketId)
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(PacketStatePaid, packet.State)
	packet, err = ShowPacket(ctx, packet.PacketId)
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(int64(0), packet.RemainingCount)
	assert.Equal("0", packet.RemainingAmount)
	assert.Len(packet.Participants, 2)
	packet, err = li.createPacket(ctx, asset, number.FromString("1"), 2, "Hello Packet")
	assert.Nil(err)
	assert.NotNil(packet)
	packet, err = PayPacket(ctx, packet.PacketId, asset.AssetId, "1")
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(PacketStatePaid, packet.State)
	_, err = session.Database(ctx).ExecContext(ctx, "UPDATE packets SET created_at=$1 WHERE packet_id=$2", time.Now().Add(-25*time.Hour), packet.PacketId)
	assert.Nil(err)
	packet, err = ShowPacket(ctx, packet.PacketId)
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(PacketStateExpired, packet.State)
	ids, err := ListExpiredPackets(ctx, 100)
	assert.Nil(err)
	assert.Len(ids, 1)
	packet, err = RefundPacket(ctx, packet.PacketId)
	assert.Nil(err)
	assert.NotNil(packet)
	assert.Equal(PacketStateRefunded, packet.State)
	packet, err = testReadPacketWithRelation(ctx, packet.PacketId)
	assert.Nil(err)
	assert.NotNil(packet)
	packet, err = testReadPacketWithRelation(ctx, bot.UuidNewV4().String())
	assert.Nil(err)
	assert.Nil(packet)
}

func testReadPacketWithRelation(ctx context.Context, packetId string) (*Packet, error) {
	var packet *Packet
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		packet, err = readPacketWithAssetAndUser(ctx, tx, packetId)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return packet, err
}

func TestPacketPreAllocation(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	totalCount := int64(100)
	amount := number.FromString("10")
	allocation, err := packetPreAllocate(totalCount, amount)
	assert.Nil(err)
	assert.Len(allocation, int(totalCount))
	t.Log(allocation)

	var dist []float64
	var allocationSum number.Decimal
	for _, allocateAmount := range allocation {
		a := number.FromString(allocateAmount)
		allocationSum = allocationSum.Add(a)
		dist = append(dist, a.Float64())
	}
	assert.True(amount.Cmp(allocationSum) >= 0)

	testNormalDistribution(t, dist)
}

func testNormalDistribution(t *testing.T, sample []float64) {
	var μ, σ float64
	var sum float64
	for _, v := range sample {
		sum += v
	}
	μ = float64(sum) / float64(len(sample))

	var variance float64
	for _, v := range sample {
		variance += math.Pow((v - μ), 2)
	}
	σ = math.Sqrt(variance / float64(len(sample)))

	t.Logf("μ: %f", μ)
	t.Logf("σ: %f", σ)
}

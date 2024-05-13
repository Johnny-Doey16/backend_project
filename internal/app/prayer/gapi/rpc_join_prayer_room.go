package gapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/steve-mir/diivix_backend/internal/app/authentication/token"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/pb"
	"github.com/steve-mir/diivix_backend/internal/app/prayer/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	OneCoinDuration = 30 // minutes
	TwoCoinDuration = 45 // minutes
)

func calculateDuration(coins int64) (time.Duration, error) {
	if coins < 1 {
		return time.Microsecond, errors.New("sorry you do not have enough coin")
	}
	switch coins {
	case 1:
		return time.Duration(OneCoinDuration) * time.Minute, nil
	case 2:
		return time.Duration(TwoCoinDuration) * time.Minute, nil
	default:
		return time.Duration(TwoCoinDuration*(coins-1)) * time.Minute, nil
	}
}

// TODO: Implement - act as tunnel to signalling server
func (s *PrayerServer) JoinPrayerRoom(ctx context.Context, req *pb.PrayerRoomId) (*pb.PrayerRoomResponse, error) {
	claims, ok := ctx.Value("payloadKey").(*token.Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unable to retrieve UID from context")
	}
	// uid, err := post_service.StrToUUID("99dc65d6-6bec-4f72-8df0-18949eea1ff8")
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "parsing uid: %s", err)
	// }

	// TODO: If user joins before the time show before time error

	//! Get prayer duration with prayer id
	prayerRoom, err := s.store.GetPrayerRoomById(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting prayer room by id %s", err)
	}

	//! Check if the prayer is still active. The end time has been elapsed
	if time.Now().After(prayerRoom.EndTime.Time) {
		return nil, status.Error(codes.Internal, "prayer meeting expired")
	}

	//! Calculate the duration of the prayer room.
	prayerRoomDuration := prayerRoom.EndTime.Time.Sub(prayerRoom.StartTime.Time)
	coinCost := int(prayerRoomDuration.Minutes() / 30)

	if int(prayerRoomDuration.Minutes())%30 != 0 {
		coinCost++
	}
	log.Println("Duration", prayerRoomDuration, "Costs", coinCost)

	//! subtract the coinCost from the total_coin
	err = subtractCoin(ctx, s.db, claims.UserId, float64(coinCost))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error subtracting %v", err)
	}

	// ! get room details
	url, err := services.JoinRoomInSignallingServer(prayerRoom.RoomID, prayerRoom.Name.String)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error joining %v", err)
	}

	// return room_id and webRTC websocket url. In signalling server while joining check if the room_id exist, if it doesn't create a room with that id and join.
	return &pb.PrayerRoomResponse{
		Msg:       "successfully joined prayer room",
		RoomId:    prayerRoom.RoomID,
		RoomName:  prayerRoom.Name.String,
		WebRtcUrl: url,
	}, nil
}

func subtractCoin(ctx context.Context, db *sql.DB, uid uuid.UUID, amount float64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "transaction error: %v", err.Error())
	}
	defer tx.Rollback()

	// Executing the SQL code within the transaction
	query := fmt.Sprintf(`
        DO $$
        BEGIN
            UPDATE accounts
            SET total_coin = CASE 
                                WHEN total_coin >= %f THEN total_coin - %f
                                ELSE total_coin 
                            END
            WHERE user_id = '%s';

            IF NOT EXISTS (SELECT 1 FROM accounts WHERE user_id = '%s' AND total_coin >= %f) THEN
                RAISE EXCEPTION 'Insufficient funds';
            END IF;
        END $$;
    `, amount, amount, uid, uid, amount)
	_, err = tx.Exec(query)

	if err != nil {
		// If there's an error, rollback the transaction
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

package services

import (
	"errors"
	"scoring_app/models"
	"scoring_app/repositories"

	"fmt"
	"strconv"
)

type MatchRoundService interface {
	Create(round *models.MatchRound) error
	CreatePairing(matchID string) ([]models.MatchRound, error)
	UpdateScores(roundID uint, court string, updates map[string]interface{}) error
	UpdateScoresBulk(roundID uint, items []models.UpdateScoreItem) ([]models.UpdateScoreResponse, error)
	Update(round *models.MatchRound) error
	GetAll(page int, pageSize int) ([]models.MatchRound, int64, error)
	GetByID(id uint) (*models.MatchRound, error)
	GetByMatchID(matchID string) ([]models.MatchRound, error)
	Delete(roundID uint) error
	DeleteByMatchID(matchID string) error
}

type matchRoundService struct {
	repo            repositories.MatchRoundRepository
	matchEventRepo  repositories.MatchEventRepository
	matchPlayerRepo repositories.MatchPlayerRepository
}

func NewMatchRoundService(
	repo repositories.MatchRoundRepository,
	matchEventRepo repositories.MatchEventRepository,
	matchPlayerRepo repositories.MatchPlayerRepository,
) MatchRoundService {
	return &matchRoundService{
		repo:            repo,
		matchEventRepo:  matchEventRepo,
		matchPlayerRepo: matchPlayerRepo,
	}
}

// CREATE
func (s *matchRoundService) Create(round *models.MatchRound) error {

	_, err := s.matchEventRepo.FindByID(round.MatchID)
	if err != nil {
		return errors.New("match event not found")
	}

	playerIDs := []uint{
		round.TeamAPlayer1ID,
		round.TeamAPlayer2ID,
		round.TeamBPlayer1ID,
		round.TeamBPlayer2ID,
	}

	for _, playerID := range playerIDs {
		_, err := s.matchPlayerRepo.FindByID(playerID)
		if err != nil {
			return errors.New("one or more players not found")
		}
	}

	existingRounds, _ := s.repo.FindByRoundNumber(round.MatchID, round.RoundNumber)
	for _, r := range existingRounds {
		if r.Court == round.Court {
			return errors.New("round and court already exists")
		}
	}

	return s.repo.Create(round)
}

// UPDATE SCORES
func (s *matchRoundService) UpdateScores(roundID uint, court string, updates map[string]interface{}) error {

	courtInt, err := strconv.Atoi(court)
	if err != nil {
		return errors.New("invalid court number")
	}

	_, err = s.repo.FindByIDAndCourt(roundID, courtInt)
	if err != nil {
		return errors.New("round not found on this court")
	}

	scoreA := 0
	scoreB := 0

	if v, ok := updates["score_a"]; ok {
		scoreA = v.(int) // 🔥 cukup int saja
	}

	if v, ok := updates["score_b"]; ok {
		scoreB = v.(int)
	}

	return s.repo.PatchScore(roundID, courtInt, scoreA, scoreB)
}

// UPDATE SCORES BULK
func (s *matchRoundService) UpdateScoresBulk(roundID uint, items []models.UpdateScoreItem) ([]models.UpdateScoreResponse, error) {
	var responses []models.UpdateScoreResponse

	for _, item := range items {
		// Validasi record ada
		existingRound, err := s.repo.FindByIDAndCourt(roundID, item.Court)
		if err != nil {
			return nil, fmt.Errorf("round not found on court %d", item.Court)
		}

		// Get current total score
		currentTotal := existingRound.ScoreA + existingRound.ScoreB

		scoreA := 0
		scoreB := 0

		if item.ScoreA != nil {
			scoreA = *item.ScoreA
		}

		if item.ScoreB != nil {
			scoreB = *item.ScoreB
		}

		// Skip kalau dua-duanya kosong
		if item.ScoreA == nil && item.ScoreB == nil {
			continue
		}

		// Check if total score already >= 21
		if currentTotal >= 21 {
			// Return current status without updating
			response := models.UpdateScoreResponse{
				Court:       item.Court,
				ScoreA:      existingRound.ScoreA,
				ScoreB:      existingRound.ScoreB,
				TotalScore:  currentTotal,
				MatchStatus: "session_end",
				CanUpdate:   false,
			}
			responses = append(responses, response)
			continue
		}

		// Validasi kalau total score akan melebihi 21, tidak akan terupdate dan langsung return
		if currentTotal+scoreA+scoreB > 21 {
			return nil, fmt.Errorf("score update on court %d would exceed total score of 21", item.Court)
		}

		// Update score
		if err := s.repo.PatchScore(roundID, item.Court, scoreA, scoreB); err != nil {
			return nil, err
		}

		// Get updated round
		updatedRound, err := s.repo.FindByIDAndCourt(roundID, item.Court)
		if err != nil {
			return nil, err
		}

		// Calculate new total and status
		newTotal := updatedRound.ScoreA + updatedRound.ScoreB
		matchStatus := s.calculateMatchStatus(newTotal)

		response := models.UpdateScoreResponse{
			Court:       item.Court,
			ScoreA:      updatedRound.ScoreA,
			ScoreB:      updatedRound.ScoreB,
			TotalScore:  newTotal,
			MatchStatus: matchStatus,
			CanUpdate:   newTotal < 21,
		}

		// If session ended, accumulate the final round score into each player's temp_score
		if newTotal >= 21 {
			// Team A players: add updatedRound.ScoreA
			if updatedRound.TeamAPlayer1ID != 0 {
				p, err := s.matchPlayerRepo.FindByID(updatedRound.TeamAPlayer1ID)
				if err != nil {
					return nil, err
				}
				p.TempScore += updatedRound.ScoreA
				if err := s.matchPlayerRepo.Update(p); err != nil {
					return nil, err
				}
			}

			if updatedRound.TeamAPlayer2ID != 0 {
				p, err := s.matchPlayerRepo.FindByID(updatedRound.TeamAPlayer2ID)
				if err != nil {
					return nil, err
				}
				p.TempScore += updatedRound.ScoreA
				if err := s.matchPlayerRepo.Update(p); err != nil {
					return nil, err
				}
			}

			// Team B players: add updatedRound.ScoreB
			if updatedRound.TeamBPlayer1ID != 0 {
				p, err := s.matchPlayerRepo.FindByID(updatedRound.TeamBPlayer1ID)
				if err != nil {
					return nil, err
				}
				p.TempScore += updatedRound.ScoreB
				if err := s.matchPlayerRepo.Update(p); err != nil {
					return nil, err
				}
			}

			if updatedRound.TeamBPlayer2ID != 0 {
				p, err := s.matchPlayerRepo.FindByID(updatedRound.TeamBPlayer2ID)
				if err != nil {
					return nil, err
				}
				p.TempScore += updatedRound.ScoreB
				if err := s.matchPlayerRepo.Update(p); err != nil {
					return nil, err
				}
			}
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// CALCULATE MATCH STATUS based on total score
func (s *matchRoundService) calculateMatchStatus(totalScore int) string {
	if totalScore >= 21 {
		return "session_end"
	} else if totalScore < 5 {
		return "session_1"
	} else if totalScore < 10 {
		return "session_2"
	} else if totalScore < 15 {
		return "session_3"
	} else {
		return "session_4"
	}
}

// UPDATE
func (s *matchRoundService) Update(round *models.MatchRound) error {

	existingRound, err := s.repo.FindByID(round.ID)
	if err != nil {
		return errors.New("round not found")
	}

	_, err = s.matchEventRepo.FindByID(round.MatchID)
	if err != nil {
		return errors.New("match event not found")
	}

	playerIDs := []uint{
		round.TeamAPlayer1ID,
		round.TeamAPlayer2ID,
		round.TeamBPlayer1ID,
		round.TeamBPlayer2ID,
	}

	for _, playerID := range playerIDs {
		_, err := s.matchPlayerRepo.FindByID(playerID)
		if err != nil {
			return errors.New("one or more players not found")
		}
	}

	existingRound.MatchID = round.MatchID
	existingRound.RoundNumber = round.RoundNumber
	existingRound.Court = round.Court
	existingRound.TeamAPlayer1ID = round.TeamAPlayer1ID
	existingRound.TeamAPlayer2ID = round.TeamAPlayer2ID
	existingRound.TeamBPlayer1ID = round.TeamBPlayer1ID
	existingRound.TeamBPlayer2ID = round.TeamBPlayer2ID
	existingRound.ScoreA = round.ScoreA
	existingRound.ScoreB = round.ScoreB

	return s.repo.Update(existingRound)
}

// GET ALL
func (s *matchRoundService) GetAll(page int, pageSize int) ([]models.MatchRound, int64, error) {
	return s.repo.GetAll(page, pageSize)
}

// GET BY ID
func (s *matchRoundService) GetByID(id uint) (*models.MatchRound, error) {
	return s.repo.FindByID(id)
}

// GET BY MATCH ID
func (s *matchRoundService) GetByMatchID(matchID string) ([]models.MatchRound, error) {
	return s.repo.FindByMatchIDWithPlayers(matchID)
}

// DELETE
func (s *matchRoundService) Delete(roundID uint) error {

	_, err := s.repo.FindByID(roundID)
	if err != nil {
		return errors.New("round not found")
	}

	return s.repo.Delete(roundID)
}

// DELETE BY MATCH ID
func (s *matchRoundService) DeleteByMatchID(matchID string) error {

	return s.repo.DeleteByMatchID(matchID)
}

// CREATE PAIRING
func (s *matchRoundService) CreatePairing(matchID string) ([]models.MatchRound, error) {
	// 1. Validasi bahwa match event ada
	matchEvent, err := s.matchEventRepo.FindByID(matchID)
	if err != nil {
		return nil, errors.New("match event not found")
	}

	// 2. Get all players untuk matchID
	players, err := s.matchPlayerRepo.FindByMatchID(matchID)
	if err != nil {
		return nil, err
	}

	// 3. Validasi jumlah player harus kelipatan 4
	totalPlayers := len(players)
	if totalPlayers == 0 || totalPlayers%4 != 0 {
		return nil, fmt.Errorf("invalid number of players: %d. Must be a multiple of 4 (4, 8, 12, etc.)", totalPlayers)
	}

	// 4. Get ranked players (sorted by temp_score DESC)
	rankedPlayers := s.getRankedPlayers(players)

	// 5. Get existing rounds untuk menentukan round number dan pattern
	existingRounds, _ := s.repo.FindByMatchID(matchID)

	// 6. Tentukan round number berikutnya
	nextRoundNumber := 1
	maxCourt := 0

	if len(existingRounds) > 0 {
		// Cari round number tertinggi
		for _, round := range existingRounds {
			if round.RoundNumber > nextRoundNumber {
				nextRoundNumber = round.RoundNumber
			}
			if round.Court > maxCourt {
				maxCourt = round.Court
			}
		}
		nextRoundNumber++ // increment untuk round baru
	}

	// 7. Hitung jumlah court yang dibutuhkan dari total_courts di match event
	totalCourts := matchEvent.TotalCourts
	if totalCourts < 1 {
		totalCourts = 1
	}

	// 8. Buat pairing berdasarkan kondisi
	var newRounds []models.MatchRound

	if len(existingRounds) == 0 {
		// Pattern 1: First round (no existing data) - vertical pattern
		newRounds = s.createPairingPattern1(rankedPlayers, matchID, nextRoundNumber, totalCourts)
	} else if maxCourt >= 1 && maxCourt <= 2 {
		// Pattern 2: Court 1-2 - alternating vertical pattern
		newRounds = s.createPairingPattern2(rankedPlayers, matchID, nextRoundNumber, totalCourts)
	} else {
		// Pattern 3: Court 3+ - horizontal pattern
		newRounds = s.createPairingPattern3(rankedPlayers, matchID, nextRoundNumber, totalCourts)
	}

	// 9. Save all rounds to database
	for i := range newRounds {
		if err := s.repo.Create(&newRounds[i]); err != nil {
			return nil, fmt.Errorf("failed to create pairing: %v", err)
		}
	}

	return newRounds, nil
}

// GET RANKED PLAYERS - Sort players by temp_score descending
func (s *matchRoundService) getRankedPlayers(players []models.MatchPlayer) []models.MatchPlayer {
	ranked := make([]models.MatchPlayer, len(players))
	copy(ranked, players)
	return ranked
}

// PATTERN 1: First round - Vertical pattern
// Multiple courts: Court 1: (1,3 vs 5,7), Court 2: (2,4 vs 6,8)
// Single court: (1,2 vs 3,4)
func (s *matchRoundService) createPairingPattern1(players []models.MatchPlayer, matchID string, roundNumber int, totalCourts int) []models.MatchRound {
	var rounds []models.MatchRound
	totalPlayers := len(players)

	if totalCourts == 1 {
		// Single court: 1,2 vs 3,4
		for i := 0; i < totalPlayers; i += 4 {
			if i+3 < totalPlayers {
				round := models.MatchRound{
					MatchID:        matchID,
					RoundNumber:    roundNumber,
					Court:          1,
					TeamAPlayer1ID: players[i].ID,
					TeamAPlayer2ID: players[i+1].ID,
					TeamBPlayer1ID: players[i+2].ID,
					TeamBPlayer2ID: players[i+3].ID,
					ScoreA:         0,
					ScoreB:         0,
				}
				rounds = append(rounds, round)
			}
		}
	} else {
		// Multiple courts: vertical pattern
		// Court 1: 1,3 vs 5,7
		// Court 2: 2,4 vs 6,8
		playersPerCourt := 4
		totalGroups := totalPlayers / (totalCourts * playersPerCourt)

		for group := 0; group < totalGroups; group++ {
			for court := 1; court <= totalCourts; court++ {
				baseIdx := group * (totalCourts * playersPerCourt)

				// Pattern: court 1 gets 1,3,5,7; court 2 gets 2,4,6,8
				idx1 := baseIdx + (court - 1)
				idx2 := baseIdx + (court - 1) + totalCourts
				idx3 := baseIdx + (court - 1) + totalCourts*2
				idx4 := baseIdx + (court - 1) + totalCourts*3

				if idx4 < totalPlayers {
					round := models.MatchRound{
						MatchID:        matchID,
						RoundNumber:    roundNumber,
						Court:          court,
						TeamAPlayer1ID: players[idx1].ID,
						TeamAPlayer2ID: players[idx2].ID,
						TeamBPlayer1ID: players[idx3].ID,
						TeamBPlayer2ID: players[idx4].ID,
						ScoreA:         0,
						ScoreB:         0,
					}
					rounds = append(rounds, round)
				}
			}
		}
	}

	return rounds
}

// PATTERN 2: Existing data with court <= 2 - Alternating vertical pattern
// Multiple courts: Court 1: (1,5 vs 3,7), Court 2: (2,6 vs 4,8)
// Single court: (1,3 vs 2,4)
func (s *matchRoundService) createPairingPattern2(players []models.MatchPlayer, matchID string, roundNumber int, totalCourts int) []models.MatchRound {
	var rounds []models.MatchRound
	totalPlayers := len(players)

	if totalCourts == 1 {
		// Single court: 1,3 vs 2,4
		for i := 0; i < totalPlayers; i += 4 {
			if i+3 < totalPlayers {
				round := models.MatchRound{
					MatchID:        matchID,
					RoundNumber:    roundNumber,
					Court:          1,
					TeamAPlayer1ID: players[i].ID,
					TeamAPlayer2ID: players[i+2].ID,
					TeamBPlayer1ID: players[i+1].ID,
					TeamBPlayer2ID: players[i+3].ID,
					ScoreA:         0,
					ScoreB:         0,
				}
				rounds = append(rounds, round)
			}
		}
	} else {
		// Multiple courts: alternating vertical pattern
		// Court 1: 1,5 vs 3,7
		// Court 2: 2,6 vs 4,8
		playersPerCourt := 4
		totalGroups := totalPlayers / (totalCourts * playersPerCourt)

		for group := 0; group < totalGroups; group++ {
			for court := 1; court <= totalCourts; court++ {
				baseIdx := group * (totalCourts * playersPerCourt)

				// Alternating pattern
				idx1 := baseIdx + (court - 1)
				idx2 := baseIdx + (court - 1) + totalCourts*2
				idx3 := baseIdx + (court - 1) + totalCourts
				idx4 := baseIdx + (court - 1) + totalCourts*3

				if idx4 < totalPlayers {
					round := models.MatchRound{
						MatchID:        matchID,
						RoundNumber:    roundNumber,
						Court:          court,
						TeamAPlayer1ID: players[idx1].ID,
						TeamAPlayer2ID: players[idx2].ID,
						TeamBPlayer1ID: players[idx3].ID,
						TeamBPlayer2ID: players[idx4].ID,
						ScoreA:         0,
						ScoreB:         0,
					}
					rounds = append(rounds, round)
				}
			}
		}
	}

	return rounds
}

// PATTERN 3: Existing data with court >= 3 - Horizontal pattern
// Multiple courts: Court 1: (1,2 vs 3,4), Court 2: (5,6 vs 7,8)
// Single court: (1,2 vs 3,4)
func (s *matchRoundService) createPairingPattern3(players []models.MatchPlayer, matchID string, roundNumber int, totalCourts int) []models.MatchRound {
	var rounds []models.MatchRound
	totalPlayers := len(players)

	// Horizontal pattern works the same for single or multiple courts
	court := 1
	for i := 0; i < totalPlayers; i += 4 {
		if i+3 < totalPlayers {
			round := models.MatchRound{
				MatchID:        matchID,
				RoundNumber:    roundNumber,
				Court:          court,
				TeamAPlayer1ID: players[i].ID,
				TeamAPlayer2ID: players[i+1].ID,
				TeamBPlayer1ID: players[i+2].ID,
				TeamBPlayer2ID: players[i+3].ID,
				ScoreA:         0,
				ScoreB:         0,
			}
			rounds = append(rounds, round)
			court++
			if court > totalCourts {
				court = 1
			}
		}
	}

	return rounds
}

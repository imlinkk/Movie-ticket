package services

import (
	"fmt"
	"log"
	"time"

	"movie-ticket/internal/models"
)

type EmailService interface {
	SendTicketConfirmation(user *models.User, ticket *models.Ticket, showtime *models.Showtime)
}

type emailService struct {
	taskQueue chan emailTask
}

type emailTask struct {
	User     *models.User
	Ticket   *models.Ticket
	Showtime *models.Showtime
}

func NewEmailService() EmailService {
	s := &emailService{
		taskQueue: make(chan emailTask, 100),
	}
	go s.startWorker()
	return s
}

func (s *emailService) startWorker() {
	log.Println("Email background worker started")
	for task := range s.taskQueue {
		// Simulate network / SMTP email delivery latency
		time.Sleep(300 * time.Millisecond)

		movieTitle := "Unknown Movie"
		if task.Showtime != nil && task.Showtime.Movie.Title != "" {
			movieTitle = task.Showtime.Movie.Title
		}

		log.Printf("[Email Notification] To: %s <%s> | Subject: Ticket Confirmation #%d | Movie: %s | Seat: %s | Price: $%.2f",
			task.User.Name,
			task.User.Email,
			task.Ticket.ID,
			movieTitle,
			task.Ticket.SeatNumber,
			task.Ticket.Price,
		)
		fmt.Printf("✓ Confirmation email successfully dispatched to %s for Ticket #%d\n", task.User.Email, task.Ticket.ID)
	}
}

func (s *emailService) SendTicketConfirmation(user *models.User, ticket *models.Ticket, showtime *models.Showtime) {
	if user == nil || ticket == nil {
		return
	}
	select {
	case s.taskQueue <- emailTask{User: user, Ticket: ticket, Showtime: showtime}:
	default:
		log.Printf("Warning: Email queue full, dropping email for ticket #%d", ticket.ID)
	}
}

package storage

import (
	"billohub/internal/model"
	"billohub/pkg/helper"
)

// CreateChat creates a new chat session for a specific user in the database.
func (s *Storage) CreateChat(chat *model.Chat) error {
	if err := s.db.Create(chat).Error; err != nil {
		return helper.WrapError(err, "failed to create new chat session")
	}
	return nil
}

// UpdateChatName updates the name of a chat session, ensuring the user has ownership.
func (s *Storage) UpdateChatName(chat *model.Chat) error {
	if err := s.db.Model(chat).Where("id = ? AND username = ?", chat.ID, chat.Username).Update("name", chat.Name).Error; err != nil {
		return helper.WrapError(err, "failed to modify chat session name")
	}
	return nil
}

// GetChats retrieves the list of chat sessions for a specific user from the database.
func (s *Storage) GetChats(username string) ([]model.Chat, error) {
	var chats []model.Chat
	if err := s.db.Where("username = ?", username).Order("created_at DESC").Find(&chats).Error; err != nil {
		return nil, helper.WrapError(err, "failed to get chat session list for user")
	}
	return chats, nil
}

func (s *Storage) DeleteChatById(chatId string) error {
	var chats model.Chat

	if err := s.db.Delete(chats, "id = ?", chatId).Error; err != nil {
		return helper.WrapError(err, "failed to delete chat for user")
	}
	return nil
}

package validators

import (
	dto "ROP_Backend/internal/dto/request"
	"errors"
)

func ValidateTagSkill(s dto.CreateTagSkill) error {
	if s.Name == "" {
		return errors.New("name required")
	}
	return nil
}

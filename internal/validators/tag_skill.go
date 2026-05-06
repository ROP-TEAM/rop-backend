package validators

import (
	"ROP_Backend/internal/dto"
	"errors"
)

func ValidateTagSkill(s dto.CreateTagSkill) error {
	if s.Name == "" {
		return errors.New("name required")
	}
	return nil
}

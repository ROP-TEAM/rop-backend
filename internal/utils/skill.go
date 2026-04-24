package utils

func IsSkillMatch(orderSkills, vehicleSkills []string) bool {
	set := make(map[string]bool)

	for _, skill := range vehicleSkills {
		set[skill] = true
	}
	for _, skill := range orderSkills {
		if !set[skill] {
			return false
		}
	}
	return true
}

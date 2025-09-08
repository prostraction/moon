package phase

import "testing"

func TestGetMoonPhase(t *testing.T) {
	tests := []struct {
		name          string
		before        float64
		current       float64
		after         float64
		lang          string
		expectedName  string
		expectedLocal string
		expectedEmoji string
	}{
		// Waxing Crescent
		{
			name:          "Waxing Crescent - normal case",
			before:        0.03,
			current:       0.1,
			after:         0.2,
			lang:          "en",
			expectedName:  "Waxing Crescent",
			expectedLocal: "Waxing Crescent",
			expectedEmoji: "🌒",
		},
		{
			name:          "Waxing Crescent - lower boundary",
			before:        0.04,
			current:       0.051,
			after:         0.3,
			lang:          "en",
			expectedName:  "Waxing Crescent",
			expectedLocal: "Waxing Crescent",
			expectedEmoji: "🌒",
		},
		{
			name:          "Waxing Crescent - upper boundary",
			before:        0.1,
			current:       0.449,
			after:         0.5,
			lang:          "en",
			expectedName:  "Waxing Crescent",
			expectedLocal: "Waxing Crescent",
			expectedEmoji: "🌒",
		},

		// First quarter
		{
			name:          "First quarter - exact lower boundary",
			before:        0.4,
			current:       0.45,
			after:         0.6,
			lang:          "en",
			expectedName:  "First quarter",
			expectedLocal: "First quarter",
			expectedEmoji: "🌓",
		},
		{
			name:          "First quarter - middle",
			before:        0.4,
			current:       0.5,
			after:         0.6,
			lang:          "en",
			expectedName:  "First quarter",
			expectedLocal: "First quarter",
			expectedEmoji: "🌓",
		},
		{
			name:          "First quarter - exact upper boundary",
			before:        0.4,
			current:       0.55,
			after:         0.6,
			lang:          "en",
			expectedName:  "First quarter",
			expectedLocal: "First quarter",
			expectedEmoji: "🌓",
		},

		// Waxing Gibbous
		{
			name:          "Waxing Gibbous - lower boundary",
			before:        0.5,
			current:       0.551,
			after:         0.7,
			lang:          "en",
			expectedName:  "Waxing Gibbous",
			expectedLocal: "Waxing Gibbous",
			expectedEmoji: "🌔",
		},
		{
			name:          "Waxing Gibbous - middle",
			before:        0.6,
			current:       0.7,
			after:         0.8,
			lang:          "en",
			expectedName:  "Waxing Gibbous",
			expectedLocal: "Waxing Gibbous",
			expectedEmoji: "🌔",
		},
		{
			name:          "Waxing Gibbous - upper boundary",
			before:        0.8,
			current:       0.949,
			after:         0.96,
			lang:          "en",
			expectedName:  "Waxing Gibbous",
			expectedLocal: "Waxing Gibbous",
			expectedEmoji: "🌔",
		},

		// Full Moon
		{
			name:          "Full Moon - exact lower boundary",
			before:        0.9,
			current:       0.95,
			after:         0.96,
			lang:          "en",
			expectedName:  "Full Moon",
			expectedLocal: "Full Moon",
			expectedEmoji: "🌕",
		},
		{
			name:          "Full Moon - above boundary",
			before:        0.9,
			current:       0.98,
			after:         0.99,
			lang:          "en",
			expectedName:  "Full Moon",
			expectedLocal: "Full Moon",
			expectedEmoji: "🌕",
		},
		{
			name:          "Full Moon - maximum",
			before:        0.9,
			current:       1.0,
			after:         0.99,
			lang:          "en",
			expectedName:  "Full Moon",
			expectedLocal: "Full Moon",
			expectedEmoji: "🌕",
		},

		// Waning Gibbous
		{
			name:          "Waning Gibbous - lower boundary",
			before:        0.96,
			current:       0.94,
			after:         0.9,
			lang:          "en",
			expectedName:  "Waning Gibbous",
			expectedLocal: "Waning Gibbous",
			expectedEmoji: "🌖",
		},
		{
			name:          "Waning Gibbous - middle",
			before:        0.8,
			current:       0.7,
			after:         0.6,
			lang:          "en",
			expectedName:  "Waning Gibbous",
			expectedLocal: "Waning Gibbous",
			expectedEmoji: "🌖",
		},
		{
			name:          "Waning Gibbous - upper boundary",
			before:        0.6,
			current:       0.551,
			after:         0.5,
			lang:          "en",
			expectedName:  "Waning Gibbous",
			expectedLocal: "Waning Gibbous",
			expectedEmoji: "🌖",
		},

		// Third quarter
		{
			name:          "Third quarter - exact upper boundary",
			before:        0.6,
			current:       0.55,
			after:         0.5,
			lang:          "en",
			expectedName:  "Third quarter",
			expectedLocal: "Third quarter",
			expectedEmoji: "🌗",
		},
		{
			name:          "Third quarter - middle",
			before:        0.6,
			current:       0.5,
			after:         0.4,
			lang:          "en",
			expectedName:  "Third quarter",
			expectedLocal: "Third quarter",
			expectedEmoji: "🌗",
		},
		{
			name:          "Third quarter - exact lower boundary",
			before:        0.6,
			current:       0.45,
			after:         0.4,
			lang:          "en",
			expectedName:  "Third quarter",
			expectedLocal: "Third quarter",
			expectedEmoji: "🌗",
		},

		// Waning Crescent
		{
			name:          "Waning Crescent - lower boundary",
			before:        0.5,
			current:       0.449,
			after:         0.4,
			lang:          "en",
			expectedName:  "Waning Crescent",
			expectedLocal: "Waning Crescent",
			expectedEmoji: "🌘",
		},
		{
			name:          "Waning Crescent - middle",
			before:        0.4,
			current:       0.3,
			after:         0.2,
			lang:          "en",
			expectedName:  "Waning Crescent",
			expectedLocal: "Waning Crescent",
			expectedEmoji: "🌘",
		},
		{
			name:          "Waning Crescent - upper boundary",
			before:        0.1,
			current:       0.051,
			after:         0.04,
			lang:          "en",
			expectedName:  "Waning Crescent",
			expectedLocal: "Waning Crescent",
			expectedEmoji: "🌘",
		},

		// New Moon
		{
			name:          "New Moon - exact boundary",
			before:        0.1,
			current:       0.05,
			after:         0.1,
			lang:          "en",
			expectedName:  "New Moon",
			expectedLocal: "New Moon",
			expectedEmoji: "🌑",
		},
		{
			name:          "New Moon - below boundary",
			before:        0.1,
			current:       0.04,
			after:         0.1,
			lang:          "en",
			expectedName:  "New Moon",
			expectedLocal: "New Moon",
			expectedEmoji: "🌑",
		},
		{
			name:          "New Moon - zero",
			before:        0.1,
			current:       0.0,
			after:         0.1,
			lang:          "en",
			expectedName:  "New Moon",
			expectedLocal: "New Moon",
			expectedEmoji: "🌑",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phaseResp := GetMoonPhase(tt.before, tt.current, tt.after, tt.lang)
			phaseName, phaseNameLocalized, phaseEmoji := phaseResp.Name, phaseResp.NameLocalized, phaseResp.Emoji

			if phaseName != tt.expectedName {
				t.Errorf("Expected phase name %s, got %s", tt.expectedName, phaseName)
			}

			if phaseNameLocalized != tt.expectedLocal {
				t.Errorf("Expected localized phase name %s, got %s", tt.expectedLocal, phaseNameLocalized)
			}

			if phaseEmoji != tt.expectedEmoji {
				t.Errorf("Expected emoji %s, got %s", tt.expectedEmoji, phaseEmoji)
			}
		})
	}
}

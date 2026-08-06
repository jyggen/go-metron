package metron_test

import (
	"testing"

	"github.com/jyggen/go-metron"
)

func TestCharacterByID(t *testing.T) {
	t.Parallel()
	testByID(t, "character", (*metron.Client).CharacterByID, []testCase[*metron.Character]{
		{
			id: 83,
			expected: &metron.Character{
				ID:    83,
				Name:  "Joker",
				Alias: new([]string{"Red Hood"}),
				Description: new(
					"The Joker is a fictional super-villain created by Bill Finger, Bob Kane, and Jerry Robinson who first appeared in the debut issue of the comic book Batman (April 25, 1940), published by DC Comics. Credit for the Joker's creation is disputed; Kane and Robinson claimed responsibility for the Joker's design, while acknowledging Finger's writing contribution. Although the Joker was planned to be killed off during his initial appearance, he was spared by editorial intervention, allowing the character to endure as the archenemy of the superhero Batman.\r\n\r\nIn his comic book appearances, the Joker is portrayed as a criminal mastermind. Introduced as a psychopath with a warped, sadistic sense of humor, the character became a goofy prankster in the late 1950's in response to regulation by the Comics Code Authority, before returning to his darker roots during the early 1970's. As Batman's nemesis, the Joker has been part of the superhero's defining stories, including the murder of Jason Todd—the second Robin and Batman's ward—and the paralysis of one of Batman's allies, Barbara Gordon. The Joker has had various possible origin stories during his decades of appearances. The most common story involves him falling into a tank of chemical waste which bleaches his skin white and turns his hair green and lips bright red; the resulting disfigurement drives him insane. The antithesis of Batman in personality and appearance, the Joker is considered by critics to be his perfect adversary.\r\n\r\nThe Joker possesses no superhuman abilities, instead using his expertise in chemical engineering to develop poisonous or lethal concoctions, and thematic weaponry, including razor-tipped playing cards, deadly joy buzzers, and acid-spraying lapel flowers. The Joker sometimes works with other Gotham City super-villains such as the Penguin and Two-Face, and groups like the Injustice Gang and Injustice League, but these relationships often collapse due to the Joker's desire for unbridled chaos. The 1990's introduced a romantic interest for the Joker in his former psychiatrist, Harley Quinn, who becomes his villainous sidekick. Although his primary obsession is Batman, the Joker has also fought other heroes including Superman and Wonder Woman.",
				),
				ImageURL: new(
					parseURL(t, "https://static.metron.cloud/media/character/2018/11/13/joker.jpg"),
				),
				Creators: []metron.CreatorList{
					{
						ID:       41,
						Name:     "Bill Finger",
						Modified: parseTime(t, "2025-02-18T10:21:42.567434-05:00"),
					},
					{
						ID:       42,
						Name:     "Bob Kane",
						Modified: parseTime(t, "2025-05-13T07:58:11.618339-04:00"),
					},
					{
						ID:       176,
						Name:     "Jerry Robinson",
						Modified: parseTime(t, "2025-02-18T12:21:01.316153-05:00"),
					},
				},
				Teams: []metron.TeamList{
					{
						ID:       42,
						Name:     "Injustice League",
						Modified: parseTime(t, "2026-04-10T09:48:04.386086-04:00"),
					},
				},
				Universes: []metron.UniverseList{
					{
						ID:       87,
						Name:     "Batman '66",
						Modified: parseTime(t, "2024-01-25T09:16:32.190698-05:00"),
					},
					{
						ID:       23,
						Name:     "Bombshells",
						Modified: parseTime(t, "2024-01-25T09:26:16.040978-05:00"),
					},
					{
						ID:       4,
						Name:     "Dark Multiverse: Batman Who Laughs",
						Modified: parseTime(t, "2024-01-25T10:00:22.865462-05:00"),
					},
					{
						ID:       8,
						Name:     "DC Animated Universe",
						Modified: parseTime(t, "2024-01-25T09:19:15.565244-05:00"),
					},
					{
						ID:       86,
						Name:     "DC vs. Vampires",
						Modified: parseTime(t, "2024-01-25T09:21:31.445093-05:00"),
					},
					{
						ID:       150,
						Name:     "Earth One",
						Modified: parseTime(t, "2024-12-10T12:12:23.216026-05:00"),
					},
					{
						ID:       51,
						Name:     "Injustice",
						Modified: parseTime(t, "2024-01-26T12:19:31.696959-05:00"),
					},
					{
						ID:       20,
						Name:     "New Frontier",
						Modified: parseTime(t, "2025-11-09T21:09:14.716552-05:00"),
					},
					{
						ID:       171,
						Name:     "Prime Earth",
						Modified: parseTime(t, "2025-10-26T22:25:13.152897-04:00"),
					},
				},
				ComicVineID: new(1702),
				ResourceURL: parseURL(t, "https://metron.cloud/character/joker/"),
				Modified:    parseTime(t, "2026-08-05T09:15:26.347117-04:00"),
			},
		},
		{
			id: 26153,
			expected: &metron.Character{
				ID:          26153,
				Name:        "176-301",
				Alias:       new([]string{}),
				Description: new(""),
				ImageURL:    nil,
				Creators:    []metron.CreatorList{},
				Teams:       []metron.TeamList{},
				Universes:   []metron.UniverseList{},
				ComicVineID: new(31838),
				ResourceURL: parseURL(t, "https://metron.cloud/character/176-301/"),
				Modified:    parseTime(t, "2024-02-24T13:07:30.187788-05:00"),
			},
		},
	})
}

func TestCharacters(t *testing.T) {
	t.Parallel()
	testList(t, "character", (*metron.Client).Characters, []testCase[*metron.CharacterList]{
		{
			id: 6211,
			expected: &metron.CharacterList{
				ID:       6211,
				Name:     "0101",
				Modified: parseTime(t, "2025-01-02T12:37:53.196726-05:00"),
			},
		},
		{
			id: 36241,
			expected: &metron.CharacterList{
				ID:       36241,
				Name:     "0-1X",
				Modified: parseTime(t, "2025-05-04T12:27:02.188492-04:00"),
			},
		},
		{
			id: 23822,
			expected: &metron.CharacterList{
				ID:       23822,
				Name:     "10th Muse",
				Modified: parseTime(t, "2023-12-02T11:32:20.584317-05:00"),
			},
		},
		{
			id: 42936,
			expected: &metron.CharacterList{
				ID:       42936,
				Name:     "1138",
				Modified: parseTime(t, "2026-01-09T14:09:26.193684-05:00"),
			},
		},
	})
}

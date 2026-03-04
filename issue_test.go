package metron_test

import (
	"net/url"
	"testing"

	"github.com/jyggen/go-metron"
)

func TestIssueByID(t *testing.T) {
	t.Parallel()
	testByID(t, "issue", (*metron.Client).IssueByID, []testCase[*metron.Issue]{
		{
			id: 112901,
			expected: &metron.Issue{
				ID: 112901,
				Publisher: metron.Reference{
					ID:   2,
					Name: "DC Comics",
				},
				Imprint: nil,
				Series: struct {
					ID        int
					Name      string
					SortName  string
					Volume    int
					YearBegan int
					Type      metron.Reference
					Genres    []metron.Reference
				}{
					ID:        7133,
					Name:      "Batman",
					SortName:  "Batman",
					Volume:    1,
					YearBegan: 2012,
					Type: metron.Reference{
						ID:   8,
						Name: "Hardcover",
					},
					Genres: []metron.Reference{
						{
							ID:   1,
							Name: "Adult",
						},
						{
							ID:   13,
							Name: "Crime",
						},
						{
							ID:   10,
							Name: "Super-Hero",
						},
					},
				},
				Number:            "1",
				AlternativeNumber: "",
				Title:             asReference("The Court of Owls"),
				Name: []string{
					"Knife Trick",
					"Trust Fall",
					"The Thirteenth Hour",
					"Face the Court, Part One",
					"Face the Court, Part Two",
					"Beneath the Glass",
					"The Talons Strike!",
				},
				CoverDate:            parseDate(t, "2012-07-01"),
				StoreDate:            asReference(parseDate(t, "2012-05-09")),
				FinalOrderCutoffDate: nil,
				Price:                asReference("24.99"),
				Rating: metron.Reference{
					ID:   1,
					Name: "Unknown",
				},
				SKU:       asReference("JAN120300"),
				ISBN:      asReference(""),
				UPC:       asReference(""),
				PageCount: asReference(176),
				Description: asReference(
					"Beware the Court of Owls, that watches all the time, ruling Gotham from a shadowed perch, behind granite and lime, they watch you at your hearth, they watch you in your bed. Speak not a whispered word of them, or they'll send the Talon for your head.\r\n\r\nBatman had heard the tales of Gotham City's Court of Owls. Meeting in the shadows and using the nocturnal bird of prey as their calling card, the members of this powerful cabal are the true rulers of Gotham.\r\n\r\nBut the Dark Knight dismissed the stories as rumors and old wives' tales. Gotham was his city.\r\n\r\nUntil now.\r\n\r\nA brutal assassin is sinking his razor-sharp talons into the city's best and brightest, as well as its most dangerous and deadly. If the dark legends are true, his masters are more powerful predators than the Batman could ever imagine—and their nests are everywhere...\r\n\r\nThe superstar team of writer SCOTT SNYDER (AMERICAN VAMPIRE, BATMAN: THE BLACK MIRROR) and GREG CAPULLO (SPAWN) proudly present BATMAN: THE COURT OF OWLS (collecting BATMAN #1-7)—a soaring saga of mystery and terror that became one of the year's best-selling comics.",
				),
				ImageURL: asReference(
					parseURL(
						t,
						"https://static.metron.cloud/media/issue/2024/02/27/e49a51e397ac4489a81785ad8098d296.jpg",
					),
				),
				CoverHash: asReference("841068ef7e313dec"),
				Arcs:      []metron.ArcList{},
				Credits: []struct {
					ID    int
					Name  string
					Roles []metron.Reference
				}{
					{
						ID:   1379,
						Name: "Bob Harras",
						Roles: []metron.Reference{
							{
								ID:   20,
								Name: "Editor In Chief",
							},
						},
					},
					{
						ID:   1148,
						Name: "Dan DiDio",
						Roles: []metron.Reference{
							{
								ID:   31,
								Name: "Publisher",
							},
						},
					},
					{
						ID:   739,
						Name: "Eddie Berganza",
						Roles: []metron.Reference{
							{
								ID:   25,
								Name: "Executive Editor",
							},
						},
					},
					{
						ID:   102,
						Name: "FCO Plascencia",
						Roles: []metron.Reference{
							{
								ID:   5,
								Name: "Colorist",
							},
							{
								ID:   7,
								Name: "Cover",
							},
						},
					},
					{
						ID:   92,
						Name: "Geoff Johns",
						Roles: []metron.Reference{
							{
								ID:   18,
								Name: "Chief Creative Officer",
							},
						},
					},
					{
						ID:   103,
						Name: "Greg Capullo",
						Roles: []metron.Reference{
							{
								ID:   3,
								Name: "Penciller",
							},
							{
								ID:   7,
								Name: "Cover",
							},
						},
					},

					{
						ID:   628,
						Name: "Harvey Richards",
						Roles: []metron.Reference{
							{
								ID:   9,
								Name: "Associate Editor",
							},
						},
					},
					{
						ID:   681,
						Name: "Janelle Asselin",
						Roles: []metron.Reference{
							{
								ID:   9,
								Name: "Associate Editor",
							},
						},
					},
					{
						ID:   172,
						Name: "Jim Lee",
						Roles: []metron.Reference{
							{
								ID:   31,
								Name: "Publisher",
							},
						},
					},
					{
						ID:   1919,
						Name: "Jimmy Betancourt",
						Roles: []metron.Reference{
							{
								ID:   6,
								Name: "Letterer",
							},
						},
					},
					{
						ID:   104,
						Name: "Jonathan Glapion",
						Roles: []metron.Reference{
							{
								ID:   4,
								Name: "Inker",
							},
						},
					},
					{
						ID:   201,
						Name: "Katie Kubert",
						Roles: []metron.Reference{
							{
								ID:   12,
								Name: "Assistant Editor",
							},
						},
					},

					{
						ID:   309,
						Name: "Mike Marts",
						Roles: []metron.Reference{
							{
								ID:   8,
								Name: "Editor",
							},
						},
					},
					{
						ID:   9428,
						Name: "Peter Hamboussi",
						Roles: []metron.Reference{
							{
								ID:   8,
								Name: "Editor",
							},
						},
					},
					{
						ID:   206,
						Name: "Richard Starkings",
						Roles: []metron.Reference{
							{
								ID:   6,
								Name: "Letterer",
							},
						},
					},
					{
						ID:   3454,
						Name: "Robbie Biederman",
						Roles: []metron.Reference{
							{
								ID:   15,
								Name: "Designer",
							},
						},
					},
					{
						ID:   148,
						Name: "Scott Snyder",
						Roles: []metron.Reference{
							{
								ID:   1,
								Name: "Writer",
							},
						},
					},
				},
				Characters: []metron.CharacterList{
					{
						ID:       261,
						Name:     "Alfred Pennyworth",
						Modified: parseTime(t, "2025-05-10T11:41:36.028997-04:00"),
					},
					{
						ID:       275,
						Name:     "Barbara Gordon",
						Modified: parseTime(t, "2025-02-18T15:21:57.440821-05:00"),
					},
					{
						ID:       12,
						Name:     "Batman",
						Modified: parseTime(t, "2025-05-10T11:32:34.186924-04:00"),
					},
					{
						ID:       2351,
						Name:     "Bluebird",
						Modified: parseTime(t, "2025-05-10T11:43:21.126838-04:00"),
					},
					{
						ID:       77,
						Name:     "Catwoman",
						Modified: parseTime(t, "2025-02-18T15:22:10.535267-05:00"),
					},
					{
						ID:       1120,
						Name:     "Clayface (Karlo)",
						Modified: parseTime(t, "2025-02-18T15:27:24.466997-05:00"),
					},
					{
						ID:       394,
						Name:     "Damian Wayne",
						Modified: parseTime(t, "2025-02-18T15:28:22.696013-05:00"),
					},
					{
						ID:       293,
						Name:     "Dick Grayson",
						Modified: parseTime(t, "2025-02-18T15:28:26.982215-05:00"),
					},
					{
						ID:       16,
						Name:     "Harley Quinn",
						Modified: parseTime(t, "2025-02-18T15:22:26.752300-05:00"),
					},
					{
						ID:       1048,
						Name:     "Harvey Bullock",
						Modified: parseTime(t, "2025-02-18T15:22:32.371011-05:00"),
					},
					{
						ID:       82,
						Name:     "James Gordon",
						Modified: parseTime(t, "2025-02-18T15:22:43.443874-05:00"),
					},
					{
						ID:       1204,
						Name:     "James Gordon Jr.",
						Modified: parseTime(t, "2025-02-18T15:26:16.097968-05:00"),
					},
					{
						ID:       1126,
						Name:     "Jeremiah Arkham",
						Modified: parseTime(t, "2025-02-18T22:11:56.059808-05:00"),
					},
					{
						ID:       83,
						Name:     "Joker",
						Modified: parseTime(t, "2025-02-18T15:23:40.393641-05:00"),
					},
					{
						ID:       763,
						Name:     "Killer Croc",
						Modified: parseTime(t, "2025-02-18T15:26:43.639475-05:00"),
					},
					{
						ID:       2234,
						Name:     "Leslie Thompkins",
						Modified: parseTime(t, "2025-02-18T15:26:55.636901-05:00"),
					},
					{
						ID:       2856,
						Name:     "Lincoln March",
						Modified: parseTime(t, "2025-02-18T15:24:18.890445-05:00"),
					},
					{
						ID:       274,
						Name:     "Mr. Freeze",
						Modified: parseTime(t, "2025-02-18T15:27:29.919010-05:00"),
					},
					{
						ID:       347,
						Name:     "Professor Pyg",
						Modified: parseTime(t, "2025-02-18T15:22:57.025190-05:00"),
					},
					{
						ID:       348,
						Name:     "Riddler",
						Modified: parseTime(t, "2025-02-18T15:23:00.534192-05:00"),
					},
					{
						ID:       280,
						Name:     "Scarecrow",
						Modified: parseTime(t, "2025-02-18T15:23:04.327602-05:00"),
					},
					{
						ID:       765,
						Name:     "Tim Drake",
						Modified: parseTime(t, "2025-02-18T15:23:13.140852-05:00"),
					},
					{
						ID:       355,
						Name:     "Two-Face",
						Modified: parseTime(t, "2025-02-18T15:46:12.108514-05:00"),
					},
					{
						ID:       2057,
						Name:     "Vicki Vale",
						Modified: parseTime(t, "2025-02-18T15:24:02.317133-05:00"),
					},
					{
						ID:       1136,
						Name:     "Victor Zsasz",
						Modified: parseTime(t, "2025-02-18T15:27:34.641889-05:00"),
					},
				},
				Teams: []metron.TeamList{
					{
						ID:       87,
						Name:     "Court of Owls",
						Modified: parseTime(t, "2025-02-19T01:10:09.195823-05:00"),
					},
					{
						ID:       88,
						Name:     "The Talons",
						Modified: parseTime(t, "2019-06-23T15:13:24.036247-04:00"),
					},
				},
				Universes: []metron.UniverseList{},
				Reprints: []struct {
					ID    int
					Issue string
				}{
					{
						ID:    6798,
						Issue: "Batman (2011) #1",
					},
					{
						ID:    6799,
						Issue: "Batman (2011) #2",
					},
					{
						ID:    6800,
						Issue: "Batman (2011) #3",
					},
					{
						ID:    6801,
						Issue: "Batman (2011) #4",
					},
					{
						ID:    6802,
						Issue: "Batman (2011) #5",
					},
					{
						ID:    6803,
						Issue: "Batman (2011) #6",
					},
					{
						ID:    6804,
						Issue: "Batman (2011) #7",
					},
				},
				Variants: []struct {
					Name     *string
					SKU      *string
					UPC      *string
					Price    *string
					ImageURL url.URL
				}{},
				ComicVineID:           nil,
				GrandComicsDatabaseID: asReference(1035895),
				ResourceURL:           parseURL(t, "https://metron.cloud/issue/batman-2012-1/"),
				Modified:              parseTime(t, "2025-01-07T08:09:42.476138-05:00"),
			},
		},
		{
			id: 2558,
			expected: &metron.Issue{
				ID: 2558,
				Publisher: metron.Reference{
					ID:   2,
					Name: "DC Comics",
				},
				Imprint: &metron.Reference{
					ID:   2,
					Name: "DC Black Label",
				},
				Series: struct {
					ID        int
					Name      string
					SortName  string
					Volume    int
					YearBegan int
					Type      metron.Reference
					Genres    []metron.Reference
				}{
					ID:        279,
					Name:      "Batman: Last Knight on Earth",
					SortName:  "Batman Last Knight on Earth",
					Volume:    1,
					YearBegan: 2019,
					Type: metron.Reference{
						ID:   11,
						Name: "Limited Series",
					},
					Genres: []metron.Reference{
						{
							ID:   10,
							Name: "Super-Hero",
						},
					},
				},
				Number: "1",
				Title:  asReference(""),
				Name: []string{
					"Book One",
				},
				CoverDate: parseDate(t, "2019-07-01"),
				StoreDate: asReference(parseDate(t, "2019-05-29")),
				Price:     asReference("5.99"),
				Rating: metron.Reference{
					ID:   5,
					Name: "Mature",
				},
				SKU:       asReference(""),
				ISBN:      asReference(""),
				UPC:       asReference("76194135390600111"),
				PageCount: asReference(60),
				Description: asReference(
					"Bruce Wayne wakes up in Arkham Asylum. Young. Sane.\r\nAnd...he was never Batman.\r\n\r\n So begins this sprawling tale of the Dark Knight as he embarks on a quest through a devastated DC landscape featuring a massive cast of familiar faces from the DC Universe. As he tries to piece together the mystery of his past, he must unravel the cause of this terrible future and track down the unspeakable force that destroyed the world as he knew it…",
				),
				ImageURL: asReference(
					parseURL(t, "https://static.metron.cloud/media/issue/2019/05/27/batman-last-knight-1.jpg"),
				),
				CoverHash: asReference("ede81312b2337ac5"),
				Arcs:      []metron.ArcList{},
				Credits: []struct {
					ID    int
					Name  string
					Roles []metron.Reference
				}{
					{
						ID:   303,
						Name: "Amedeo Turturro",
						Roles: []metron.Reference{
							{
								ID:   9,
								Name: "Associate Editor",
							},
						},
					},
					{
						ID:   102,
						Name: "FCO Plascencia",
						Roles: []metron.Reference{
							{
								ID:   5,
								Name: "Colorist",
							},
							{
								ID:   7,
								Name: "Cover",
							},
						},
					},
					{
						ID:   103,
						Name: "Greg Capullo",
						Roles: []metron.Reference{
							{
								ID:   3,
								Name: "Penciller",
							},
							{
								ID:   7,
								Name: "Cover",
							},
						},
					},
					{
						ID:   740,
						Name: "Jock",
						Roles: []metron.Reference{
							{
								ID:   7,
								Name: "Cover",
							},
						},
					},
					{
						ID:   104,
						Name: "Jonathan Glapion",
						Roles: []metron.Reference{
							{
								ID:   4,
								Name: "Inker",
							},
						},
					},
					{
						ID:   173,
						Name: "Mark Doyle",
						Roles: []metron.Reference{
							{
								ID:   8,
								Name: "Editor",
							},
						},
					},

					{
						ID:   148,
						Name: "Scott Snyder",
						Roles: []metron.Reference{
							{
								ID:   1,
								Name: "Writer",
							},
						},
					},
					{
						ID:   62,
						Name: "Tom Napolitano",
						Roles: []metron.Reference{
							{
								ID:   6,
								Name: "Letterer",
							},
						},
					},
				},
				Characters: []metron.CharacterList{
					{
						ID:       261,
						Name:     "Alfred Pennyworth",
						Modified: parseTime(t, "2025-05-10T11:41:36.028997-04:00"),
					},

					{
						ID:       12,
						Name:     "Batman",
						Modified: parseTime(t, "2025-05-10T11:32:34.186924-04:00"),
					},
					{
						ID:       297,
						Name:     "Huntress (Bertinelli)",
						Modified: parseTime(t, "2025-02-18T15:46:25.289114-05:00"),
					},
					{
						ID:       83,
						Name:     "Joker",
						Modified: parseTime(t, "2025-02-18T15:23:40.393641-05:00"),
					},
					{
						ID:       22,
						Name:     "Poison Ivy",
						Modified: parseTime(t, "2025-02-18T15:27:40.557163-05:00"),
					},
					{
						ID:       103,
						Name:     "Supergirl (Kara Zor-El)",
						Modified: parseTime(t, "2025-02-18T15:41:12.282439-05:00"),
					},
					{
						ID:       14,
						Name:     "Wonder Woman",
						Modified: parseTime(t, "2025-02-18T15:41:16.582787-05:00"),
					},
				},
				Teams:     []metron.TeamList{},
				Universes: []metron.UniverseList{},
				Reprints: []struct {
					ID    int
					Issue string
				}{},
				Variants: []struct {
					Name     *string
					SKU      *string
					UPC      *string
					Price    *string
					ImageURL url.URL
				}{
					{
						Name:  asReference("Variant Cover"),
						SKU:   asReference(""),
						UPC:   asReference(""),
						Price: nil,
						ImageURL: parseURL(
							t,
							"https://static.metron.cloud/media/variants/2019/05/27/batman-last-knight-1a.jpg",
						),
					},
				},
				ComicVineID:           asReference(710046),
				GrandComicsDatabaseID: asReference(1970616),
				ResourceURL:           parseURL(t, "https://metron.cloud/issue/batman-last-knight-earth-2019-1/"),
				Modified:              parseTime(t, "2024-12-21T14:45:34.863850-05:00"),
			},
		},
	})
}

func TestIssuesByArcID(t *testing.T) {
	t.Parallel()
	testListByID(t, "arc", 659, "issue", (*metron.Client).IssuesByArcID, issueListTestCases(t))
}

func TestIssuesByCharacterID(t *testing.T) {
	t.Parallel()
	testListByID(t, "character", 83, "issue", (*metron.Client).IssuesByCharacterID, issueListTestCases(t))
}

func TestIssuesBySeriesID(t *testing.T) {
	t.Parallel()
	testListByID(t, "series", 793, "issue", (*metron.Client).IssuesBySeriesID, issueListTestCases(t))
}

func TestIssuesByTeamID(t *testing.T) {
	t.Parallel()
	testListByID(t, "team", 180, "issue", (*metron.Client).IssuesByTeamID, issueListTestCases(t))
}

func TestIssues(t *testing.T) {
	t.Parallel()
	testList(t, "issue", (*metron.Client).Issues, issueListTestCases(t))
}

func issueListTestCases(t *testing.T) []testCase[*metron.IssueList] {
	return []testCase[*metron.IssueList]{
		{
			id: 89088,
			expected: &metron.IssueList{
				ID: 89088,
				Series: struct {
					Name      string
					Volume    int
					YearBegan int
				}{
					Name:      "'68",
					Volume:    1,
					YearBegan: 2006,
				},
				Number:    "1",
				Name:      "'68 (2006) #1",
				CoverDate: parseDate(t, "2006-12-01"),
				StoreDate: nil,
				ImageURL: asReference(
					parseURL(
						t,
						"https://static.metron.cloud/media/issue/2023/12/02/395d1af45859497183f0b868485831aa.jpg",
					),
				),
				CoverHash: asReference("de9a207b768189f4"),
				Modified:  parseTime(t, "2023-12-02T10:08:26.125751-05:00"),
			},
		},
		{
			id: 89089,
			expected: &metron.IssueList{
				ID: 89089,
				Series: struct {
					Name      string
					Volume    int
					YearBegan int
				}{
					Name:      "'68",
					Volume:    2,
					YearBegan: 2011,
				},
				Number:    "1",
				Name:      "'68 (2011) #1",
				CoverDate: parseDate(t, "2011-04-01"),
				StoreDate: asReference(parseDate(t, "2011-04-27")),
				ImageURL: asReference(
					parseURL(
						t,
						"https://static.metron.cloud/media/issue/2023/12/02/b8a86c0422e749968f74c2fe2605ea4a.jpg",
					),
				),
				CoverHash: asReference("ccd6331a22959d3b"),
				Modified:  parseTime(t, "2023-12-02T10:10:55.924868-05:00"),
			},
		},
		{
			id: 89090,
			expected: &metron.IssueList{
				ID: 89090,
				Series: struct {
					Name      string
					Volume    int
					YearBegan int
				}{
					Name:      "'68",
					Volume:    2,
					YearBegan: 2011,
				},
				Number:    "2",
				Name:      "'68 (2011) #2",
				CoverDate: parseDate(t, "2011-05-01"),
				StoreDate: asReference(parseDate(t, "2011-05-25")),
				ImageURL: asReference(
					parseURL(
						t,
						"https://static.metron.cloud/media/issue/2023/12/02/ceddeae0aac84df7a509d66003c8007c.jpg",
					),
				),
				CoverHash: asReference("c0471e7273cf1ccc"),
				Modified:  parseTime(t, "2023-12-02T10:11:19.577639-05:00"),
			},
		},
		{
			id: 89091,
			expected: &metron.IssueList{
				ID: 89091,
				Series: struct {
					Name      string
					Volume    int
					YearBegan int
				}{
					Name:      "'68",
					Volume:    2,
					YearBegan: 2011,
				},
				Number:    "3",
				Name:      "'68 (2011) #3",
				CoverDate: parseDate(t, "2011-08-01"),
				StoreDate: asReference(parseDate(t, "2011-08-03")),
				ImageURL: asReference(
					parseURL(
						t,
						"https://static.metron.cloud/media/issue/2023/12/02/fb0072f3b7094758a3488f4cc8b6e551.jpg",
					),
				),
				CoverHash: asReference("cc97960b5be6a684"),
				Modified:  parseTime(t, "2023-12-02T10:11:24.329496-05:00"),
			},
		},
	}
}

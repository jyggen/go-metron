package metron_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jyggen/go-metron"
	"github.com/stretchr/testify/require"
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
				Series: metron.IssueSeries{
					ID:               7133,
					Name:             "Batman",
					AlternativeNames: []string{},
					SortName:         "Batman",
					Volume:           1,
					YearBegan:        2012,
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
				Title:             new("The Court of Owls"),
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
				StoreDate:            new(parseDate(t, "2012-05-09")),
				FinalOrderCutoffDate: nil,
				Price:                new("24.99"),
				PriceCurrency:        "USD",
				Rating: metron.Reference{
					ID:   1,
					Name: "Unknown",
				},
				SKU:       new("JAN120300"),
				ISBN:      new(""),
				UPC:       new(""),
				PageCount: new(176),
				Description: new(
					"Beware the Court of Owls, that watches all the time, ruling Gotham from a shadowed perch, behind granite and lime, they watch you at your hearth, they watch you in your bed. Speak not a whispered word of them, or they'll send the Talon for your head.\r\n\r\nBatman had heard the tales of Gotham City's Court of Owls. Meeting in the shadows and using the nocturnal bird of prey as their calling card, the members of this powerful cabal are the true rulers of Gotham.\r\n\r\nBut the Dark Knight dismissed the stories as rumors and old wives' tales. Gotham was his city.\r\n\r\nUntil now.\r\n\r\nA brutal assassin is sinking his razor-sharp talons into the city's best and brightest, as well as its most dangerous and deadly. If the dark legends are true, his masters are more powerful predators than the Batman could ever imagine—and their nests are everywhere...\r\n\r\nThe superstar team of writer SCOTT SNYDER (AMERICAN VAMPIRE, BATMAN: THE BLACK MIRROR) and GREG CAPULLO (SPAWN) proudly present BATMAN: THE COURT OF OWLS (collecting BATMAN #1-7)—a soaring saga of mystery and terror that became one of the year's best-selling comics.",
				),
				ImageURL: new(
					parseURL(
						t,
						"https://static.metron.cloud/media/issue/2024/02/27/e49a51e397ac4489a81785ad8098d296.jpg",
					),
				),
				CoverHash:     new("841068ef7e313dec"),
				AverageRating: nil,
				RatingCount:   0,
				Arcs:          []metron.ArcList{},
				Credits: []metron.IssueCredit{
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
						Modified: parseTime(t, "2026-08-02T15:12:29.254628-04:00"),
					},
					{
						ID:       275,
						Name:     "Barbara Gordon",
						Modified: parseTime(t, "2026-08-02T15:07:50.635149-04:00"),
					},
					{
						ID:       12,
						Name:     "Batman",
						Modified: parseTime(t, "2026-08-05T09:39:56.696076-04:00"),
					},
					{
						ID:       2351,
						Name:     "Bluebird",
						Modified: parseTime(t, "2025-11-09T11:31:37.118434-05:00"),
					},
					{
						ID:       77,
						Name:     "Catwoman",
						Modified: parseTime(t, "2026-08-05T09:28:49.406962-04:00"),
					},
					{
						ID:       1120,
						Name:     "Clayface (Karlo)",
						Modified: parseTime(t, "2026-06-14T16:58:12.246869-04:00"),
					},
					{
						ID:       394,
						Name:     "Damian Wayne",
						Modified: parseTime(t, "2026-08-05T09:12:31.176305-04:00"),
					},
					{
						ID:       293,
						Name:     "Dick Grayson",
						Modified: parseTime(t, "2026-08-02T15:12:29.254628-04:00"),
					},
					{
						ID:       16,
						Name:     "Harley Quinn",
						Modified: parseTime(t, "2026-08-01T14:47:53.412361-04:00"),
					},
					{
						ID:       1048,
						Name:     "Harvey Bullock",
						Modified: parseTime(t, "2026-08-02T15:11:22.944954-04:00"),
					},
					{
						ID:       82,
						Name:     "James Gordon",
						Modified: parseTime(t, "2026-08-05T09:07:13.282845-04:00"),
					},
					{
						ID:       1204,
						Name:     "James Gordon Jr.",
						Modified: parseTime(t, "2026-03-13T13:25:33.554299-04:00"),
					},
					{
						ID:       1126,
						Name:     "Jeremiah Arkham",
						Modified: parseTime(t, "2025-11-09T12:30:42.840247-05:00"),
					},
					{
						ID:       83,
						Name:     "Joker",
						Modified: parseTime(t, "2026-08-05T09:15:26.347117-04:00"),
					},
					{
						ID:       763,
						Name:     "Killer Croc",
						Modified: parseTime(t, "2026-06-22T09:16:05.969582-04:00"),
					},
					{
						ID:       2234,
						Name:     "Leslie Thompkins",
						Modified: parseTime(t, "2026-08-02T15:11:22.944954-04:00"),
					},
					{
						ID:       2856,
						Name:     "Lincoln March",
						Modified: parseTime(t, "2025-11-09T10:54:18.059163-05:00"),
					},
					{
						ID:       274,
						Name:     "Mr. Freeze",
						Modified: parseTime(t, "2026-06-22T09:17:05.769693-04:00"),
					},
					{
						ID:       347,
						Name:     "Professor Pyg",
						Modified: parseTime(t, "2026-08-01T02:10:53.299159-04:00"),
					},
					{
						ID:       348,
						Name:     "Riddler",
						Modified: parseTime(t, "2026-08-02T15:00:36.769818-04:00"),
					},
					{
						ID:       280,
						Name:     "Scarecrow (DC)",
						Modified: parseTime(t, "2026-08-05T09:39:56.696076-04:00"),
					},
					{
						ID:       765,
						Name:     "Tim Drake",
						Modified: parseTime(t, "2026-08-01T14:44:21.729825-04:00"),
					},
					{
						ID:       355,
						Name:     "Two-Face",
						Modified: parseTime(t, "2026-08-02T15:00:36.769818-04:00"),
					},
					{
						ID:       2057,
						Name:     "Vicki Vale",
						Modified: parseTime(t, "2026-08-02T15:07:09.562963-04:00"),
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
						Modified: parseTime(t, "2025-11-09T10:47:58.678221-05:00"),
					},
				},
				Universes: []metron.UniverseList{},
				Reprints: []metron.IssueReprint{
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
				Variants:              []metron.IssueVariant{},
				ComicVineID:           nil,
				GrandComicsDatabaseID: new(1035895),
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
				Series: metron.IssueSeries{
					ID:               279,
					Name:             "Batman: Last Knight on Earth",
					AlternativeNames: []string{},
					SortName:         "Batman Last Knight on Earth",
					Volume:           1,
					YearBegan:        2019,
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
				Title:  new(""),
				Name: []string{
					"Book One",
				},
				CoverDate:     parseDate(t, "2019-07-01"),
				StoreDate:     new(parseDate(t, "2019-05-29")),
				Price:         new("5.99"),
				PriceCurrency: "USD",
				Rating: metron.Reference{
					ID:   5,
					Name: "Mature",
				},
				SKU:       new(""),
				ISBN:      new(""),
				UPC:       new("76194135390600111"),
				PageCount: new(60),
				Description: new(
					"Bruce Wayne wakes up in Arkham Asylum. Young. Sane.\r\nAnd...he was never Batman.\r\n\r\n So begins this sprawling tale of the Dark Knight as he embarks on a quest through a devastated DC landscape featuring a massive cast of familiar faces from the DC Universe. As he tries to piece together the mystery of his past, he must unravel the cause of this terrible future and track down the unspeakable force that destroyed the world as he knew it…",
				),
				ImageURL: new(
					parseURL(t, "https://static.metron.cloud/media/issue/2019/05/27/batman-last-knight-1.jpg"),
				),
				CoverHash:     new("ede81312b2337ac5"),
				AverageRating: nil,
				RatingCount:   0,
				Arcs:          []metron.ArcList{},
				Credits: []metron.IssueCredit{
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
						Modified: parseTime(t, "2026-08-02T15:12:29.254628-04:00"),
					},

					{
						ID:       12,
						Name:     "Batman",
						Modified: parseTime(t, "2026-08-05T09:39:56.696076-04:00"),
					},
					{
						ID:       297,
						Name:     "Huntress (Bertinelli)",
						Modified: parseTime(t, "2026-07-08T09:49:44.287660-04:00"),
					},
					{
						ID:       83,
						Name:     "Joker",
						Modified: parseTime(t, "2026-08-05T09:15:26.347117-04:00"),
					},
					{
						ID:       22,
						Name:     "Poison Ivy",
						Modified: parseTime(t, "2026-06-24T11:38:36.039298-04:00"),
					},
					{
						ID:       103,
						Name:     "Supergirl (Kara Zor-El)",
						Modified: parseTime(t, "2026-08-05T11:09:32.424517-04:00"),
					},
					{
						ID:       14,
						Name:     "Wonder Woman",
						Modified: parseTime(t, "2026-08-03T09:03:51.967056-04:00"),
					},
				},
				Teams:     []metron.TeamList{},
				Universes: []metron.UniverseList{},
				Reprints:  []metron.IssueReprint{},
				Variants: []metron.IssueVariant{
					{
						Name:  new("Variant Cover"),
						SKU:   new(""),
						UPC:   new(""),
						Price: nil,
						ImageURL: parseURL(
							t,
							"https://static.metron.cloud/media/variants/2019/05/27/batman-last-knight-1a.jpg",
						),
					},
				},
				ComicVineID:           new(710046),
				GrandComicsDatabaseID: new(1970616),
				ResourceURL:           parseURL(t, "https://metron.cloud/issue/batman-last-knight-earth-2019-1/"),
				Modified:              parseTime(t, "2024-12-21T14:45:34.863850-05:00"),
			},
		},
	})
}

func TestIssuesByArcID(t *testing.T) {
	t.Parallel()
	testListByID(t, "arc", 659, "issue", (*metron.Client).IssuesByArcID, issueListTestCases(t, "arc"))
}

func TestIssuesByCharacterID(t *testing.T) {
	t.Parallel()
	testListByID(t, "character", 83, "issue", (*metron.Client).IssuesByCharacterID, issueListTestCases(t, "character"))
}

func TestIssuesBySeriesID(t *testing.T) {
	t.Parallel()
	testListByID(t, "series", 793, "issue", (*metron.Client).IssuesBySeriesID, issueListTestCases(t, "series"))
}

func TestIssuesByTeamID(t *testing.T) {
	t.Parallel()
	testListByID(t, "team", 180, "issue", (*metron.Client).IssuesByTeamID, issueListTestCases(t, "team"))
}

func TestIssues(t *testing.T) {
	t.Parallel()
	testList(t, "issue", (*metron.Client).Issues, issueListTestCases(t, "issue"))
}

// issueListRecord is one expected IssueList in flattened form. The five
// issue-list endpoints each return different records, and spelling every one
// out as a struct literal repeated the same fifteen lines twenty times. An
// empty string means the field is absent upstream, so the pointer stays nil.
type issueListRecord struct {
	kind      string
	id        int
	seriesID  int
	series    string
	volume    int
	yearBegan int
	name      string
	number    string
	coverDate string
	storeDate string
	image     string
	coverHash string
	modified  string
}

// issueListRecords holds the expectations for every endpoint returning an
// issue list. One flat table rather than a map of per-kind groups: each
// endpoint returns different records, but the shape is identical.
var issueListRecords = []issueListRecord{
	{kind: "issue", id: 52529, seriesID: 3856, series: "The 06 Protocol", volume: 1, yearBegan: 2022, name: "The 06 Protocol (2022) #1", number: "1", coverDate: "2022-09-01", storeDate: "2022-09-14", image: "https://static.metron.cloud/media/issue/2022/09/19/6-protocol-1.jpg", coverHash: "e9529e3586ca6b34", modified: "2024-12-23T16:17:29.333263-05:00"},
	{kind: "issue", id: 52530, seriesID: 3856, series: "The 06 Protocol", volume: 1, yearBegan: 2022, name: "The 06 Protocol (2022) #2", number: "2", coverDate: "2022-11-01", storeDate: "2022-11-23", image: "https://static.metron.cloud/media/issue/2022/09/19/6-protocol-2.jpg", coverHash: "c3a90fdd304b66e8", modified: "2024-12-23T16:17:31.360018-05:00"},
	{kind: "issue", id: 61168, seriesID: 3856, series: "The 06 Protocol", volume: 1, yearBegan: 2022, name: "The 06 Protocol (2022) #3", number: "3", coverDate: "2023-06-01", storeDate: "2023-06-28", image: "https://static.metron.cloud/media/issue/2023/04/07/cfe56a98a0cd40bdbc5a4ed22956ca75.jpg", coverHash: "90944f7344c87f3e", modified: "2024-12-23T16:17:33.770796-05:00"},
	{kind: "issue", id: 136787, seriesID: 10307, series: "8-Bit Zombie", volume: 1, yearBegan: 2013, name: "8-Bit Zombie (2013) #1", number: "1", coverDate: "2013-10-01", storeDate: "2013-11-27", image: "https://static.metron.cloud/media/issue/2025/02/25/Cover-A_AcSDvLN.jpg", coverHash: "ae8379d4aa671478", modified: "2025-02-25T22:34:04.333183-05:00"},
	{kind: "arc", id: 48264, seriesID: 3408, series: "Deadpool", volume: 6, yearBegan: 2015, name: "Deadpool (2015) #28", number: "28", coverDate: "2017-05-01", storeDate: "2017-03-01", image: "https://static.metron.cloud/media/issue/2022/05/21/deadpool-28.jpg", coverHash: "8dc8b6b2628c39b7", modified: "2024-12-29T13:12:34.104918-05:00"},
	{kind: "arc", id: 48128, seriesID: 3399, series: "Deadpool & The Mercs for Money", volume: 2, yearBegan: 2016, name: "Deadpool & The Mercs for Money (2016) #9", number: "9", coverDate: "2017-05-01", storeDate: "2017-03-29", image: "https://static.metron.cloud/media/issue/2022/05/17/deadpool-the-mercs-for-money-9.jpg", coverHash: "eab9c05bc5568f12", modified: "2024-12-29T13:18:45.506840-05:00"},
	{kind: "arc", id: 43011, seriesID: 3094, series: "Spider-Man / Deadpool", volume: 1, yearBegan: 2016, name: "Spider-Man / Deadpool (2016) #15", number: "15", coverDate: "2017-05-01", storeDate: "2017-03-08", image: "https://static.metron.cloud/media/issue/2022/01/16/spider-man-deadpool-15.jpg", coverHash: "ec94bcb6416a4acb", modified: "2025-01-12T18:12:50.219180-05:00"},
	{kind: "arc", id: 48265, seriesID: 3408, series: "Deadpool", volume: 6, yearBegan: 2015, name: "Deadpool (2015) #29", number: "29", coverDate: "2017-06-01", storeDate: "2017-04-19", image: "https://static.metron.cloud/media/issue/2022/05/21/deadpool-29.jpg", coverHash: "eabdc7126c07b116", modified: "2024-12-29T13:12:36.439218-05:00"},
	{kind: "character", id: 34067, seriesID: 2481, series: "Batman", volume: 1, yearBegan: 1940, name: "Batman (1940) #1", number: "1", coverDate: "1940-04-01", image: "https://static.metron.cloud/media/issue/2021/07/10/batman-1.jpg", coverHash: "ea2a919decc5934c", modified: "2024-12-20T14:50:51.216877-05:00"},
	{kind: "character", id: 34068, seriesID: 2481, series: "Batman", volume: 1, yearBegan: 1940, name: "Batman (1940) #2", number: "2", coverDate: "1940-07-01", image: "https://static.metron.cloud/media/issue/2021/07/10/batman-2.jpg", coverHash: "cf33b5c81d266c49", modified: "2024-12-20T14:50:53.328383-05:00"},
	{kind: "character", id: 29257, seriesID: 2102, series: "Detective Comics", volume: 1, yearBegan: 1937, name: "Detective Comics (1937) #45", number: "45", coverDate: "1940-11-01", image: "https://static.metron.cloud/media/issue/2021/04/15/detective-comics-45.jpg", coverHash: "f5acdb06707089e6", modified: "2024-12-23T14:04:38.901516-05:00"},
	{kind: "character", id: 34070, seriesID: 2481, series: "Batman", volume: 1, yearBegan: 1940, name: "Batman (1940) #4", number: "4", coverDate: "1941-01-01", image: "https://static.metron.cloud/media/issue/2021/07/10/batman-4.jpg", coverHash: "c3b73ccb2f0e3860", modified: "2024-12-20T14:50:58.372170-05:00"},
	{kind: "series", id: 7057, seriesID: 793, series: "Fables", volume: 1, yearBegan: 2002, name: "Fables (2002) #1", number: "1", coverDate: "2002-07-01", storeDate: "2002-05-08", image: "https://static.metron.cloud/media/issue/2019/11/01/fables-1.jpg", coverHash: "fc1952ece86782f0", modified: "2025-01-01T23:35:18.056873-05:00"},
	{kind: "series", id: 7058, seriesID: 793, series: "Fables", volume: 1, yearBegan: 2002, name: "Fables (2002) #2", number: "2", coverDate: "2002-08-01", storeDate: "2002-06-12", image: "https://static.metron.cloud/media/issue/2019/11/01/fables-2.jpg", coverHash: "9816e2a59bc62a3f", modified: "2025-01-01T23:35:21.122309-05:00"},
	{kind: "series", id: 7059, seriesID: 793, series: "Fables", volume: 1, yearBegan: 2002, name: "Fables (2002) #3", number: "3", coverDate: "2002-09-01", storeDate: "2002-07-17", image: "https://static.metron.cloud/media/issue/2019/11/01/fables-3.jpg", coverHash: "e59a3a6e697a3038", modified: "2025-01-01T23:35:24.039083-05:00"},
	{kind: "series", id: 7060, seriesID: 793, series: "Fables", volume: 1, yearBegan: 2002, name: "Fables (2002) #4", number: "4", coverDate: "2002-10-01", storeDate: "2002-08-14", image: "https://static.metron.cloud/media/issue/2019/11/01/fables-4.jpg", coverHash: "9524faf04ade2d38", modified: "2025-01-01T23:35:26.583714-05:00"},
	{kind: "team", id: 4123, seriesID: 475, series: "A-Force", volume: 1, yearBegan: 2015, name: "A-Force (2015) #1", number: "1", coverDate: "2015-07-01", storeDate: "2015-05-20", image: "https://static.metron.cloud/media/issue/2019/07/27/aforce-v1-1.jpg", coverHash: "ef8d84d08b1fd02b", modified: "2024-12-27T17:31:13.930112-05:00"},
	{kind: "team", id: 4124, seriesID: 475, series: "A-Force", volume: 1, yearBegan: 2015, name: "A-Force (2015) #2", number: "2", coverDate: "2015-09-01", storeDate: "2015-07-01", image: "https://static.metron.cloud/media/issue/2019/07/27/aforce-v1-2.jpg", coverHash: "8253c34687463f6f", modified: "2024-12-27T17:31:16.420503-05:00"},
	{kind: "team", id: 4125, seriesID: 475, series: "A-Force", volume: 1, yearBegan: 2015, name: "A-Force (2015) #3", number: "3", coverDate: "2015-10-01", storeDate: "2015-08-12", image: "https://static.metron.cloud/media/issue/2019/07/27/aforce-v1-3.jpg", coverHash: "d0955f2978dc4b4a", modified: "2024-12-27T17:31:19.430087-05:00"},
	{kind: "team", id: 4126, seriesID: 475, series: "A-Force", volume: 1, yearBegan: 2015, name: "A-Force (2015) #4", number: "4", coverDate: "2015-11-01", storeDate: "2015-09-09", image: "https://static.metron.cloud/media/issue/2019/07/27/aforce-v1-4.jpg", coverHash: "c2824bd7af3e2c0b", modified: "2024-12-27T17:31:21.662023-05:00"},
}

func issueListTestCases(t *testing.T, kind string) []testCase[*metron.IssueList] {
	var cases []testCase[*metron.IssueList]

	for _, r := range issueListRecords {
		if r.kind != kind {
			continue
		}

		expected := &metron.IssueList{
			ID: r.id,
			Series: metron.IssueListSeries{
				ID:        r.seriesID,
				Name:      r.series,
				Volume:    r.volume,
				YearBegan: r.yearBegan,
			},
			Name:      r.name,
			Number:    r.number,
			CoverDate: parseDate(t, r.coverDate),
			Modified:  parseTime(t, r.modified),
		}

		if r.storeDate != "" {
			expected.StoreDate = new(parseDate(t, r.storeDate))
		}

		if r.image != "" {
			expected.ImageURL = new(parseURL(t, r.image))
		}

		if r.coverHash != "" {
			expected.CoverHash = new(r.coverHash)
		}

		cases = append(cases, testCase[*metron.IssueList]{id: r.id, expected: expected})
	}

	return cases
}

// TestIssueOptionalFieldsAreGuarded covers the fields the spec allows the API
// to omit and the mapper rejects anyway, in both the absent and null forms.
func TestIssueOptionalFieldsAreGuarded(t *testing.T) {
	t.Parallel()

	for field, expected := range map[string]string{
		"alt_number":   "AltNumber",
		"name":         "Name",
		"rating_count": "RatingCount",
	} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()

			for _, form := range []string{"absent", "null"} {
				t.Run(form, func(t *testing.T) {
					t.Parallel()

					raw, err := fs.ReadFile("fixtures/issue_2558.json")
					require.NoError(t, err)

					var payload map[string]any

					require.NoError(t, json.Unmarshal(raw, &payload))

					if form == "absent" {
						delete(payload, field)
					} else {
						payload[field] = nil
					}

					body, err := json.Marshal(payload)
					require.NoError(t, err)

					c := newTestClient(t, []requestMock{
						{expectedURL: "https://metron.cloud/api/issue/2558/", responseBody: string(body)},
					})

					_, err = c.IssueByID(context.Background(), 2558)

					var mapErr *metron.MapError

					require.ErrorAs(t, err, &mapErr)
					require.Equal(t, "issue", mapErr.Kind)
					require.Equal(t, 2558, mapErr.ID)
					require.Equal(t, expected, mapErr.Field)

					require.EqualError(t, err, "metron: issue 2558: nil "+expected)
				})
			}
		})
	}
}

// TestIssueListUnparseableImage pins the list mapper's image failure to
// MapError, so a listing can keep skipping past it like any other bad record.
func TestIssueListUnparseableImage(t *testing.T) {
	t.Parallel()

	body := `{"count":2,"next":null,"previous":null,"results":[
		{"id":1,"series":{"id":9,"name":"S","volume":1,"year_began":2000},"number":"1","cover_date":"2022-09-01","image":"https://x/%zz","modified":"2024-12-23T16:17:29.333263-05:00"},
		{"id":2,"series":{"id":9,"name":"S","volume":1,"year_began":2000},"number":"2","cover_date":"2022-09-01","modified":"2024-12-23T16:17:29.333263-05:00"}
	]}`

	c := newTestClient(t, []requestMock{
		{expectedURL: "https://metron.cloud/api/issue/?page=1", responseBody: body},
	})

	var ids []int
	var errs []error

	for issue, err := range c.Issues(context.Background()) {
		if err != nil {
			errs = append(errs, err)

			continue
		}

		ids = append(ids, issue.ID)
	}

	require.Equal(t, []int{2}, ids)
	require.Len(t, errs, 1)
	require.ErrorIs(t, errs[0], &metron.MapError{Kind: "issue", ID: 1, Field: "Image"})

	var mapErr *metron.MapError

	require.ErrorAs(t, errs[0], &mapErr)
	require.ErrorContains(t, mapErr.Err, `invalid URL escape "%zz"`)
	require.ErrorContains(t, errs[0], "metron: issue 1: Image: ")
}

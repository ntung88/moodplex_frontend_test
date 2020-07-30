package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/jzelinskie/geddit"
	hn "github.com/ytkhs/hackernews-api-go"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type DataSource interface {
	// Acquires posts from a given a data source and adds them to the database
	fillDB()
}
type TwitterSrc struct{ source string }
type RedditSrc struct{ source string }
type YouTubeSrc struct{ source string }
type ImgurSrc struct{ source string }
type HackerNewsSrc struct{ source string }

var allMoods = []string{"happy", "funny", "informative", "motivational", "sad",
	"cute", "educational", "angry", "uplifting", "scary", "artistic", "news",
	"romantic"}
// var twitterSrc = TwitterSrc{"twitter"}
var redditSrc = RedditSrc{"reddit"}
var youtubeSrc = YouTubeSrc{"youtube"}
// var imgurSrc = ImgurSrc{"imgur"}
var hackernewsSrc = HackerNewsSrc{"hackernews"}
var dataSources = []DataSource{
	// twitterSrc,
	redditSrc,
	youtubeSrc,
	// imgurSrc,
	hackernewsSrc,
}
var utcLocation, _ = time.LoadLocation("UTC")

const SeparateMedia = false
const InitRating = 1000

// Adds multiple entries to the database for a single social media post with
// varying media URLs
func addMediaPosts(mediaUrls []string, postUrl, src string, nsfw bool,
	misc, addDate, publishDate string) {
	for _, mediaUrl := range mediaUrls {
		addPost(InitRating, mediaUrl, postUrl, src, allMoods, nsfw, misc,
			addDate, publishDate)
	}
}

// Adds a slice of Tweets to the database
func (src TwitterSrc) addTweetSlice(tweetArr []anaconda.Tweet, moods []string) {
	for _, tweet := range tweetArr {
		postUrl := "https://twitter.com/" + tweet.User.ScreenName +
			"/status/" + tweet.IdStr
		var mediaUrl = ""

		addDate := strings.Split(time.Now().In(utcLocation).String(),
			".")[0] + " UTC"

		// Twitter API format: "Wed Oct 10 20:19:24 +0000 2018"
		splitTimeStamp := strings.SplitAfter(tweet.CreatedAt, " ")
		var publishDate string
		for i := 1; i < 4; i++ {
			publishDate += splitTimeStamp[i]
		}
		publishDate += splitTimeStamp[5] + " UTC"
		if mediaArr := tweet.Entities.Media; len(mediaArr) > 0 &&
			SeparateMedia {
			var mediaUrls []string
			for _, entity := range mediaArr {
				mediaUrls = append(mediaUrls, entity.Media_url)
			}
			mediaUrl = ""
			addMediaPosts(mediaUrls, postUrl, src.source,
				tweet.PossiblySensitive, "", addDate, publishDate)
		}
		addPost(InitRating, mediaUrl, postUrl, src.source, moods,
			tweet.PossiblySensitive, "", addDate,
			publishDate)
	}
}

// Acquires Tweets from various Twitter lists and adds them to the database
func (src TwitterSrc) addTweetsFromList(api *anaconda.TwitterApi,
	tweetsPerList int64, listIDs []int64, moods []string) {
	// # of Tweets per call defaults to 20 when not specified
	searchParams := url.Values{}
	searchParams.Set("count", strconv.FormatInt(tweetsPerList, 10))
	for _, twitterList := range listIDs {
		// API Limit: 60 requests/min
		tweetArr, err := api.GetListTweets(twitterList, false,
			searchParams)
		if err != nil {
			panic("Cannot get Tweets from List #" + strconv.FormatInt(
				twitterList, 10))
		}
		src.addTweetSlice(tweetArr, moods)
	}
}

// Acquires Tweets obtained from various Twitter searches and adds them to the
// database
func (src TwitterSrc) addTweetsFromSearches(api *anaconda.TwitterApi,
	tweetsPerSearch int64, searchPhrases []string, moods []string) {
	// # of Tweets per call defaults to 15 when not specified
	// Max: 100
	searchParams := url.Values{}
	searchParams.Set("count", strconv.FormatInt(tweetsPerSearch, 10))
	for _, searchPhrase := range searchPhrases {
		// API Limit: 30 requests/min
		tweetSearch, err := api.GetSearch(searchPhrase, searchParams)
		if err != nil {
			panic("Cannot get Tweets using seach phrase: " + searchPhrase)
		}
		src.addTweetSlice(tweetSearch.Statuses, moods)
	}
}

func (src TwitterSrc) fillDB() {
	api := anaconda.NewTwitterApiWithCredentials(
		"2309313565-tz09wVjG4WNWLGyaOF78GGZtkiytJuACvbsiZ0O",
		"oV3dyAV8G1SUGpPXE0J09mPzcbrkOJy8uv4odPLN7pk9Z",
		"5k9CUEdLVfeZnKSx7u9SBuvrW",
		"kP2l9ZT1pwyAeQjqLm4HhcYX3xkUTAx5cZss2YmfkLnX7mLu6T")

	// HAPPY, SAD, & ANGRY
	// Note: Can't find happy-/angry-specific content,
	// so these are some general Twitter lists.
	happySadAngryLists := []int64{1860883, 4343, 61394, 100791150, 21189472}
	src.addTweetsFromList(api, int64(30), happySadAngryLists,
		[]string{"happy", "sad", "angry"})

	// FUNNY
	src.addTweetsFromList(api, int64(150), []int64{1278068055671783427},
		[]string{"funny"})

	// INFORMATIVE & EDUCATIONAL
	src.addTweetsFromList(api, int64(150), []int64{1278078570410192899},
		[]string{"informative", "educational"})

	// MOTIVATIONAL & UPLIFTING
	src.addTweetsFromList(api, int64(150), []int64{1278347525502783488},
		[]string{"motivational", "uplifting"})

	// CUTE
	src.addTweetsFromList(api, int64(30), []int64{1278333876918779905},
		[]string{"cute"})
	src.addTweetsFromSearches(api, int64(40), []string{"%23wholesome",
		"%23cute", "%23adorable"}, []string{"cute"})

	// SCARY
	src.addTweetsFromList(api, int64(22), []int64{1278389767760236544},
		[]string{"scary"})
	src.addTweetsFromSearches(api, int64(32), []string{"%23scarystories",
		"%23scary", "%23horror", "%23creepy"}, []string{"scary"})

	// ARTISTIC
	src.addTweetsFromList(api, int64(150), []int64{1278393560883253249},
		[]string{"artistic"})

	// NEWS
	src.addTweetsFromList(api, int64(150), []int64{1278059091584258049},
		[]string{"news"})

	// ROMANTIC
	src.addTweetsFromSearches(api, int64(30), []string{"%23love", "%23romantic",
		"%23romance", "%23relationship", "%23relationshipgoals"},
		[]string{"romantic"})
}

func (src RedditSrc) fillDB() {
	// API Limit: 30 requests/min

	// Creates a new session by logging in
	session, err := geddit.NewOAuthSession(
		"nAd-b5xh_dzqRQ",
		"s3jgYRaDU9YhtHR4HmWXmjfHOhA",
		"Moodplex by u/moodplex_api v0.1 see source https://github."+
			"com/sohamsankaran/moodplex",
		"http://localhost:8080",
	)
	if err != nil {
		panic("Cannot authenticate Reddit session")
	}

	// Creates new auth token for script apps (
	// apps that don't require specific user data).
	err = session.LoginAuth("moodplex_api", "soham.sankaran")
	if err != nil {
		panic("Cannot acquire new token")
	}

	/*
		Hot: Posts that are getting lots of upvotes & comments recently
		Top: Posts with the most upvotes given a set period (
		currently set to 24 hrs)
		Rising: Posts getting a lot of activity RIGHT NOW
	*/
	const postCt = 100 // max: 100
	hotSubmissions, err := session.SubredditSubmissions("popular",
		geddit.HotSubmissions, geddit.ListingOptions{Limit: postCt})
	if err != nil {
		panic("Cannot acquire hot submissions")
	}
	src.addSubredditPosts(hotSubmissions)

	topSubmissions, err := session.SubredditSubmissions("popular",
		geddit.TopSubmissions, geddit.ListingOptions{Time: geddit.ThisDay,
			Limit: postCt})
	if err != nil {
		panic("Cannot acquire top submissions with set period: " +
			geddit.ThisDay)
	}
	src.addSubredditPosts(topSubmissions)

	risingSubmissions, err := session.SubredditSubmissions("popular",
		geddit.TopSubmissions, geddit.ListingOptions{Limit: postCt})
	if err != nil {
		panic("Cannot acquire rising submissions")
	}
	src.addSubredditPosts(risingSubmissions)
}

func (r RedditSrc) addSubredditPosts(posts []*geddit.Submission) {
	const redditLinkPrefix = "https://www.reddit.com"
	for _, post := range posts {
		if post != nil {
			postUrl := redditLinkPrefix + post.Permalink
			addDate := strings.Split(time.Now().In(utcLocation).String(),
				".")[0] + " UTC"
			// Reddit API gives the UNIX time in the current location
			publishDateChars := []rune(time.Unix(int64(post.DateCreated),
				0).In(utcLocation).String())
			publishDate := string(publishDateChars[:19]) + " UTC"
			// post.URL is the URL of the media attached,
			// so if they are equal, then there is no media attached
			if post.URL != postUrl && SeparateMedia {
				addPost(InitRating, post.URL, postUrl, r.source, allMoods,
					post.IsNSFW, "", addDate,
					publishDate)
			}
			addPost(InitRating, "", postUrl, r.source, allMoods,
				post.IsNSFW, "", addDate,
				publishDate)
		}
	}
}

func (src YouTubeSrc) fillDB() {
	// API Limit: 10,000 units/day. Each videos list call costs 1 + (6 units per
	// video for its content details, snippet,
	// and status since each of the three fields costs 2 units).
	// 50 videos are acquired per call,
	// so each videos list call is 301 units. So, this method uses 1204 units.
	const youTubePrefix = "https://www.youtube.com/watch?v="
	const videoCt = 50 // max: 50

	ctx := context.Background()
	// Creates new Service and authenticates with API key
	youtubeService, err := youtube.NewService(ctx,
		option.WithAPIKey("AIzaSyAtYSWTapEx2fd_H7htyvmrI9xpYOV8tJA"))
	if err != nil {
		panic("Cannot connect to YouTube API")
	}

	videosService := youtube.NewVideosService(youtubeService)
	videosListCall := videosService.List([]string{"id",
		"contentDetails", "snippet",
		"status"})
	videosListCall.Chart("mostPopular")
	videosListCall.MaxResults(videoCt)
	videosListCall.RegionCode("us")
	nextPageToken := ""
	// There are 200 results, and each call allows for a max. of 50 videos
	for i := 0; i < 4; i++ {
		videosListCall.PageToken(nextPageToken)
		resp, err := videosListCall.Do()
		if err != nil {
			panic("Cannot perform search for most popular YouTube videos")
		}
		for _, video := range resp.Items {
			var nsfw bool
			// Labeling videos with a YouTube rating restriction as NSFW
			switch video.ContentDetails.ContentRating.YtRating {
			case "ytUnspecified":
				nsfw = false
			case "ytAgeRestricted":
				nsfw = true
			default:
				nsfw = false
			}
			source := youTubePrefix + video.Id
			// For YouTube, the actual post is the media, so there's no need to
			// add a media URL
			addDate := strings.Split(time.Now().In(utcLocation).String(),
				".")[0] + " UTC"
			addPost(InitRating, "", source, src.source, allMoods, nsfw, "",
				addDate, video.Snippet.PublishedAt)
			nextPageToken = resp.NextPageToken
		}
	}
}

type ImgurAlbumImage struct {
	Image string `json:"link"'`
}

type ImgurObject struct {
	Id       string `json:"id"`
	DateTime int    `json:"datetime"`
	Nsfw     bool   `json:"nsfw"`
	/*
		An ImgurObject can either be a post or an album. A post must only
		contain 1 image, but an album doesn't have to be multiple images.
		It's unclear what the difference is between a public Imgur post and
		1-image Imgur album
	*/
	Images []ImgurAlbumImage `json:"images"`
}

type ImgurResponse struct {
	Data []ImgurObject `json:"data"`
}

func (src ImgurSrc) fillDB() {
	// API Limit: ~12,500 requests/day
	/*
		Section:
		Hot - most viral (sorted by popularity)
		NOTE: 1st page (page=0) has 60 results. 2nd page has ~500
		Top - most viral (sorted by highest scoring)
		NOTE: Each page has 60 results
	*/
	const apiUrlPrefix = "https://api.imgur.com/3/gallery/"
	const urlPrefix = "https://imgur.com/gallery/"
	sectionParams := [2]string{"hot", "top"}
	method := "GET"

	client := &http.Client{}
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.Close()
	if err != nil {
		panic("Cannot close writer")
	}
	for i := 0; i < len(sectionParams); i++ {
		for pageNum := 0; pageNum < 2; pageNum++ {
			apiURL := apiUrlPrefix + sectionParams[i] + "/viral/day/" +
				strconv.FormatInt(int64(pageNum), 10) +
				"/?showViral=true&mature=true&album_previews=true"
			req, err := http.NewRequest(method, apiURL, payload)
			if err != nil {
				panic("Cannot perform Imgur API GET request")
			}
			req.Header.Add("Authorization",
				"Client-ID c70cd8d165c75b7")
			req.Header.Set("Content-Type", writer.FormDataContentType())
			res, err := client.Do(req)
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			var posts ImgurResponse
			err = json.Unmarshal(body, &posts)
			if err != nil {
				panic(err)
			}
			for _, post := range posts.Data {
				postURL := urlPrefix + post.Id
				addDate := strings.Split(time.Now().In(utcLocation).String(),
					".")[0] + " UTC"
				// Imgur API gives UNIX time in local time zone
				publishDateChars := []rune(time.Unix(int64(post.DateTime),
					0).In(utcLocation).String())
				publishDate := string(publishDateChars[:19]) + " UTC"
				if imageArr := post.Images; len(imageArr) > 1 && SeparateMedia {
					var imageURLs []string
					for _, albumImage := range imageArr {
						imageURLs = append(imageURLs, albumImage.Image)
					}
					addMediaPosts(imageURLs, postURL, src.source, post.Nsfw,
						"", addDate, publishDate)
				} else {
					// Since each post/album on Imgur contains an image, there
					// is no need to create a separate entry if there's only 1
					// image, nor is there a need to add a URL for the media.
					addPost(InitRating, "", postURL, src.source, allMoods,
						post.Nsfw, "", addDate, publishDate)
				}
			}
		}
	}
}

func (src HackerNewsSrc) fillDB() {
	// API Limit: No limit
	const hackerNewsPrefix = "https://news.ycombinator.com/item?id="
	deadList := []string{"https://www.bloomberg.com",
		"https://www.nasa.gov/", "https://coil.com", "https://www.vice.com"}
	// Gets newest 500 stories
	itemIDArr, err := hn.GetStories("new")
	fmt.Println()
	if err != nil {
		panic("Cannot get new Hacker News stories")
	}
	for _, itemID := range itemIDArr {
		itemIDString := strconv.FormatInt(int64(itemID), 10)
		item, err := hn.GetItem(itemID)
		if err != nil {
			panic("Cannot get story #" + itemIDString)
		}
		// Dead posts are flagged by moderators
		// If the title is empty, the JSON is null
		if !item.Dead && item.Title != "" {
			addDate := strings.Split(time.Now().In(utcLocation).String(),
				".")[0] + " UTC"
			// HN API gives Unix time in local time zone
			publishDateChars := []rune(time.Unix(item.Time,
				0).In(utcLocation).String())
			publishDate := string(publishDateChars[:19]) + " UTC"

			hnURL := hackerNewsPrefix+itemIDString
			var allowIFrame, iframeURL string
			if item.URL == "" {
				iframeURL = hnURL
			} else {
				iframeURL = item.URL
			}
			isDeadLink := false
			for _, link := range deadList {
				if strings.Contains(iframeURL, link) {
					isDeadLink = true
					break
				}
			}
			// Check if it can be displayed in an iframe
			resp, err := http.Get(iframeURL)
			if err != nil {
				log.Println("Cannot perform GET request to " + iframeURL)
				allowIFrame = "false"
			} else {
				xFrameOptions := resp.Header.Get("X-Frame-Options")
				if xFrameOptions == "" && !isDeadLink  {
					allowIFrame = "true"
				} else {
					// Possible values are SAMEORIGIN, DENY, ALLOW-FROM *insert uri*
					allowIFrame = "false"
				}
			}
			miscJSON := "{\"title\": \"" + item.Title + "\", " +
				"\"allowIFrame\": " + allowIFrame + "}"

			// Hacker News API does not provide a field for NSFW/sensitive
			// content, so it's set to false
			// Since there's always <= 1 media link per post AND we are
			// embedding the attached media onto Moodplex,
			// it's always added on regardless of the value of separateMedia
			addPost(InitRating, item.URL, hnURL,
				src.source, allMoods, false, miscJSON,
				addDate, publishDate)
		}
	}
}

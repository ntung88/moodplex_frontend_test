let NUM_RESULTS = 10;
EXPIRY_TIME = 288000000

function checkExpiry() {
    json = JSON.parse(localStorage.getItem("ranking_count"))
    const now = new Date()
    console.log("current time: " + String(now.getTime()))
    console.log("expiry: " + String(json["expiry"]))
    if (now.getTime() > json["expiry"]) {
        json = {}
        json["expiry"] = now.getTime() + EXPIRY_TIME
        json[window.site] = 0
    }
    localStorage.setItem("ranking_count", JSON.stringify(json))
}

var ranking_count = localStorage.getItem('ranking_count')
if (ranking_count == null) {
    json = {}
    json[window.site] = 0
    const now = new Date()
    json["expiry"] = now.getTime() + EXPIRY_TIME
    localStorage.setItem('ranking_count', JSON.stringify(json));
    console.log("ranking_count reinitialized")
} else if (JSON.parse(ranking_count)[window.site] == null) {
    json = JSON.parse(ranking_count)
    json[window.site] = 0
    localStorage.setItem('ranking_count', JSON.stringify(json));
    console.log("ranking_count reinitialized")
}
checkExpiry()


let posts;
let post1 = document.getElementById('post_div1');
let post2 = document.getElementById('post_div2');

// function expandListener()
// {
//     console.log("clicked");
//     var element = document.getElementById("post1");
//     element.classList.add("expanded");
// }

function appendResultsButtons() {
    var storage = JSON.parse(localStorage.getItem("ranking_count"))
    let count = storage[window.site]
    if (count >= 5 && count % 5 === 0) {
        let top = count / 5 * NUM_RESULTS
        let navbar = document.getElementById("navbar")
        let siteName;
        switch (window.site) {
            case "twitter":
                siteName = "Twitter";
                break;
            case "reddit":
                siteName = "Reddit";
                break;
            case "imgur":
                siteName = "Imgur";
                break;
            case "youtube":
                siteName = "YouTube";
                break;
            case "hackernews":
                siteName = "Hacker News";
                break;
        }
        navbar.innerHTML = navbar.innerHTML + "<a class=\"nohover\">Rankings:</a>"
        for (let i = NUM_RESULTS; i <= top; i += NUM_RESULTS) {
            navbar.innerHTML = navbar.innerHTML + `<a href="/rankings/` + window.mood + `/` + window.site + `/` + i + `">`
                + siteName + `: ` + (i - NUM_RESULTS + 1) + `-` + i + `</a>`
        }
    }
}

function addHNPosts(response) {
    const IFRAME_WARNING = "This website does not allow for it to be" +
        " on Moodplex. Click on the links above to view the post.";

    let miscObj1 = JSON.parse(response['misc1']);
    let cat_1_url = response['source1'];
    if (cat_1_url === "") {
        cat_1_url = response['agridata_source1'];
    }
    let titleLink1 = '<h2><a href="' + cat_1_url + '" ' +
        'target="_blank">' + miscObj1.title + '</a></h2>';
    let hnLink1 = `<h3><a href="` + response['agridata_source1'] +
        `" target="_blank" id="hnLink1">Hacker News Comments</a></h3>`
    post1.innerHTML = titleLink1 + hnLink1;

    let miscObj2 = JSON.parse(response['misc2']);
    let cat_2_url  = response['source2'];
    if (cat_2_url === "") {
        cat_2_url = response['agridata_source2'];
    }
    let titleLink2 = '<h2><a href="' + cat_2_url + '" ' +
        'target="_blank">' + miscObj2.title + '</a></h2>';
    let hnLink2 = `<h3><a href="` + response['agridata_source2']
        + `" target="_blank" id="hnLink2">Hacker News Comments</a></h3>`
    post2.innerHTML = titleLink2 + hnLink2;

    let iframe1 = document.getElementById("post1");
    let iframeAlert1 = document.getElementById("iframe-alert1");
    if (miscObj1.allowIFrame === true) {
        if (iframe1 === null) {
            let iframe1 = `<iframe src="` + cat_1_url + `" class="posts"
                            id="post1"></iframe>`;
            post1.innerHTML += iframe1;
        } else {
            iframe1.setAttribute("src", cat_1_url);
        }
        if (iframeAlert1 != null) {
            iframeAlert1.remove();
        }
    } else {
        if (iframeAlert1 === null) {
            post1.innerHTML += "<p id='iframe-alert1'>" + IFRAME_WARNING + "</p>";
        }
        if (iframe1 != null) {
            iframe1.remove();
        }
    }

    let iframe2 = document.getElementById("post2");
    let iframeAlert2 = document.getElementById("iframe-alert2");
    if (miscObj2.allowIFrame === true) {
        if (iframe2 === null) {
            let iframe2 = `<iframe src="` + cat_2_url + `" class="posts"
                            id="post2"></iframe>`;
            post2.innerHTML += iframe2;
        } else {
            iframe2.setAttribute("src", cat_2_url);
        }
        if (iframeAlert2 != null) {
            iframeAlert2.remove();
        }
    } else {
        if (iframeAlert2 === null) {
            post2.innerHTML += "<p id='iframe-alert2'>" + IFRAME_WARNING + "</p>";
        }
        if (iframe2 != null) {
            iframe2.remove();
        }
    }

    console.log("metadata: " + cat_1_url + " " + cat_2_url);
}

function addTwitterPosts(response) {
    let container1 = `<div id="container1"></div>`;
    let container2 = `<div id="container2"></div>`;
    post1.innerHTML = container1;
    post2.innerHTML = container2;
    let element1 = document.getElementById("container1");
    let element2 = document.getElementById("container2");

    let url1 = response['agridata_source1'];
    let id1 = url1.substring(url1.lastIndexOf('/') + 1);

    let url2 = response['agridata_source2'];
    let id2 = url2.substring(url2.lastIndexOf('/') + 1);

    twttr.widgets.createTweet(
        id1, post1,
        {
            conversation: 'all',    // or none
            cards: 'visible',  // or hidden
            linkColor: '#cc0000', // default is blue
            theme: 'light',    // or dark
            align: 'center' // or left or right
        })

    twttr.widgets.createTweet(
        id2, post2,
        {
            conversation : 'all',    // or none
            cards        : 'visible',  // or hidden
            linkColor    : '#cc0000', // default is blue
            theme        : 'light',    // or dark
            align: 'center' // or left or right
        })
}

function addRedditPosts(response) {
    let url1 = response['agridata_source1'].split("/comments/")[0];
    let url2 = response['agridata_source2'].split("/comments/")[0];
    let block_1 = "<blockquote class=\"reddit-card\"><a href=\"" + response["agridata_source1"] + "\"></a><a href=\"" + url1 + "\"></a></blockquote>";
    let block_2 = "<blockquote class=\"reddit-card\"><a href=\"" + response["agridata_source2"] + "\"></a><a href=\"" + url2 + "\"></a></blockquote>";
// <blockquote class="reddit-card" data-card-created="1595219912"><a href="https://www.reddit.com/r/Fencing/comments/htuvtw/whats_your_thought_when_you_first_started_to_fence/">What’s your thought when you first started to fence ?</a> from <a href="http://www.reddit.com/r/Fencing">r/Fencing</a></blockquote>
//     <script async src="//embed.redditmedia.com/widgets/platform.js" charset="UTF-8"></script>
    post1.innerHTML = block_1;
    post2.innerHTML = block_2;
}

function addYouTubePosts(response) {
    const aspectRatio = 9/16;
    let videoWidth = 0.95 * window.innerWidth / 2;
    let videoHeight = aspectRatio * videoWidth;
    let splitURL1 = response['agridata_source1'].split("watch?v=");
    let embedURL1 = splitURL1[0] + "embed/" + splitURL1[1];
    post1.innerHTML = "<iframe width='" + videoWidth.toString(10) + "' height='" + videoHeight.toString(10) + "' src=\"" + embedURL1 + "\"" +
        " frameborder=\"0\"" + " allow=\"accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture\" allowfullscreen></iframe>"
    let splitURL2 = response['agridata_source2'].split("watch?v=");
    let embedURL2 = splitURL2[0] + "embed/" + splitURL2[1];
    post2.innerHTML = "<iframe width='" + videoWidth.toString(10) + "' height='" + videoHeight.toString(10) + "' src=\"" + embedURL2 + "\"" +
        " frameborder=\"0\"" + " allow=\"accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture\" allowfullscreen></iframe>"
}

function addPostsToPage(response) {
    console.log(localStorage);
    console.log(response);
    posts = response;

    appendResultsButtons()
    switch (window.site) {
        case "hackernews":
            addHNPosts(response);
            break;
        case "twitter":
           addTwitterPosts(response);
            break;
        case "reddit":
            addRedditPosts(response);
            break;
        case "youtube":
            addYouTubePosts(response);
            break;
    }
}

function sendcatJSON() {
    var data = JSON.stringify({"mood": window.mood, "site": window.site });

    $.ajax
    ({
        type: "POST",
        //the url where you want to sent the userName and password to
        url: '/posts',
        dataType: 'json',
        contentType: 'application/json',
        async: false,
        //json object to sent to the authentication url
        data: data,
        success: addPostsToPage
    })
}

function sendMatchJSON(score1, score2) {
    let row_id1 = posts['row_id1'];
    let row_id2 = posts['row_id2'];

    console.log("id1: "+ row_id1);
    console.log("id2: "+ row_id2);
    console.log("score1: " + score1);
    console.log("score2: "+ score2);

    // Converting JSON data to string
    let data = JSON.stringify({ "rowid1": row_id1, "rowid2": row_id2, "score1": score1, "score2": score2});

    $.ajax
    ({
        type: "POST",
        //the url where you want to sent the userName and password to
        url: '/match',
        dataType: 'json',
        contentType: 'application/json',
        async: false,
        //json object to sent to the authentication url
        data: data,
        success: function() {
            console.log('Match Recorded')
            location.reload()
        }})

    let count_json = JSON.parse(localStorage.getItem('ranking_count'))
    count = count_json[window.site]
    let new_count = count + 1
    count_json[window.site] = new_count
    updateBar(new_count)
    localStorage.setItem('ranking_count', JSON.stringify(count_json))
    if ((new_count)%5 === 0) {
        // new_count = count + 1
        // localStorage.setItem('ranking_count', JSON.stringify({'count': 0}))
        location.href = "/rankings/" + window.mood + "/" + window.site + "/none"
    } else {
        sendcatJSON();
    }
}

function updateBar(new_count) {
    let relative_val = (new_count % 5) + 1
    document.getElementById('file').value = relative_val * 20;
    document.getElementById('postNumber').innerHTML = String(relative_val) + "/5"
    var disclaimer = document.getElementById("num_left")
    let top = (Math.trunc(new_count / 5) + 1) * NUM_RESULTS
    disclaimer.innerHTML = "Compare " + (5 - (new_count%5)) + " more pairs to see results " + (top - NUM_RESULTS + 1) + "-" + top
}

function updateTitle() {
    switch (window.site) {
        case "twitter":
            $("title").text("Moodplex | Twiiter");
            break;
        case "reddit":
            $("title").text("Moodplex | Reddit");
            break;
        case "youtube":
            $("title").text("Moodplex | YouTube");
            break;
        case "hackernews":
            $("title").text("Moodplex | Hacker News");
            break;
    }
}

sendcatJSON();
updateBar(JSON.parse(localStorage.getItem('ranking_count'))[window.site])
updateTitle();

let NUM_RESULTS = 10;
let data = JSON.stringify({"mood": window.mood, "site": window.site, "num":
    NUM_RESULTS});

function appendResultsButtons() {
    var count = JSON.parse(localStorage.getItem("ranking_count"))[window.site]
    if (count >= 5 && count % 5 === 0) {
        let top = count / 5 * NUM_RESULTS
        var navbar = document.getElementById("navbar")
        let siteName
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
            navbar.innerHTML = navbar.innerHTML + `<a href="/results/` + window.mood + `/` + window.site + `/` + i + `">` 
                + siteName + `: ` + (i - NUM_RESULTS + 1) + `-` + i + `</a>`
        }
        navbar.innerHTML = navbar.innerHTML + `<a href="/category/` + window.mood + `/` + window.site + `">Continue Comparing to See More of the Rankings</a>`
    }
}

function addResultsToPage(response) {
    appendResultsButtons()
    console.log(response)
    var container = document.getElementById('results')
    let count = JSON.parse(localStorage.getItem('ranking_count'))[window.site]
    let top = (Math.trunc(count / 5) - 1) * NUM_RESULTS
    let start = top
    let end = (top + NUM_RESULTS)
    if (window.upto !== "none") {
        start = window.upto - NUM_RESULTS
        end = parseInt(window.upto)
    }
    for (let i = start; i < end; i++) {
        let element;
        switch (window.site) {
            case "twitter":
                element = `<button class="pairing-accordion" href = "">
                        #` + (i + 1) + `
                        </button>
                        <div id="result` + (i + 1) + `" class = "b-item" alt = "placeholder" width = "350" height = "350"></div>
                        <br>`
                break;
            case "hackernews":
                let url = response['sources'][i]
                if (url === "") {
                    url = response['agridata_sources'][i]
                }
                let comments_url = response['agridata_sources'][i]
                element = `#` + (i + 1) + `: <a style="font-size:20px;" class="pairing-accordion" href="` + url + `">`
                         + JSON.parse(response['miscs'][i]).title + `
                        </a>
                        <br>
                        <br>
                        <a class="comments" href=\"` + comments_url + `\"" class = "b-item" alt = "placeholder" width = "350" height = "350">Hacker News Comments</a>
                        <br>
                        <br>
                        <br>`
                break;
            case "reddit":
                let cat_url = response['agridata_sources'][i].split("/comments/")[0]
                element = `<button class="pairing-accordion" href = "">
                            #` + (i + 1) + `
                            </button>
                            <blockquote class="reddit-card">
                            <a href="` + response['agridata_sources'][i] + `"></a>
                            <a href="` + cat_url + `"></a>
                            </blockquote>
                            <br>`
                break;
            case "youtube":
                const aspectRatio = 9/16;
                let videoWidth = 0.95 * window.innerWidth / 2;
                let videoHeight = aspectRatio * videoWidth;
                let splitURL = response['agridata_sources'][i].split("watch?v=");
                let embedURL = splitURL[0] + "embed/" + splitURL[1];
                let iframe = "<iframe width='" + videoWidth.toString(10) +
                    "' height='" + videoHeight.toString(10) + "' src=\"" + embedURL + "\"" +
                    " frameborder=\"0\"" + " allow=\"accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture\" allowfullscreen></iframe>"
                element = `<button class="pairing-accordion" href = "">
                            #` + (i + 1) + `
                            </button><br><br>` + iframe + "\n<br><br>"
                break;
        }
        container.innerHTML = container.innerHTML + element;
    }
    if (window.site === "twitter") {
        for (var i = top; i < (top + NUM_RESULTS); i++) {
            var container = document.getElementById('result' + (i + 1));
            var url = response['agridata_sources'][i]
            var id = url.substring(url.lastIndexOf('/') + 1);

            twttr.widgets.createTweet(
                id, container,
                {
                    conversation: 'all',    // or none
                    cards: 'visible',  // or hidden
                    linkColor: '#cc0000', // default is blue
                    theme: 'light',    // or dark
                    align: 'center' // or left or right
                })
        }
    }
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

updateTitle();
$.ajax
({
    type: "POST",
    //the url where you want to sent the userName and password to
    url: '/getresults',
    dataType: 'json',
    contentType: 'application/json',
    async: false,
    //json object to sent to the authentication url
    data: data,
    success: addResultsToPage
});

document.getElementById('rank_button').onclick = function(){
    location.href = "/category/" + window.mood + "/" + window.site;
};
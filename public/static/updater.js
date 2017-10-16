window.checkForUpdates = function(version) {
    var cachedXHR = function (url, cacheSeconds, callback, errorCallback) {
        var cacheKey = "xhrcache|" + url + "|";
        var cacheTimeKey = cacheKey + "time";
        var cacheValueKey = cacheKey + "value";

        var currentTime = Math.round(new Date().getTime() / 1000);
        var cacheValue = localStorage.getItem(cacheValueKey);
        var cacheTime = 0;

        try {
            cacheTime = parseInt(localStorage.getItem(cacheTimeKey));
        } catch (e) {
            localStorage.setItem(cacheTimeKey, 0);
            cacheTime = 0;
        }

        if (cacheValue === null || (currentTime - cacheTime) > cacheSeconds) {
            var xhttp = new XMLHttpRequest();
            xhttp.onreadystatechange = function () {
                if (this.readyState == 4 && this.status == 200) {
                    var resp = this.responseText;
                    localStorage.setItem(cacheTimeKey, currentTime);
                    localStorage.setItem(cacheValueKey, resp);
                    callback(resp);
                } else if (this.readyState == 4 && this.status > 399) {
                    errorCallback("Error: HTTP status " + this.status.toString());
                }
            };
            xhttp.onerror = function () {
                errorCallback("Error: Network error");
            };
            xhttp.open("GET", url, true);
            xhttp.send();
        } else {
            callback(cacheValue);
        }
    };
    var HOUR = 3600;
    cachedXHR("https://api.github.com/repos/geek1011/BookBrowser/releases", HOUR / 2, function (respa) {
        cachedXHR("https://api.github.com/repos/geek1011/BookBrowser/releases/latest", HOUR / 2, function (respb) {
            try {
                var releases = JSON.parse(respa);
                var current = JSON.parse(respb);

                var currentVersion = version;

                var isDev = (currentVersion.indexOf("dev") > -1) || (currentVersion.indexOf("+") > -1);
                if (isDev) {
                    console.warn("You are using a development version of BookBrowser");
                    return;
                }

                var latestVersion = current["tag_name"];
                if (latestVersion == currentVersion) {
                    console.info("You are using the latest version of BookBrowser: " + latestVersion);
                    return;
                }

                console.info("You are not using the latest version of BookBrowser. Your current version is " + currentVersion + ", but the latest version is " + latestVersion);
                var releaseNotes = "";
                for (var i = 0; i < releases.length; i++) {
                    var release = releases[i];
                    if (release["tag_name"] == currentVersion) {
                        break;
                    }
                    releaseNotes += "<div class=\"release\"><div class=\"version\">" + release["tag_name"] + "</div><div class=\"changelog\">" + release.body.split("## Usage")[0].split("\n").filter(function (l) {
                        return !(l.indexOf("Changes for") > -1) && (l !== "");
                    }).map(function (l) {
                        return l + "<br>";
                    }).join("\n") + "</div></div>";
                }
                releaseNotes += [
                    "<style>",
                    ".release .changelog {",
                    "    display: block;",
                    "",
                    "}",
                    ".release .version {",
                    "    display: block;",
                    "    font-weight: bold;",
                    "}",
                    ".release {",
                    "    display: block;",
                    "    margin: 20px 0;",
                    "    padding: 10px;",
                    "    border: 1px solid #CCCCCC;",
                    "    background: #F0F0F0;",
                    "    border-radius: 5px;",
                    "}",
                    "</style>"
                ].join("\n");

                var message = "<b>You are not using the latest version of BookBrowser. Your current version is " + currentVersion + ", but the latest version is " + latestVersion + ".</b><br><br>You can download the latest version <a href=\"https://github.com/geek1011/BookBrowser/releases/latest\" target=\"_blank\">here</a>.<br><br>The release notes for the versions up to " + latestVersion + " are below.<br><br>";
                message += releaseNotes;

                console.log(message);

                if (window.location.pathname == "/books/") {
                    picoModal('<h2>BookBrowser Update Available</h2>' + message + '<br><a type="button" target="_blank" href="https://github.com/geek1011/BookBrowser/releases/latest" class="btn btn-primary">Update</a>').show();
                }
            } catch (err) {
                console.warn(err);
            }
        }, function (err) {
            console.warn(err);
        });
    }, function (err) {
        console.warn(err);
    });
};

window.checkForUpdates(BookBrowserVersion);
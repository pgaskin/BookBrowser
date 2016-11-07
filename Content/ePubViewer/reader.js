EPUBJS.Hooks.register("beforeChapterDisplay").wgxpath = function (callback, renderer) {

    wgxpath.install(renderer.render.window);

    if (callback) callback();
};

EPUBJS.Hooks.register('beforeChapterDisplay').swipeDetection = function (callback, renderer) {
    function detectSwipe() {
        var script = renderer.doc.createElement('script');
        script.text = "\
      var swiper = new Hammer(document);\
      swiper.on('swipeleft', function() {\
        parent.Book.nextPage();\
      });\
      swiper.on('swiperight', function() {\
        parent.Book.prevPage();\
      });";
        renderer.doc.head.appendChild(script);
    }
    EPUBJS.core.addScript('http://geek1011.github.io/ePubViewer/epubjs/libs/hammer.min.js', detectSwipe, renderer.doc.head);
    if (callback) {
        callback();
    }
};

wgxpath.install(window);

EPUBJS.Hooks.register("beforeChapterDisplay").pageTurns = function (callback, renderer) {

    var lock = false;
    var arrowKeys = function (e) {
        e.preventDefault();
        if (lock) return;

        if (e.keyCode == 37) {
            Book.prevPage();
            lock = true;
            setTimeout(function () {
                lock = false;
            }, 100);
            return false;
        }

        if (e.keyCode == 39) {
            Book.nextPage();
            lock = true;
            setTimeout(function () {
                lock = false;
            }, 100);
            return false;
        }

    };
    renderer.doc.addEventListener('keydown', arrowKeys, false);
    if (callback) callback();
}

Book = null;
BookID = "";
BookToc = null;
appid = "ePubViewer"
initSettingsDone = false;

/* Toggles a sidebar and hides the others. Pass no argument to hide all sidebars */
doSidebar = function (sidebarName) {
    if (sidebarName != null) {
        var isHidden = document.getElementById("sidebar-" + sidebarName).classList.contains("hidden");
    }

    var sidebars = document.querySelectorAll(".reader > main > aside")
    for (var i = 0; i < sidebars.length; i++) {
        sidebars[i].classList.add("hidden");
    }

    if (sidebarName != null) {
        if (isHidden) {
            document.getElementById("sidebar-" + sidebarName).classList.remove("hidden");
        } else {
            document.getElementById("sidebar-" + sidebarName).classList.add("hidden");
        }
    }
}

doPrev = function () {
    Book.prevPage();
}

doNext = function () {
    Book.nextPage();
}

doCfi = function (cfi) {
    Book.gotoCfi(cfi);
}

doChapter = function (chaptercfi) {
    Book.displayChapter(chaptercfi);
}
getCoverAsDataURL = function (book, callback) {
    book.coverUrl().then(function (blobUrl) {
        console.log(blobUrl);
        var xhr = new XMLHttpRequest;
        xhr.responseType = 'blob';
        xhr.onload = function () {
            var recoveredBlob = xhr.response;
            var reader = new FileReader;
            reader.onload = function () {
                callback(reader.result);
            };
            reader.readAsDataURL(recoveredBlob);
        };
        xhr.open('GET', blobUrl);
        xhr.send();
    });
}


doBook = function (url) {
    var bookel = document.getElementById("book");
    bookel.innerHTML = '<div class="sk-fading-circle"> <div class="sk-circle1 sk-circle"></div><div class="sk-circle2 sk-circle"></div><div class="sk-circle3 sk-circle"></div><div class="sk-circle4 sk-circle"></div><div class="sk-circle5 sk-circle"></div><div class="sk-circle6 sk-circle"></div><div class="sk-circle7 sk-circle"></div><div class="sk-circle8 sk-circle"></div><div class="sk-circle9 sk-circle"></div><div class="sk-circle10 sk-circle"></div><div class="sk-circle11 sk-circle"></div><div class="sk-circle12 sk-circle"></div></div>';
    document.getElementById("curpercent").innerText = "";

    Book = ePub({
        storage: false
    });

    Book.on('book:loadFailed', function () {
        bookel.innerHTML = "<div class=\"message error\">Error loading book</div>";
    });

    Book.open(url);

    Book.getMetadata().then(function (meta) {
        try {
            Book.nextPage(); /* Fix first page not showing issue */
        } catch (e) {}
        document.title = meta.bookTitle + " – " + meta.creator;
        document.getElementById("booktitle").innerHTML = meta.bookTitle;
        document.getElementById("bookauthor").innerHTML = meta.creator;
        try {
            getCoverAsDataURL(Book, function (u) {
                document.getElementById("bookcover").src = u;
            })
        } catch (e) {}
        BookID = [meta.bookTitle, meta.creator, meta.identifier, meta.publisher].join(":");
        var curpostmp = localStorage.getItem(appid + "|" + BookID + "|curPosCfi");
        if (curpostmp) {
            Book.goto(curpostmp)
        }

        Book.on('renderer:locationChanged', function (locationCfi) {
            localStorage.setItem(appid + "|" + BookID + "|curPosCfi", Book.getCurrentLocationCfi())
        });

        Book.locations.generate().then(function () {
            doUpdateProgressIndicators();
            Book.on('renderer:locationChanged', function (locationCfi) {
                doUpdateProgressIndicators();
            });
        });

        Book.on('renderer:locationChanged', function (locationCfi) {
            try {
                var toclist = document.getElementById("toc-container").getElementsByClassName("toc-entry");
                for (var e = 0; e < toclist.length; e++) {
                    if (toclist[e].getAttribute("data-cfi") == "epubcfi(" + Book.currentChapter.cfiBase + ")") {
                        toclist[e].classList.add("active");
                    } else {
                        toclist[e].classList.remove("active");
                    }
                }
            } catch (e) {}
        });

        bookel.innerHTML = "";
        updateSettings();
        Book.renderTo(bookel);
        document.body.classList.remove("not-loaded")
    });

    Book.getToc().then(function (toc) {
        BookToc = toc;
        var containerel = document.getElementById("toc-container");
        for (var i = 0; i < toc.length; i++) {
            var entryel = document.createElement("a");
            entryel.classList.add("toc-entry");
            entryel.innerText = toc[i].label;
            entryel.setAttribute("data-cfi", toc[i].cfi);
            entryel.href = "javascript:void(0);";
            entryel.onclick = function (e) {
                doChapter(e.target.getAttribute("data-cfi"));
            };
            containerel.appendChild(entryel);
        }
    });
}

doFileFromFileObject = function (fileObj) {
    var reader = new FileReader();
    reader.addEventListener("load", function () {
        var arr = (new Uint8Array(reader.result)).subarray(0, 2);
        var header = "";
        for (var i = 0; i < arr.length; i++) {
            header += arr[i].toString(16);
        }
        console.log(header);
        if (header == "504b") {
            doBook(reader.result);
        } else {
            document.getElementById("book").innerHTML = "<div class=\"message error\">The file you chose is not a valid ePub ebook. Please try choosing a new file.</div>";
        }
    }, false);
    if (fileObj) {
        reader.readAsArrayBuffer(fileObj);
    }
}

doHandleFileInput = function (el) {
    var el = el || document.getElementById("bookChooser");
    doFileFromFileObject(el.files[0]);
}

/**
 * detect IE
 * returns version of IE or false, if browser is not Internet Explorer
 */
function detectIE() {
    var ua = window.navigator.userAgent;

    // Test values; Uncomment to check result …

    // IE 10
    // ua = 'Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Trident/6.0)';

    // IE 11
    // ua = 'Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko';

    // Edge 12 (Spartan)
    // ua = 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36 Edge/12.0';

    // Edge 13
    // ua = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586';

    var msie = ua.indexOf('MSIE ');
    if (msie > 0) {
        // IE 10 or older => return version number
        return parseInt(ua.substring(msie + 5, ua.indexOf('.', msie)), 10);
    }

    var trident = ua.indexOf('Trident/');
    if (trident > 0) {
        // IE 11 => return version number
        var rv = ua.indexOf('rv:');
        return parseInt(ua.substring(rv + 3, ua.indexOf('.', rv)), 10);
    }

    var edge = ua.indexOf('Edge/');
    if (edge > 0) {
        // Edge (IE 12+) => return version number
        return parseInt(ua.substring(edge + 5, ua.indexOf('.', edge)), 10);
    }

    // other browser
    return false;
}

checkCompatibility = function () {
    if (detectIE() === false) {} else {
        return false;
    }
    if (window.FileReader && window.FileReader.prototype.readAsArrayBuffer) {} else {
        return false;
    }
    if (document.createElement("p").style.flex == null) {
        return false;
    }
    if (document.createElement("p").style['flex-direction'] == null) {
        return false;
    }
    if (document.createElement("p").style['justify-content'] == null) {
        return false;
    }
    if (document.createElement("p").style['opacity'] == null) {
        return false;
    }
    if (document.createElement("p").style['white-space'] == null) {
        return false;
    }
    if (document.createElement("p").style['vertical-align'] == null) {
        return false;
    }
    if (document.createElement("p").style['min-width'] == null) {
        return false;
    }
    if (window.Int16Array) {} else {
        return false;
    }
    if (window.Worker) {} else {
        return false;
    }
    if (window.localStorage) {} else {
        return false;
    }
    return true;
}

elID = function (i) {
    return document.getElementById(i);
}

mappingsValueCSS = {
    'font-family': elID("family"),
    'font-size': elID("size"),
    'line-height': elID("lineheight"),
    'margin': elID("margin")
}
mappingCheckedInit = {
    'spreads': elID("spreads")
}
themeselect = elID("theme");
themes = {
    "white": {
        "bg": "white",
        "fg": "black"
    },
    "black": {
        "bg": "black",
        "fg": "white"
    },
    "darkgray": {
        "bg": "rgb(64,64,64)",
        "fg": "rgb(220,220,220)"
    },
    "sepia": {
        "bg": "wheat",
        "fg": "black"
    },
    "solarizedDark": {
        "bg": "#002b36",
        "fg": "#839496"
    },
    "solarizedLight": {
        "bg": "#fdf6e3",
        "fg": "#657b83"
    }
}

initSettings = function () {
    for (var i in mappingsValueCSS) {
        if (mappingsValueCSS.hasOwnProperty(i)) {
            if (localStorage[appid + "|" + i]) {
                mappingsValueCSS[i].value = localStorage[appid + "|" + i];
            } else {
                localStorage[appid + "|" + i] = mappingsValueCSS[i].value;
            }
        }
    }
    if (localStorage[appid + "|" + "spreads"]) {
        mappingCheckedInit["spreads"].checked = localStorage[appid + "|" + "spreads"];
    } else {
        localStorage[appid + "|" + "spreads"] = mappingCheckedInit["spreads"].checked;
    }

    if (localStorage[appid + "|" + "theme"]) {
        themeselect.value = localStorage[appid + "|" + "theme"];
    } else {
        localStorage[appid + "|" + "theme"] = themeselect.value;
    }
    initSettingsDone = true;
}

updateSettings = function () {
    for (var i in mappingsValueCSS) {
        if (mappingsValueCSS.hasOwnProperty(i)) {
            localStorage[appid + "|" + i] = mappingsValueCSS[i].value;
            console.log(i + ":" + mappingsValueCSS[i].value)
            Book.setStyle(i, mappingsValueCSS[i].value);
        }
    }
    localStorage[appid + "|" + "spreads"] = mappingCheckedInit["spreads"].checked;
    try {
        Book.setStyle("background-color", themes[themeselect.value].bg);
        Book.setStyle("color", themes[themeselect.value].fg);
    } catch (e) {
        console.error("Error applying theme", e)
    }
    localStorage[appid + "|" + "theme"] = themeselect.value;
}

doBookReset = function () {
    if (confirm("Do you want to reset the book position?")) {
        if (confirm("Are you sure?")) {
            delete localStorage[appid + "|" + BookID + "|curPosCfi"];
            location.reload();
        }
    }
}

doSettingsReset = function () {
    if (confirm("Do you want to reset the settings (this will not erase your book progress)?")) {
        if (confirm("Are you sure?")) {
            delete localStorage[appid + "|" + "spreads"];
            for (var i in mappingsValueCSS) {
                if (mappingsValueCSS.hasOwnProperty(i)) {
                    delete localStorage[appid + "|" + i];
                    mappingsValueCSS[i].value = mappingsValueCSS[i].defaultValue || mappingsValueCSS[i].querySelector("option[selected]").value;
                }
            }
            initSettings();
            updateSettings();
        }
    }
}

doAllReset = function () {
    if (confirm("Do you want to reset all your book progress and all settings?")) {
        if (confirm("Are you sure?")) {
            localStorage.clear();
            try {
                document.getElementById("bookChooser").value = "";
            } catch (ex) {

            }
            location.reload();
        }
    }
}

doUpdateProgressIndicators = function () {
    var progressint = Math.round(Book.locations.percentageFromCfi(Book.getCurrentLocationCfi()).toFixed(2) * 100);
    document.getElementById("curpercent").innerText = String(progressint) + "%";
    document.getElementById("bookprogresstext").innerText = String(progressint) + "% read";
    document.getElementById("bookprogressbar").setAttribute("value", String(progressint));
    document.getElementById("bookcurrentcfi").innerText = "Current cfi: " + Book.getCurrentLocationCfi();
    try {
        document.getElementById("currentchapter").innerText = "Chapter: " + BookToc[Book.currentChapter.spinePos].label;
    } catch (e) {
        document.getElementById("currentchapter").innerText = "";

    }
}
document.getElementById("book").innerHTML = "<div class=\"message info\">Please click the middle button on the toolbar below or <a href=\"javascript:void(0);\" onclick=\"document.getElementById('bookChooser').click()\">click here</a> to open a book.</div>";
if (checkCompatibility()) {

} else {
    alert("You are using an incompatible browser. Try using a different browser such as Google Chrome or Mozilla Firefox.");
    document.getElementById("book").innerHTML = "<div class=\"message error\">You are using an incompatible browser. Try using a different browser such as Google Chrome or Mozilla Firefox. If you think this was a mistake, then you can <a href=\"http://github.com/geek1011/ePubViewer/issues\">report an issue</a>.</div>";
    document.querySelector("nav").style.display = "none";
}
document.body.classList.add("not-loaded")

initSettings();

var ufn = location.search.replace("?", "") || location.hash.replace("#","");
if (ufn) {
    doBook(ufn);
} else {
    doHandleFileInput();
}

doSidebar();

(function nwjsfunctions() {
    if (typeof nw != "undefined") {
        var gui = require('nw.gui');
        var fs = require('fs');
        var uto = gui.App.argv[0];
        fs.stat(uto, function (err, stat) {
            if (err == null) {
                doBook("file://" + uto);
            } else if (err.code == 'ENOENT') {} else {}
        });
    }
})();

EPUBJS.Hooks.register("beforeChapterDisplay").wgxpath = function(callback, renderer) {

    wgxpath.install(renderer.render.window);

    if (callback) callback();
};

wgxpath.install(window);

EPUBJS.Hooks.register("beforeChapterDisplay").pageTurns = function(callback, renderer) {

    var lock = false;
    var arrowKeys = function(e) {
        e.preventDefault();
        if (lock) return;

        if (e.keyCode == 37) {
            Book.prevPage();
            lock = true;
            setTimeout(function() {
                lock = false;
            }, 100);
            return false;
        }

        if (e.keyCode == 39) {
            Book.nextPage();
            lock = true;
            setTimeout(function() {
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
appid = "ePubViewer"
initSettingsDone = false;

/* Toggles a sidebar and hides the others. Pass no argument to hide all sidebars */
doSidebar = function(sidebarName) {
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

doPrev = function() {
    Book.prevPage();
}

doNext = function() {
    Book.nextPage();
}

doCfi = function(cfi) {
    Book.gotoCfi(cfi);
}

doChapter = function(chaptercfi) {
    Book.displayChapter(chaptercfi);
}

doBook = function(url) {
    var bookel = document.getElementById("book");
    bookel.innerHTML = '<div class="sk-fading-circle"> <div class="sk-circle1 sk-circle"></div><div class="sk-circle2 sk-circle"></div><div class="sk-circle3 sk-circle"></div><div class="sk-circle4 sk-circle"></div><div class="sk-circle5 sk-circle"></div><div class="sk-circle6 sk-circle"></div><div class="sk-circle7 sk-circle"></div><div class="sk-circle8 sk-circle"></div><div class="sk-circle9 sk-circle"></div><div class="sk-circle10 sk-circle"></div><div class="sk-circle11 sk-circle"></div><div class="sk-circle12 sk-circle"></div></div>';
    document.getElementById("curpercent").innerText = "";

    Book = ePub({
        storage: false
    });

    Book.on('book:loadFailed', function() {
        bookel.innerHTML = "<div class=\"message error\">Error loading book</div>";
    });

    Book.open(url);

    Book.getMetadata().then(function(meta) {
        document.title = meta.bookTitle + " – " + meta.creator;
        BookID = [meta.bookTitle, meta.creator, meta.identifier, meta.publisher].join(":");
        var curpostmp = localStorage.getItem(appid + "|" + BookID + "|curPosCfi");
        if (curpostmp) {
            Book.goto(curpostmp)
        }
        Book.on('renderer:locationChanged', function(locationCfi) {
            localStorage.setItem(appid + "|" + BookID + "|curPosCfi", Book.getCurrentLocationCfi())
        });

        Book.locations.generate().then(function() {
            doUpdateProgressIndicators();
            Book.on('renderer:locationChanged', function(locationCfi) {
                doUpdateProgressIndicators();
            });
        });

        bookel.innerHTML = "";
        updateSettings();
        Book.renderTo(bookel);
        document.body.classList.remove("not-loaded")
    });

    Book.getToc().then(function(toc) {
        var containerel = document.getElementById("toc-container");
        for (var i = 0; i < toc.length; i++) {
            var entryel = document.createElement("a");
            entryel.classList.add("toc-entry");
            entryel.innerText = toc[i].label;
            entryel.setAttribute("data-cfi", toc[i].cfi);
            entryel.href = "javascript:void(0);";
            entryel.onclick = function(e) {
                doChapter(e.target.getAttribute("data-cfi"));
            };
            containerel.appendChild(entryel);
        }
    });
}

doFileFromFileObject = function(fileObj) {
    var reader = new FileReader();
    reader.addEventListener("load", function() {
        var arr = (new Uint8Array(reader.result)).subarray(0, 2);
        var header = "";
        for(var i = 0; i < arr.length; i++) {
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

doHandleFileInput = function(el) {
    var el = el || document.getElementById("bookChooser");
    doFileFromFileObject(el.files[0]);
}

checkCompatibility = function() {
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

elID = function(i) {
    return document.getElementById(i);
}

mappingsValueCSS = {
    'font-family': elID("family"),
    'font-size': elID("size"),
    'color': elID("color"),
    'background-color': elID("background"),
    'line-height': elID("lineheight"),
    'margin': elID("margin")
}
mappingCheckedInit = {
    'spreads': elID("spreads")
}

initSettings = function() {
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
    initSettingsDone = true;
}

updateSettings = function() {
    for (var i in mappingsValueCSS) {
        if (mappingsValueCSS.hasOwnProperty(i)) {
            localStorage[appid + "|" + i] = mappingsValueCSS[i].value;
            console.log(i + ":" + mappingsValueCSS[i].value)
            Book.setStyle(i, mappingsValueCSS[i].value);
        }
    }
    localStorage[appid + "|" + "spreads"] = mappingCheckedInit["spreads"].checked;
}

doBookReset = function() {
    if (confirm("Do you want to reset the book position?")) {
        if (confirm("Are you sure?")) {
            delete localStorage[appid + "|" + BookID + "|curPosCfi"];
            location.reload();
        }
    }
}

doSettingsReset = function() {
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

doAllReset = function() {
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

doUpdateProgressIndicators = function() {
    document.getElementById("curpercent").innerText = String(Book.locations.percentageFromCfi(Book.getCurrentLocationCfi()).toFixed(2) * 100) + "%";
}

if (checkCompatibility()) {

} else {
    alert("You are using an incompatible browser. Try using a modern browser such as Google Chrome or Mozilla Firefox.");
    document.getElementById("book").innerHTML = "<div class=\"message error\">You are using an incompatible browser. Try using a modern browser such as Google Chrome or Mozilla Firefox. If you think this was a mistake, then you can <a href=\"http://github.com/geek1011/ePubViewer/issues\">report an issue</a>.</div>";
    document.querySelector("nav").style.display = "none";
}
document.body.classList.add("not-loaded")
document.getElementById("book").innerHTML = "<div class=\"message info\">Please click the middle button on the toolbar below or <a href=\"javascript:void(0);\" onclick=\"document.getElementById('bookChooser').click()\">click here</a> to open a book.</div>";
initSettings();
doHandleFileInput();
doSidebar();

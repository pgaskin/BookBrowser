ePubViewer = {};
ePubViewer.state = {
    "loaded": false,
    "book-title": "",
    "current-page": 0,
    "total-pages": 0,
    "book-author": "",
    "cover-url": "",
    "book-id": "",
    "percent-read": 0,
    "current-cfi": "",
    "current-chapter": "",
    "book": null,
    "toc": []
};
ePubViewer.themes = {
    "SepiaLight": {
        "background-color": "#FBF0D9",
        "color": "#704214"
    },
    "SepiaDark": {
        "color": "#FBF0D9",
        "background-color": "#704214"
    },
    "White": {
        "color": "#000000",
        "background-color": "#FFFFFF"
    },
    "Black": {
        "background-color": "#000000",
        "color": "#FFFFFF"
    }
};
ePubViewer.fonts = {
    "ArbutusSlab": {
        "link": "https://fonts.googleapis.com/css?family=Arbutus+Slab",
        "font-family": "'Arbutus Slab', Georgia, serif"
    },
    "DroidSerif": {
        "link": "https://fonts.googleapis.com/css?family=Droid+Serif:400,400i,700,700i",
        "font-family": "'Droid Serif', Georgia, serif"
    },
    "OpenSans": {
        "link": "https://fonts.googleapis.com/css?family=Open+Sans:400,400i,700,700i",
        "font-family": "'Open Sans', Ubuntu, Trebuchet, sans-serif"
    }
}
ePubViewer.settings = {
    "theme": "White",
    "font": "OpenSans",
    "line-height": "1.5",
    "font-size": "11pt",
    "margin": "5%"
}
ePubViewer.elements = {};
ePubViewer.events = {};
ePubViewer.functions = {};
ePubViewer.functions.showFatalError = function (message) {
    ePubViewer.state = {
        "loaded": false,
        "book-title": "",
        "current-page": 0,
        "total-pages": 0,
        "book-author": "",
        "cover-url": "",
        "book-id": "",
        "percent-read": 0,
        "current-cfi": "",
        "current-chapter": "",
        "book": null,
        "toc": []
    };
    ePubViewer.elements.content.innerHTML = [
        "<div class=\"welcome\">",
        "<div class=\"welcome-inner\">",
        "<div class=\"title\">ePubViewer</div>",
        "<div class=\"menu\">",
        "A fatal error has occured: " + message,
        "</div>",
        "</div>",
        "</div>"
    ].join("\n");
};
ePubViewer.functions.updateIndicators = function (message) {
    if (ePubViewer.state.loaded) {
        try {
            ePubViewer.state["percent-read"] = Math.round(ePubViewer.state.book.locations.percentageFromCfi(ePubViewer.state.book.getCurrentLocationCfi()).toFixed(2) * 100);
        } catch (e) {}
        try {
            ePubViewer.state["current-cfi"] = ePubViewer.state.book.getCurrentLocationCfi();
        } catch (e) {}
        try {
            ePubViewer.state["current-chapter"] = ePubViewer.state.book.toc[ePubViewer.state.book.spinePos].label
        } catch (e) {}
        try {
            document.title = ePubViewer.state["book-title"] + " - " + ePubViewer.state["book-author"];
        } catch (e) {}
        if (ePubViewer.state.book.pagination.totalPages) {
            ePubViewer.state["current-page"] = ePubViewer.state.book.pagination.pageFromCfi(ePubViewer.state.book.getCurrentLocationCfi());
            ePubViewer.state["total-pages"] = ePubViewer.state.book.pagination.totalPages;
        }
    }

    var els = document.querySelectorAll("[data-text]");
    for (var i = 0; i < els.length; i++) {
        try {
            var nv = ePubViewer.state[els[i].getAttribute("data-text")];
            if (els[i].innerHTML != nv) {
                els[i].innerHTML = nv;
            }
        } catch (e) {}
    }

    var els = document.querySelectorAll("[data-href]");
    for (var i = 0; i < els.length; i++) {
        try {
            var nv = ePubViewer.state[els[i].getAttribute("data-href")];
            if (els[i].href != nv) {
                els[i].href = nv;
            }
        } catch (e) {}
    }

    var els = document.querySelectorAll("[data-src]");
    for (var i = 0; i < els.length; i++) {
        try {
            var nv = ePubViewer.state[els[i].getAttribute("data-src")];
            if (els[i].src != nv) {
                els[i].src = nv;
            }
        } catch (e) {}
    }

    var els = document.querySelectorAll("[data-if]");
    for (var i = 0; i < els.length; i++) {
        try {
            var e = ePubViewer.state[els[i].getAttribute("data-if")];
            if (e) {
                if (e != "" && e != 0 && e != false) {
                    els[i].classList.remove("hidden");
                } else {
                    els[i].classList.add("hidden");
                }
            } else {
                els[i].classList.add("hidden");
            }
        } catch (e) {}
    }

    var els = document.querySelectorAll("[data-if-not]");
    for (var i = 0; i < els.length; i++) {
        try {
            var e = ePubViewer.state[els[i].getAttribute("data-if-not")];
            if (e) {
                if (e != "" && e != 0 && e != false) {
                    els[i].classList.add("hidden");
                } else {
                    els[i].classList.remove("hidden");
                }
            } else {
                els[i].classList.remove("hidden");
            }
        } catch (e) {}
    }
};

ePubViewer.functions.loadSettings = function () {
    for (k in ePubViewer.settings) {
        if (ePubViewer.settings.hasOwnProperty(k)) {
            var v = localStorage.getItem("ePubViewer|settings|" + k);
            if (v !== null) {
                ePubViewer.settings[k] = v;
                console.log("Loaded setting: ", k, v);
            }

            var el = document.querySelector(".reader [data-setting=" + k + "]");
            if (el) {
                el.value = ePubViewer.settings[k];
                console.log("Updated setting chooser: ", el, k, v);
            }
        }
    }
    ePubViewer.functions.applySettings();
};

ePubViewer.functions.updateSettingsFromSelectors = function () {
    for (k in ePubViewer.settings) {
        if (ePubViewer.settings.hasOwnProperty(k)) {
            var el = document.querySelector(".reader [data-setting=" + k + "]");
            if (el) {
                if (el.value) {
                    var v = el.value;
                    if (el.tagName.toLowerCase() == "select") {
                        v = el.options[el.selectedIndex].value;
                    }
                    ePubViewer.settings[k] = v;
                    console.log("Updated setting: ", el, k, v);
                    ePubViewer.functions.saveSettings();
                }
            }
        }
    }
    ePubViewer.functions.applySettings();
};

ePubViewer.functions.saveSettings = function () {
    for (k in ePubViewer.settings) {
        if (ePubViewer.settings.hasOwnProperty(k)) {
            var v = ePubViewer.settings[k];
            localStorage.setItem("ePubViewer|settings|" + k, v)
            console.log("Saved setting: ", k, v);
        }
    }
};

ePubViewer.functions.applySettings = function () {
    var font = ePubViewer.fonts[ePubViewer.settings.font] || ePubViewer.fonts.ArbutusSlab;
    var theme = ePubViewer.themes[ePubViewer.settings.theme] || ePubViewer.themes.SepiaLight;

    try {
        var doc = ePubViewer.state.book.renderer.doc;
        if (doc.getElementById("ePubViewerSettings") === null) {
            doc.body.appendChild(doc.createElement("style")).id = "ePubViewerSettings";
        }
        var styleEl = doc.getElementById("ePubViewerSettings");
        styleEl.innerHTML = [
            "html, body {",
            "font-family: " + font["font-family"] + ";",
            "font-size: " + ePubViewer.settings["font-size"] + ";",
            "color: " + theme["color"] + " !important;",
            "background-color: " + theme["background-color"] + " !important;",
            "line-height: " + ePubViewer.settings["line-height"] + " !important;",
            "}",
            "p {",
            "font-family: " + font["font-family"] + " !important;",
            "font-size: " + ePubViewer.settings["font-size"] + " !important;",
            "}"
        ].join("\n");
        if (font.link) {
            var el = doc.body.appendChild(doc.createElement("link"));
            el.setAttribute("rel", "stylesheet");
            el.setAttribute("href", font.link)
        }
    } catch (e) {}

    if (document.getElementById("ePubViewerAppSettings") === null) {
        document.body.appendChild(document.createElement("style")).id = "ePubViewerAppSettings";
    }
    var styleEla = document.getElementById("ePubViewerAppSettings");
    styleEla.innerHTML = [
        ".reader {",
        "font-family: " + font["font-family"] + ";",
        "color: " + theme["color"] + ";",
        "background-color: " + theme["background-color"] + ";",
        "}",
        ".reader .main .content {",
        "margin: 0 " + ePubViewer.settings["margin"] + ";",
        "}"
    ].join("\n");
    if (font.link) {
        var el = document.body.appendChild(document.createElement("link"));
        el.setAttribute("rel", "stylesheet");
        el.setAttribute("href", font.link)
    }
};

ePubViewer.functions.getCoverURL = function (callback) {
    ePubViewer.state.book.coverUrl().then(function (blobUrl) {
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
ePubViewer.actions = {};
ePubViewer.actions.settingsReset = function() {
    for (k in ePubViewer.settings) {
        if (ePubViewer.settings.hasOwnProperty(k)) {
            try {
                delete localStorage["ePubViewer|settings|" + k];
            } catch (e) {}
            console.log("Deleted setting: ", k);
        }
    }
    window.location.reload();
};
ePubViewer.actions.allReset = function() {
    if(confirm("Really delete all settings and book progress information?")) {
        for (var i = 0; i < localStorage.length; i++){
            var k = localStorage.key(i);
            if (k.startsWith("ePubViewer|")) {
                try {
                    delete localStorage[k];
                } catch (e) {}
            }
        }
        window.location.reload();
    }
};
ePubViewer.actions.showSidebar = function (sbname) {
    var sbels = document.querySelectorAll(".reader [data-sidebar]");
    for (var i = 0; i < sbels.length; i++) {
        try {
            if (sbels[i].getAttribute("data-sidebar") == sbname) {
                sbels[i].classList.add("visible");
            } else {
                sbels[i].classList.remove("visible");
            }
        } catch (e) {}
    }
};
ePubViewer.actions.closeSidebars = function () {
    var sbels = document.querySelectorAll(".reader [data-sidebar]");
    for (var i = 0; i < sbels.length; i++) {
        sbels[i].classList.remove("visible");
    }
};
ePubViewer.actions.prevPage = function () {
    ePubViewer.state.book.prevPage();
};
ePubViewer.actions.nextPage = function () {
    ePubViewer.state.book.nextPage();
};
ePubViewer.actions.gotoChapter = function (chapter) {
    if (chapter.indexOf("epubcfi") > -1) {
        ePubViewer.state.book.gotoCfi(chapter);
    } else {
        ePubViewer.state.book.gotoHref(chapter);
    }
};
ePubViewer.actions.loadBook = function (urlOrArrayBuffer) {
    ePubViewer.elements.content.innerHTML = "";
    ePubViewer.state.book = ePub({
        "storage": false
    });
    ePubViewer.state.book.on("book:loadFailed", function () {
        ePubViewer.state.loaded = false;
        ePubViewer.functions.updateIndicators();
        ePubViewer.functions.showFatalError("Error loading book");
    });
    ePubViewer.state.book.open(urlOrArrayBuffer);

    ePubViewer.state.book.getMetadata().then(function (meta) {
        try {
            ePubViewer.state.book.nextPage();
        } catch (e) {}

        ePubViewer.state["book-title"] = meta.bookTitle;
        ePubViewer.state["book-author"] = meta.creator;

        try {
            ePubViewer.functions.getCoverURL(function (u) {
                ePubViewer.state["cover-url"] = u;
            });
        } catch (e) {}

        ePubViewer.state["book-id"] = [meta.bookTitle, meta.creator, meta.identifier, meta.publisher].join(":");

        var curpostmp = localStorage.getItem("ePubViewer|" + ePubViewer.state["book-id"] + "|curPosCfi");
        if (curpostmp) {
            ePubViewer.state.book.goto(curpostmp)
        }

        ePubViewer.state.book.on('renderer:locationChanged', function (locationCfi) {
            localStorage.setItem("ePubViewer|" + ePubViewer.state["book-id"] + "|curPosCfi", ePubViewer.state.book.getCurrentLocationCfi())
        });

        ePubViewer.state.book.locations.generate().then(function () {
            ePubViewer.functions.updateIndicators();
        });

        var w = 600;
        var h = 800;
        ePubViewer.state.book.generatePagination(w, h).then(function () {
            ePubViewer.functions.updateIndicators();
        });

        ePubViewer.state.book.on('renderer:locationChanged', function (locationCfi) {
            ePubViewer.functions.updateIndicators();
        });

        ePubViewer.state.book.getToc().then(function (toc) {
            ePubViewer.state.toc = toc;
            var containerel = document.querySelector(".reader .toc");
            for (var i = 0; i < toc.length; i++) {
                console.log(toc[i])
                var entryel = document.createElement("a");
                entryel.classList.add("toc-entry");
                entryel.innerText = toc[i].label;
                entryel.setAttribute("data-cfi", toc[i].href);
                entryel.href = "javascript:void(0);";
                entryel.onclick = function (e) {
                    ePubViewer.actions.gotoChapter(e.target.getAttribute("data-cfi"));
                };
                containerel.appendChild(entryel);
                if (toc[i].subitems) {
                    for (var j = 0; j < toc[i].subitems.length; j++) {
                        var entryel = document.createElement("a");
                        entryel.classList.add("toc-entry");
                        entryel.style.paddingLeft = "20px";
                        entryel.innerText = toc[i].subitems[j].label;
                        entryel.setAttribute("data-cfi", toc[i].subitems[j].href);
                        entryel.href = "javascript:void(0);";
                        entryel.onclick = function (e) {
                            ePubViewer.actions.gotoChapter(e.target.getAttribute("data-cfi"));
                        };
                        containerel.appendChild(entryel);
                    }
                }
            }
        });

        ePubViewer.state.book.on('renderer:locationChanged', function (locationCfi) {
            try {
                var toclist = document.querySelectorAll(".reader .toc .toc-entry");
                for (var e = 0; e < toclist.length; e++) {
                    if (toclist[e].getAttribute("data-cfi") == "epubcfi(" + ePubViewer.state.book.currentChapter.cfiBase + ")") {
                        toclist[e].classList.add("active");
                    } else {
                        toclist[e].classList.remove("active");
                    }
                }
            } catch (e) {}
        });

        ePubViewer.functions.applySettings();
    });

    ePubViewer.state.loaded = true;
    ePubViewer.functions.updateIndicators();
    ePubViewer.state.book.renderTo(ePubViewer.elements.content);
};
ePubViewer.actions.openBook = function () {
    var fi = document.createElement("input");
    fi.accept = "application/epub+zip";
    fi.style.display = "none";
    fi.type = "file";
    fi.onchange = function (event) {
        var reader = new FileReader();
        reader.addEventListener("load", function () {
            var arr = (new Uint8Array(reader.result)).subarray(0, 2);
            var header = "";
            for (var i = 0; i < arr.length; i++) {
                header += arr[i].toString(16);
            }
            console.log(header);
            if (header == "504b") {
                ePubViewer.actions.loadBook(reader.result);
            } else {
                ePubViewer.functions.showFatalError("The book you chose is not a valid epub book. Please try again.")
            }
        }, false);
        if (fi.files[0]) {
            reader.readAsArrayBuffer(fi.files[0]);
        }
    };
    document.body.appendChild(fi);
    fi.click();
};
ePubViewer.init = function () {
    ePubViewer.elements.content = document.querySelector(".reader .content");
    ePubViewer.elements.openButton = document.querySelector(".reader .header .open-button");

    EPUBJS.Hooks.register('beforeChapterDisplay').swipeDetection = function (callback, renderer) {
        function detectSwipe() {
            var script = renderer.doc.createElement('script');
            script.text = "\
      var swiper = new Hammer(document);\
      swiper.on('swipeleft', function() {\
        parent.ePubViewer.actions.nextPage();\
      });\
      swiper.on('swiperight', function() {\
        parent.ePubViewer.actions.nextPage();\
      });";
            renderer.doc.head.appendChild(script);
        }
        EPUBJS.core.addScript('http://geek1011.github.io/ePubViewer/epubjs/libs/hammer.min.js', detectSwipe, renderer.doc.head);
        if (callback) {
            callback();
        }
    };

    EPUBJS.Hooks.register("beforeChapterDisplay").pageTurnKey = function (callback, renderer) {
        var lock = false;
        var arrowKeys = function (e) {
            e.preventDefault();
            if (lock) return;

            if (e.keyCode == 37) {
                ePubViewer.actions.prevPage();
                lock = true;
                setTimeout(function () {
                    lock = false;
                }, 100);
                return false;
            }

            if (e.keyCode == 39) {
                ePubViewer.actions.nextPage();
                lock = true;
                setTimeout(function () {
                    lock = false;
                }, 100);
                return false;
            }

        };
        renderer.doc.addEventListener('keydown', arrowKeys, false);
        if (callback) callback();
    };

    EPUBJS.Hooks.register("beforeChapterDisplay").clickHalfPageTurn = function (callback, renderer) {
        renderer.doc.addEventListener('click', function (event) {
            try {
                if (event.target.tagName.toLowerCase() == "a") return;
                if (event.target.parentNode.tagName.toLowerCase() == "a") return;
            } catch (e) {}
            var x = event.clientX;
            var width = document.body.clientWidth;
            var third = width / 3;
            if (x < third) {
                ePubViewer.actions.prevPage();
            } else if (x > (third * 2)) {
                ePubViewer.actions.nextPage();
            }
        });
        if (callback) callback();
    };

    EPUBJS.Hooks.register("beforeChapterDisplay").noSelection = function (callback, renderer) {
        renderer.doc.body.appendChild(document.createElement("style")).innerHTML = [
            "* {",
            "    -webkit-user-select: none;",
            "    -moz-user-select: none;",
            "    -ms-user-select: none;",
            "    user-select: none;",
            "    -webkit-user-drag: none;",
            "    -moz-user-drag: none;",
            "    -ms-user-drag: none;",
            "    user-drag: none;",
            "}"
        ].join("\n");
        if (callback) callback();
    };

    EPUBJS.Hooks.register("beforeChapterDisplay").styles = function (callback, renderer) {
        renderer.doc.body.appendChild(document.createElement("style")).innerHTML = [
            "a:link, a:visited {",
            "    color: inherit;",
            "    background: rgba(0,0,0,0.05);",
            "}",
            "",
            "html {",
            "    line-height: 1.5;",
            "    column-rule: 1px inset rgba(0,0,0,0.05);",
            "}"
        ].join("\n");
        if (callback) callback();
    };

    EPUBJS.Hooks.register("beforeChapterDisplay").settings = function (callback, renderer) {
        ePubViewer.functions.applySettings();
        if (callback) callback();
    };

    ePubViewer.functions.updateIndicators();
    window.setInterval(ePubViewer.functions.updateIndicators, 1000);

    ePubViewer.functions.loadSettings();

    (function loadFromURL() {
        var ufn = location.search.replace("?", "") || location.hash.replace("#", "");
        if (ufn.startsWith("!")) {
            ufn = ufn.replace("!", "");
            document.getElementById("openbutton").style = "display: none !important";
        }
        if (ufn) {
            ePubViewer.actions.loadBook(ufn);
        }
    })();
    (function nwjsfunctions() {
        if (typeof nw != "undefined") {
            var gui = require('nw.gui');
            var fs = require('fs');
            var uto = gui.App.argv[0];
            fs.stat(uto, function (err, stat) {
                if (err == null) {
                    ePubViewer.actions.loadBook("file://" + uto);
                } else if (err.code == 'ENOENT') {} else {}
            });
        }
    })();
};
ePubViewer.init();
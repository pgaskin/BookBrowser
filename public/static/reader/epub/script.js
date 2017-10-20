(function(d){
  var c = " ", f = "flex", fw = "-webkit-"+f, e = d.createElement('b');
  try { 
    e.style.display = fw; 
    e.style.display = f; 
    c += (e.style.display == f || e.style.display == fw) ? f : "no-"+f; 
  } catch(ex) { 
    c += "no-"+f; 
  }
  d.documentElement.className += c; 
})(document);

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
        "color": "#704214",
        "light": true
    },
    "White": {
        "color": "#000000",
        "background-color": "#FFFFFF",
        "light": true
    },
    "Black": {
        "background-color": "#000000",
        "color": "#FFFFFF",
        "light": false
    },
    "Gray": {
        "background-color": "#333333",
        "color": "#EEEEEE",
        "light": false
    },
    "Dark": {
        "background-color": "#262c2e",
        "color": "#f0f2f3",
        "light": false
    },
    "SolarizedLight": {
        "background-color": "#fdf6e3",
        "color": "#657b83",
        "light": true
    },
    "SolarizedDark": {
        "color": "#839496",
        "background-color": "#002b36",
        "light": false
    },
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
    },
    "SourceCodePro": {
        "link": "https://fonts.googleapis.com/css?family=Source+Code+Pro:200,300,400,500,600,700,900",
        "font-family": "'Source Code Pro', 'Open Sans', sans-serif"
    },
    "SourceSansPro": {
        "link": "https://fonts.googleapis.com/css?family=Source+Sans+Pro:200,200i,300,300i,400,400i,600,600i,700,700i,900,900i&subset=cyrillic,cyrillic-ext,greek,greek-ext,latin-ext,vietnamese",
        "font-family": "'Source Sans Pro', sans-serif"
    }
};
ePubViewer.settings = {
    "theme": "White",
    "font": "OpenSans",
    "line-height": "1.5",
    "font-size": "11pt",
    "margin": "5%"
};
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
            ePubViewer.state["current-chapter"] = ePubViewer.state.book.toc[ePubViewer.state.book.spinePos].label;
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
    var i = 0;
    var nv = null;
    for (i = 0; i < els.length; i++) {
        try {
            nv = ePubViewer.state[els[i].getAttribute("data-text")];
            if (els[i].innerHTML != nv) {
                els[i].innerHTML = nv;
            }
        } catch (e) {}
    }

    els = document.querySelectorAll("[data-href]");
    for (i = 0; i < els.length; i++) {
        try {
            nv = ePubViewer.state[els[i].getAttribute("data-href")];
            if (els[i].href != nv) {
                els[i].href = nv;
            }
        } catch (e) {}
    }

    els = document.querySelectorAll("[data-src]");
    for (i = 0; i < els.length; i++) {
        try {
            nv = ePubViewer.state[els[i].getAttribute("data-src")];
            if (els[i].src != nv) {
                els[i].src = nv;
            }
        } catch (e) {}
    }

    els = document.querySelectorAll("[data-if]");
    for (i = 0; i < els.length; i++) {
        try {
            var ea = ePubViewer.state[els[i].getAttribute("data-if")];
            if (ea) {
                if (ea != "" && ea != 0 && ea != false) {
                    els[i].classList.remove("hidden");
                } else {
                    els[i].classList.add("hidden");
                }
            } else {
                els[i].classList.add("hidden");
            }
        } catch (ex) {}
    }

    els = document.querySelectorAll("[data-if-not]");
    for (i = 0; i < els.length; i++) {
        try {
            var eb = ePubViewer.state[els[i].getAttribute("data-if-not")];
            if (eb) {
                if (eb != "" && eb != 0 && eb != false) {
                    els[i].classList.add("hidden");
                } else {
                    els[i].classList.remove("hidden");
                }
            } else {
                els[i].classList.remove("hidden");
            }
        } catch (ex) {}
    }
};

ePubViewer.functions.loadSettings = function () {
    for (var k in ePubViewer.settings) {
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
    for (var k in ePubViewer.settings) {
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
    for (var k in ePubViewer.settings) {
        if (ePubViewer.settings.hasOwnProperty(k)) {
            var v = ePubViewer.settings[k];
            localStorage.setItem("ePubViewer|settings|" + k, v);
            console.log("Saved setting: ", k, v);
        }
    }
};

ePubViewer.functions.applySettings = function () {
    var font = ePubViewer.fonts[ePubViewer.settings.font] || ePubViewer.fonts.ArbutusSlab;
    var theme = ePubViewer.themes[ePubViewer.settings.theme] || ePubViewer.themes.SepiaLight;

    try {
        if (theme.light) {
            document.body.classList.remove("dark");
            document.body.classList.add("light");
        } else {
            document.body.classList.add("dark");
            document.body.classList.remove("light");
        }
    } catch (ex) {}

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
            "color: " + theme.color + " !important;",
            "background-color: " + theme["background-color"] + " !important;",
            "line-height: " + ePubViewer.settings["line-height"] + " !important;",
            "}",
            "p {",
            "font-family: " + font["font-family"] + " !important;",
            "font-size: " + ePubViewer.settings["font-size"] + " !important;",
            "}"
        ].join("\n");
        if (font.link) {
            if (doc.getElementById("ePubViewerFontLink") === null) {
                doc.body.appendChild(doc.createElement("link")).id = "ePubViewerFontLink";
            }
            var el = document.getElementById("ePubViewerFontLink");
            el.setAttribute("rel", "stylesheet");
            el.setAttribute("href", font.link);
        }
    } catch (e) {}

    if (document.getElementById("ePubViewerAppSettings") === null) {
        document.body.appendChild(document.createElement("style")).id = "ePubViewerAppSettings";
    }
    var styleEla = document.getElementById("ePubViewerAppSettings");
    styleEla.innerHTML = [
        ".reader {",
        "font-family: " + font["font-family"] + ";",
        "color: " + theme.color + ";",
        "background-color: " + theme["background-color"] + ";",
        "}",
        ".reader .main .content {",
        "margin: 5px " + ePubViewer.settings.margin + ";",
        "}",
        ".reader .main .sidebar.overlay {",
        "color: " + theme.color + ";",
        "background: " + theme["background-color"] + " !important;",
        "}",
    ].join("\n");
    if (font.link) {
        if (document.getElementById("ePubViewerAppFontLink") === null) {
            document.body.appendChild(document.createElement("link")).id = "ePubViewerAppFontLink";
        }
        var ela = document.getElementById("ePubViewerAppFontLink");
        ela.setAttribute("rel", "stylesheet");
        ela.setAttribute("href", font.link);
    }
};

ePubViewer.functions.getCoverURL = function (callback) {
    ePubViewer.state.book.coverUrl().then(function (blobUrl) {
        console.log(blobUrl);
        var xhr = new XMLHttpRequest();
        xhr.responseType = 'blob';
        xhr.onload = function () {
            var recoveredBlob = xhr.response;
            var reader = new FileReader();
            reader.onload = function () {
                callback(reader.result);
            };
            reader.readAsDataURL(recoveredBlob);
        };
        xhr.open('GET', blobUrl);
        xhr.send();
    });
};
ePubViewer.actions = {};
ePubViewer.actions.settingsReset = function () {
    for (var k in ePubViewer.settings) {
        if (ePubViewer.settings.hasOwnProperty(k)) {
            try {
                delete localStorage["ePubViewer|settings|" + k];
            } catch (e) {}
            console.log("Deleted setting: ", k);
        }
    }
    window.location.reload();
};
ePubViewer.actions.allReset = function () {
    if (confirm("Really delete all settings and book progress information?")) {
        for (var i = 0; i < localStorage.length; i++) {
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
    var sbels = document.querySelectorAll(".reader .sidebar [data-sidebar]");
    for (var i = 0; i < sbels.length; i++) {
        try {
            if (sbels[i].getAttribute("data-sidebar") == sbname) {
                sbels[i].classList.add("visible");
            } else {
                sbels[i].classList.remove("visible");
            }
        } catch (e) {}
    }
    var sb = document.querySelector(".reader .sidebar");
    sb.classList.add("visible");
    sb.classList.remove("hidden");
};
ePubViewer.actions.closeSidebars = function () {
    var sbels = document.querySelectorAll(".reader .sidebar [data-sidebar]");
    for (var i = 0; i < sbels.length; i++) {
        sbels[i].classList.remove("visible");
    }
    var sb = document.querySelector(".reader .sidebar");
    sb.classList.remove("visible");
    sb.classList.add("hidden");
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
ePubViewer.actions.doSearch = function(q) {
    return new Promise(function(resolve, reject) {
        var r = ePubViewer.elements.searchResults;

        r.innerHTML = "";

        var resultPromises = [];

        q = q.replace(/^\s+|\s+$/g,'');

        if (q.length < 3) {
            r.innerHTML = '<a class="result">Please enter at least 3 characters</a>';
            resolve([]);
            return;
        }

        for (var i = 0; i < ePubViewer.state.book.spine.length; i++) {
            var spineItem = ePubViewer.state.book.spine[i];
            resultPromises.push(new Promise(function(resolve, reject) {
                new Promise(function(resolve, reject) {
                    resolve(new EPUBJS.Chapter(spineItem, ePubViewer.state.book.store, ePubViewer.state.book.credentials));
                }).then(function(chapter) {
                    return new Promise(function(resolve, reject) {
                        chapter.load().then(function() {
                        resolve(chapter);
                        }).catch(reject);
                    });
                }).then(function(chapter) {
                    return Promise.resolve(chapter.find(q));
                }).then(function(result) {
                    resolve(result);
                });
            }));
        }
        Promise.all(resultPromises).then(function(results) {
            return new Promise(function(resolve, reject) {
                resolve(results);
                var mergedResults = [].concat.apply([], results);
                console.log(mergedResults);
                var max = mergedResults.length;
                max = max > 100 ? 100 : max;
                var fragment = document.createDocumentFragment()
                for (var i = 0; i < max; i++) {
                    try {
                        var er = document.createElement("a");
                        er.classList.add("result");
                        er.href = "javascript:void(0);";
                        er.addEventListener("click", function() {
                        console.log(this.getAttribute("data-location"));
                        ePubViewer.state.book.goto(this.getAttribute("data-location"));
                        });
                        er.setAttribute("data-location", mergedResults[i].cfi);
                        er.innerHTML = mergedResults[i].excerpt;
                        fragment.appendChild(er);
                    } catch (e) {
                        console.warn(e);
                    }
                }
                r.appendChild(fragment);
            });
        });
    });
};
ePubViewer.actions.loadBook = function (urlOrArrayBuffer) {
    ePubViewer.actions.clearSearch();
    ePubViewer.elements.tocList.innerHTML = "";

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
            ePubViewer.state.book.goto(curpostmp);
        }

        ePubViewer.state.book.on('renderer:locationChanged', function (locationCfi) {
            localStorage.setItem("ePubViewer|" + ePubViewer.state["book-id"] + "|curPosCfi", ePubViewer.state.book.getCurrentLocationCfi());
        });

        ePubViewer.state.book.locations.generate().then(function () {
            ePubViewer.functions.updateIndicators();
        });

        var ismobile = (/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent));
        if (!ismobile) {
            window.setTimeout(function () {
                var w = 600;
                var h = 800;
                ePubViewer.state.book.generatePagination(w, h).then(function () {
                    ePubViewer.functions.updateIndicators();
                });
            }, 1000);
        }

        ePubViewer.state.book.on('renderer:locationChanged', function (locationCfi) {
            ePubViewer.functions.updateIndicators();
        });

        ePubViewer.state.book.getToc().then(function (toc) {
            ePubViewer.state.toc = toc;
            var containerel = ePubViewer.elements.tocList;
            containerel.innerHTML = "";
            for (var i = 0; i < toc.length; i++) {
                console.log(toc[i]);
                var entryel = document.createElement("a");
                entryel.classList.add("toc-entry");
                entryel.innerHTML = toc[i].label;
                entryel.setAttribute("data-cfi", toc[i].href);
                entryel.href = "javascript:void(0);";
                entryel.onclick = function (e) {
                    ePubViewer.actions.gotoChapter(e.target.getAttribute("data-cfi"));
                };
                containerel.appendChild(entryel);
                if (toc[i].subitems) {
                    for (var j = 0; j < toc[i].subitems.length; j++) {
                        var entryela = document.createElement("a");
                        entryela.classList.add("toc-entry");
                        entryela.style.paddingLeft = "20px";
                        entryela.innerHTML = toc[i].subitems[j].label;
                        entryela.setAttribute("data-cfi", toc[i].subitems[j].href);
                        entryela.href = "javascript:void(0);";
                        entryela.onclick = function (e) {
                            ePubViewer.actions.gotoChapter(e.target.getAttribute("data-cfi"));
                        };
                        containerel.appendChild(entryela);
                        if (toc[i].subitems[j].subitems) {
                            for (var k = 0; k < toc[i].subitems[j].subitems.length; k++) {
                                var entryelb = document.createElement("a");
                                entryelb.classList.add("toc-entry");
                                entryelb.style.paddingLeft = "40px";
                                entryelb.innerHTML = toc[i].subitems[j].subitems[k].label;
                                entryelb.setAttribute("data-cfi", toc[i].subitems[j].subitems[k].href);
                                entryelb.href = "javascript:void(0);";
                                entryelb.onclick = function (e) {
                                    ePubViewer.actions.gotoChapter(e.target.getAttribute("data-cfi"));
                                };
                                containerel.appendChild(entryelb);
                            }
                        }
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
            } catch (ex) {}
        });

        ePubViewer.functions.applySettings();
    });

    ePubViewer.state.loaded = true;
    ePubViewer.functions.updateIndicators();
    ePubViewer.state.book.renderTo(ePubViewer.elements.content);
};
ePubViewer.actions.handleSearch = function() {
    ePubViewer.actions.doSearch(ePubViewer.elements.searchBox.value);
};
ePubViewer.actions.clearSearch = function() {
    ePubViewer.elements.searchResults.innerHTML = "";
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
                ePubViewer.functions.showFatalError("The book you chose is not a valid epub book. Please try again.");
            }
        }, false);
        if (fi.files[0]) {
            reader.readAsArrayBuffer(fi.files[0]);
        }
    };
    document.body.appendChild(fi);
    fi.click();
};
ePubViewer.actions.fullScreen = function () {
    document.fullscreenEnabled = document.fullscreenEnabled || document.mozFullScreenEnabled || document.documentElement.webkitRequestFullScreen;
    
    var requestFullscreen = function (element) {
        if (element.requestFullscreen) {
            element.requestFullscreen();
        } else if (element.mozRequestFullScreen) {
            element.mozRequestFullScreen();
        } else if (element.webkitRequestFullScreen) {
            element.webkitRequestFullScreen(Element.ALLOW_KEYBOARD_INPUT);
        }
    };

    if (document.fullscreenEnabled) {
        requestFullscreen(document.documentElement);
    }
};
ePubViewer.init = function () {
    ePubViewer.elements.content = document.querySelector(".reader .content");
    ePubViewer.elements.openButton = document.querySelector(".reader .header .open-button");
    ePubViewer.elements.searchResults = document.querySelector(".reader .search-results");
    ePubViewer.elements.searchBox = document.querySelector('.reader .search-box');
    ePubViewer.elements.tocList = document.querySelector(".reader .toc");

    EPUBJS.Hooks.register('beforeChapterDisplay').swipeDetection = function (callback, renderer) {
        var swiper = renderer.doc.createElement('script');
        swiper.innerHTML = 'function Swiper(f,g,h,k,l){var b=null,c=null;f.addEventListener("touchstart",function(a){b=a.touches[0].clientX;c=a.touches[0].clientY},!1);f.addEventListener("touchmove",function(a){if(b&&c){var d=b-a.touches[0].clientX;a=c-a.touches[0].clientY;var e=Math.abs(d)>Math.abs(a);e&&30>Math.abs(d)||!e&&30>Math.abs(a)||(e?0<d?g():h():0<a?k():l(),b=c=null)}},!1)};Swiper(document,function(){parent.ePubViewer.actions.nextPage()},function(){parent.ePubViewer.actions.prevPage()},function(){},function(){});';
        renderer.doc.head.appendChild(swiper);
        if (callback) callback();
    };

    EPUBJS.Hooks.register("beforeChapterDisplay").pageTurnKey = function (callback, renderer) {
        var lock = false;
        renderer.doc.addEventListener('keydown', function (e) {
            e.preventDefault();
            if (lock) return;
            var d = false;
            switch (e.keyCode) {
                case 37:
                    d = true;
                    ePubViewer.actions.prevPage();
                    break;
                case 39:
                    d = true;
                    ePubViewer.actions.nextPage();
                    break;
            }
            if (d) {
                setTimeout(function () {
                    lock = false;
                }, 100);
                return false;
            }
        }, false);
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
    window.clearTimeout(ePubViewerLoadError);
    document.body.parentElement.classList.remove("load-error");
};
ePubViewer.init();
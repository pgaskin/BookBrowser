(function(d){
  var c = " ", f = "flex", fw = "-webkit-"+f, e = d.createElement('b');
  try { 
    e.style.display = fw; 
    e.style.display = f; 
    c += (e.style.display == f || e.style.display == fw) ? f : "no-"+f; 
  } catch(e) { 
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
        if (theme["light"]) {
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
ePubViewer.actions.settingsReset = function () {
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

        var ismobile = (/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent));
        if (!ismobile) {
            var w = 600;
            var h = 800;
            ePubViewer.state.book.generatePagination(w, h).then(function () {
                ePubViewer.functions.updateIndicators();
            });
        }

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

    EPUBJS.Hooks.register('beforeChapterDisplay').swipeDetection = function (callback, renderer) {
        function detectSwipe() {
            var hammer = renderer.doc.createElement('script');
            hammer.text = '!function(a,b,c,d){"use strict";function e(a,b,c){return setTimeout(j(a,c),b)}function f(a,b,c){return Array.isArray(a)?(g(a,c[b],c),!0):!1}function g(a,b,c){var e;if(a)if(a.forEach)a.forEach(b,c);else if(a.length!==d)for(e=0;e<a.length;)b.call(c,a[e],e,a),e++;else for(e in a)a.hasOwnProperty(e)&&b.call(c,a[e],e,a)}function h(b,c,d){var e="DEPRECATED METHOD: "+c+"\n"+d+" AT \n";return function(){var c=new Error("get-stack-trace"),d=c&&c.stack?c.stack.replace(/^[^\(]+?[\n$]/gm,"").replace(/^\s+at\s+/gm,"").replace(/^Object.<anonymous>\s*\(/gm,"{anonymous}()@"):"Unknown Stack Trace",f=a.console&&(a.console.warn||a.console.log);return f&&f.call(a.console,e,d),b.apply(this,arguments)}}function i(a,b,c){var d,e=b.prototype;d=a.prototype=Object.create(e),d.constructor=a,d._super=e,c&&hb(d,c)}function j(a,b){return function(){return a.apply(b,arguments)}}function k(a,b){return typeof a==kb?a.apply(b?b[0]||d:d,b):a}function l(a,b){return a===d?b:a}function m(a,b,c){g(q(b),function(b){a.addEventListener(b,c,!1)})}function n(a,b,c){g(q(b),function(b){a.removeEventListener(b,c,!1)})}function o(a,b){for(;a;){if(a==b)return!0;a=a.parentNode}return!1}function p(a,b){return a.indexOf(b)>-1}function q(a){return a.trim().split(/\s+/g)}function r(a,b,c){if(a.indexOf&&!c)return a.indexOf(b);for(var d=0;d<a.length;){if(c&&a[d][c]==b||!c&&a[d]===b)return d;d++}return-1}function s(a){return Array.prototype.slice.call(a,0)}function t(a,b,c){for(var d=[],e=[],f=0;f<a.length;){var g=b?a[f][b]:a[f];r(e,g)<0&&d.push(a[f]),e[f]=g,f++}return c&&(d=b?d.sort(function(a,c){return a[b]>c[b]}):d.sort()),d}function u(a,b){for(var c,e,f=b[0].toUpperCase()+b.slice(1),g=0;g<ib.length;){if(c=ib[g],e=c?c+f:b,e in a)return e;g++}return d}function v(){return qb++}function w(b){var c=b.ownerDocument||b;return c.defaultView||c.parentWindow||a}function x(a,b){var c=this;this.manager=a,this.callback=b,this.element=a.element,this.target=a.options.inputTarget,this.domHandler=function(b){k(a.options.enable,[a])&&c.handler(b)},this.init()}function y(a){var b,c=a.options.inputClass;return new(b=c?c:tb?M:ub?P:sb?R:L)(a,z)}function z(a,b,c){var d=c.pointers.length,e=c.changedPointers.length,f=b&Ab&&d-e===0,g=b&(Cb|Db)&&d-e===0;c.isFirst=!!f,c.isFinal=!!g,f&&(a.session={}),c.eventType=b,A(a,c),a.emit("hammer.input",c),a.recognize(c),a.session.prevInput=c}function A(a,b){var c=a.session,d=b.pointers,e=d.length;c.firstInput||(c.firstInput=D(b)),e>1&&!c.firstMultiple?c.firstMultiple=D(b):1===e&&(c.firstMultiple=!1);var f=c.firstInput,g=c.firstMultiple,h=g?g.center:f.center,i=b.center=E(d);b.timeStamp=nb(),b.deltaTime=b.timeStamp-f.timeStamp,b.angle=I(h,i),b.distance=H(h,i),B(c,b),b.offsetDirection=G(b.deltaX,b.deltaY);var j=F(b.deltaTime,b.deltaX,b.deltaY);b.overallVelocityX=j.x,b.overallVelocityY=j.y,b.overallVelocity=mb(j.x)>mb(j.y)?j.x:j.y,b.scale=g?K(g.pointers,d):1,b.rotation=g?J(g.pointers,d):0,b.maxPointers=c.prevInput?b.pointers.length>c.prevInput.maxPointers?b.pointers.length:c.prevInput.maxPointers:b.pointers.length,C(c,b);var k=a.element;o(b.srcEvent.target,k)&&(k=b.srcEvent.target),b.target=k}function B(a,b){var c=b.center,d=a.offsetDelta||{},e=a.prevDelta||{},f=a.prevInput||{};(b.eventType===Ab||f.eventType===Cb)&&(e=a.prevDelta={x:f.deltaX||0,y:f.deltaY||0},d=a.offsetDelta={x:c.x,y:c.y}),b.deltaX=e.x+(c.x-d.x),b.deltaY=e.y+(c.y-d.y)}function C(a,b){var c,e,f,g,h=a.lastInterval||b,i=b.timeStamp-h.timeStamp;if(b.eventType!=Db&&(i>zb||h.velocity===d)){var j=b.deltaX-h.deltaX,k=b.deltaY-h.deltaY,l=F(i,j,k);e=l.x,f=l.y,c=mb(l.x)>mb(l.y)?l.x:l.y,g=G(j,k),a.lastInterval=b}else c=h.velocity,e=h.velocityX,f=h.velocityY,g=h.direction;b.velocity=c,b.velocityX=e,b.velocityY=f,b.direction=g}function D(a){for(var b=[],c=0;c<a.pointers.length;)b[c]={clientX:lb(a.pointers[c].clientX),clientY:lb(a.pointers[c].clientY)},c++;return{timeStamp:nb(),pointers:b,center:E(b),deltaX:a.deltaX,deltaY:a.deltaY}}function E(a){var b=a.length;if(1===b)return{x:lb(a[0].clientX),y:lb(a[0].clientY)};for(var c=0,d=0,e=0;b>e;)c+=a[e].clientX,d+=a[e].clientY,e++;return{x:lb(c/b),y:lb(d/b)}}function F(a,b,c){return{x:b/a||0,y:c/a||0}}function G(a,b){return a===b?Eb:mb(a)>=mb(b)?0>a?Fb:Gb:0>b?Hb:Ib}function H(a,b,c){c||(c=Mb);var d=b[c[0]]-a[c[0]],e=b[c[1]]-a[c[1]];return Math.sqrt(d*d+e*e)}function I(a,b,c){c||(c=Mb);var d=b[c[0]]-a[c[0]],e=b[c[1]]-a[c[1]];return 180*Math.atan2(e,d)/Math.PI}function J(a,b){return I(b[1],b[0],Nb)+I(a[1],a[0],Nb)}function K(a,b){return H(b[0],b[1],Nb)/H(a[0],a[1],Nb)}function L(){this.evEl=Pb,this.evWin=Qb,this.allow=!0,this.pressed=!1,x.apply(this,arguments)}function M(){this.evEl=Tb,this.evWin=Ub,x.apply(this,arguments),this.store=this.manager.session.pointerEvents=[]}function N(){this.evTarget=Wb,this.evWin=Xb,this.started=!1,x.apply(this,arguments)}function O(a,b){var c=s(a.touches),d=s(a.changedTouches);return b&(Cb|Db)&&(c=t(c.concat(d),"identifier",!0)),[c,d]}function P(){this.evTarget=Zb,this.targetIds={},x.apply(this,arguments)}function Q(a,b){var c=s(a.touches),d=this.targetIds;if(b&(Ab|Bb)&&1===c.length)return d[c[0].identifier]=!0,[c,c];var e,f,g=s(a.changedTouches),h=[],i=this.target;if(f=c.filter(function(a){return o(a.target,i)}),b===Ab)for(e=0;e<f.length;)d[f[e].identifier]=!0,e++;for(e=0;e<g.length;)d[g[e].identifier]&&h.push(g[e]),b&(Cb|Db)&&delete d[g[e].identifier],e++;return h.length?[t(f.concat(h),"identifier",!0),h]:void 0}function R(){x.apply(this,arguments);var a=j(this.handler,this);this.touch=new P(this.manager,a),this.mouse=new L(this.manager,a)}function S(a,b){this.manager=a,this.set(b)}function T(a){if(p(a,dc))return dc;var b=p(a,ec),c=p(a,fc);return b&&c?dc:b||c?b?ec:fc:p(a,cc)?cc:bc}function U(a){this.options=hb({},this.defaults,a||{}),this.id=v(),this.manager=null,this.options.enable=l(this.options.enable,!0),this.state=gc,this.simultaneous={},this.requireFail=[]}function V(a){return a&lc?"cancel":a&jc?"end":a&ic?"move":a&hc?"start":""}function W(a){return a==Ib?"down":a==Hb?"up":a==Fb?"left":a==Gb?"right":""}function X(a,b){var c=b.manager;return c?c.get(a):a}function Y(){U.apply(this,arguments)}function Z(){Y.apply(this,arguments),this.pX=null,this.pY=null}function $(){Y.apply(this,arguments)}function _(){U.apply(this,arguments),this._timer=null,this._input=null}function ab(){Y.apply(this,arguments)}function bb(){Y.apply(this,arguments)}function cb(){U.apply(this,arguments),this.pTime=!1,this.pCenter=!1,this._timer=null,this._input=null,this.count=0}function db(a,b){return b=b||{},b.recognizers=l(b.recognizers,db.defaults.preset),new eb(a,b)}function eb(a,b){this.options=hb({},db.defaults,b||{}),this.options.inputTarget=this.options.inputTarget||a,this.handlers={},this.session={},this.recognizers=[],this.element=a,this.input=y(this),this.touchAction=new S(this,this.options.touchAction),fb(this,!0),g(this.options.recognizers,function(a){var b=this.add(new a[0](a[1]));a[2]&&b.recognizeWith(a[2]),a[3]&&b.requireFailure(a[3])},this)}function fb(a,b){var c=a.element;c.style&&g(a.options.cssProps,function(a,d){c.style[u(c.style,d)]=b?a:""})}function gb(a,c){var d=b.createEvent("Event");d.initEvent(a,!0,!0),d.gesture=c,c.target.dispatchEvent(d)}var hb,ib=["","webkit","Moz","MS","ms","o"],jb=b.createElement("div"),kb="function",lb=Math.round,mb=Math.abs,nb=Date.now;hb="function"!=typeof Object.assign?function(a){if(a===d||null===a)throw new TypeError("Cannot convert undefined or null to object");for(var b=Object(a),c=1;c<arguments.length;c++){var e=arguments[c];if(e!==d&&null!==e)for(var f in e)e.hasOwnProperty(f)&&(b[f]=e[f])}return b}:Object.assign;var ob=h(function(a,b,c){for(var e=Object.keys(b),f=0;f<e.length;)(!c||c&&a[e[f]]===d)&&(a[e[f]]=b[e[f]]),f++;return a},"extend","Use `assign`."),pb=h(function(a,b){return ob(a,b,!0)},"merge","Use `assign`."),qb=1,rb=/mobile|tablet|ip(ad|hone|od)|android/i,sb="ontouchstart"in a,tb=u(a,"PointerEvent")!==d,ub=sb&&rb.test(navigator.userAgent),vb="touch",wb="pen",xb="mouse",yb="kinect",zb=25,Ab=1,Bb=2,Cb=4,Db=8,Eb=1,Fb=2,Gb=4,Hb=8,Ib=16,Jb=Fb|Gb,Kb=Hb|Ib,Lb=Jb|Kb,Mb=["x","y"],Nb=["clientX","clientY"];x.prototype={handler:function(){},init:function(){this.evEl&&m(this.element,this.evEl,this.domHandler),this.evTarget&&m(this.target,this.evTarget,this.domHandler),this.evWin&&m(w(this.element),this.evWin,this.domHandler)},destroy:function(){this.evEl&&n(this.element,this.evEl,this.domHandler),this.evTarget&&n(this.target,this.evTarget,this.domHandler),this.evWin&&n(w(this.element),this.evWin,this.domHandler)}};var Ob={mousedown:Ab,mousemove:Bb,mouseup:Cb},Pb="mousedown",Qb="mousemove mouseup";i(L,x,{handler:function(a){var b=Ob[a.type];b&Ab&&0===a.button&&(this.pressed=!0),b&Bb&&1!==a.which&&(b=Cb),this.pressed&&this.allow&&(b&Cb&&(this.pressed=!1),this.callback(this.manager,b,{pointers:[a],changedPointers:[a],pointerType:xb,srcEvent:a}))}});var Rb={pointerdown:Ab,pointermove:Bb,pointerup:Cb,pointercancel:Db,pointerout:Db},Sb={2:vb,3:wb,4:xb,5:yb},Tb="pointerdown",Ub="pointermove pointerup pointercancel";a.MSPointerEvent&&!a.PointerEvent&&(Tb="MSPointerDown",Ub="MSPointerMove MSPointerUp MSPointerCancel"),i(M,x,{handler:function(a){var b=this.store,c=!1,d=a.type.toLowerCase().replace("ms",""),e=Rb[d],f=Sb[a.pointerType]||a.pointerType,g=f==vb,h=r(b,a.pointerId,"pointerId");e&Ab&&(0===a.button||g)?0>h&&(b.push(a),h=b.length-1):e&(Cb|Db)&&(c=!0),0>h||(b[h]=a,this.callback(this.manager,e,{pointers:b,changedPointers:[a],pointerType:f,srcEvent:a}),c&&b.splice(h,1))}});var Vb={touchstart:Ab,touchmove:Bb,touchend:Cb,touchcancel:Db},Wb="touchstart",Xb="touchstart touchmove touchend touchcancel";i(N,x,{handler:function(a){var b=Vb[a.type];if(b===Ab&&(this.started=!0),this.started){var c=O.call(this,a,b);b&(Cb|Db)&&c[0].length-c[1].length===0&&(this.started=!1),this.callback(this.manager,b,{pointers:c[0],changedPointers:c[1],pointerType:vb,srcEvent:a})}}});var Yb={touchstart:Ab,touchmove:Bb,touchend:Cb,touchcancel:Db},Zb="touchstart touchmove touchend touchcancel";i(P,x,{handler:function(a){var b=Yb[a.type],c=Q.call(this,a,b);c&&this.callback(this.manager,b,{pointers:c[0],changedPointers:c[1],pointerType:vb,srcEvent:a})}}),i(R,x,{handler:function(a,b,c){var d=c.pointerType==vb,e=c.pointerType==xb;if(d)this.mouse.allow=!1;else if(e&&!this.mouse.allow)return;b&(Cb|Db)&&(this.mouse.allow=!0),this.callback(a,b,c)},destroy:function(){this.touch.destroy(),this.mouse.destroy()}});var $b=u(jb.style,"touchAction"),_b=$b!==d,ac="compute",bc="auto",cc="manipulation",dc="none",ec="pan-x",fc="pan-y";S.prototype={set:function(a){a==ac&&(a=this.compute()),_b&&this.manager.element.style&&(this.manager.element.style[$b]=a),this.actions=a.toLowerCase().trim()},update:function(){this.set(this.manager.options.touchAction)},compute:function(){var a=[];return g(this.manager.recognizers,function(b){k(b.options.enable,[b])&&(a=a.concat(b.getTouchAction()))}),T(a.join(" "))},preventDefaults:function(a){if(!_b){var b=a.srcEvent,c=a.offsetDirection;if(this.manager.session.prevented)return void b.preventDefault();var d=this.actions,e=p(d,dc),f=p(d,fc),g=p(d,ec);if(e){var h=1===a.pointers.length,i=a.distance<2,j=a.deltaTime<250;if(h&&i&&j)return}if(!g||!f)return e||f&&c&Jb||g&&c&Kb?this.preventSrc(b):void 0}},preventSrc:function(a){this.manager.session.prevented=!0,a.preventDefault()}};var gc=1,hc=2,ic=4,jc=8,kc=jc,lc=16,mc=32;U.prototype={defaults:{},set:function(a){return hb(this.options,a),this.manager&&this.manager.touchAction.update(),this},recognizeWith:function(a){if(f(a,"recognizeWith",this))return this;var b=this.simultaneous;return a=X(a,this),b[a.id]||(b[a.id]=a,a.recognizeWith(this)),this},dropRecognizeWith:function(a){return f(a,"dropRecognizeWith",this)?this:(a=X(a,this),delete this.simultaneous[a.id],this)},requireFailure:function(a){if(f(a,"requireFailure",this))return this;var b=this.requireFail;return a=X(a,this),-1===r(b,a)&&(b.push(a),a.requireFailure(this)),this},dropRequireFailure:function(a){if(f(a,"dropRequireFailure",this))return this;a=X(a,this);var b=r(this.requireFail,a);return b>-1&&this.requireFail.splice(b,1),this},hasRequireFailures:function(){return this.requireFail.length>0},canRecognizeWith:function(a){return!!this.simultaneous[a.id]},emit:function(a){function b(b){c.manager.emit(b,a)}var c=this,d=this.state;jc>d&&b(c.options.event+V(d)),b(c.options.event),a.additionalEvent&&b(a.additionalEvent),d>=jc&&b(c.options.event+V(d))},tryEmit:function(a){return this.canEmit()?this.emit(a):void(this.state=mc)},canEmit:function(){for(var a=0;a<this.requireFail.length;){if(!(this.requireFail[a].state&(mc|gc)))return!1;a++}return!0},recognize:function(a){var b=hb({},a);return k(this.options.enable,[this,b])?(this.state&(kc|lc|mc)&&(this.state=gc),this.state=this.process(b),void(this.state&(hc|ic|jc|lc)&&this.tryEmit(b))):(this.reset(),void(this.state=mc))},process:function(){},getTouchAction:function(){},reset:function(){}},i(Y,U,{defaults:{pointers:1},attrTest:function(a){var b=this.options.pointers;return 0===b||a.pointers.length===b},process:function(a){var b=this.state,c=a.eventType,d=b&(hc|ic),e=this.attrTest(a);return d&&(c&Db||!e)?b|lc:d||e?c&Cb?b|jc:b&hc?b|ic:hc:mc}}),i(Z,Y,{defaults:{event:"pan",threshold:10,pointers:1,direction:Lb},getTouchAction:function(){var a=this.options.direction,b=[];return a&Jb&&b.push(fc),a&Kb&&b.push(ec),b},directionTest:function(a){var b=this.options,c=!0,d=a.distance,e=a.direction,f=a.deltaX,g=a.deltaY;return e&b.direction||(b.direction&Jb?(e=0===f?Eb:0>f?Fb:Gb,c=f!=this.pX,d=Math.abs(a.deltaX)):(e=0===g?Eb:0>g?Hb:Ib,c=g!=this.pY,d=Math.abs(a.deltaY))),a.direction=e,c&&d>b.threshold&&e&b.direction},attrTest:function(a){return Y.prototype.attrTest.call(this,a)&&(this.state&hc||!(this.state&hc)&&this.directionTest(a))},emit:function(a){this.pX=a.deltaX,this.pY=a.deltaY;var b=W(a.direction);b&&(a.additionalEvent=this.options.event+b),this._super.emit.call(this,a)}}),i($,Y,{defaults:{event:"pinch",threshold:0,pointers:2},getTouchAction:function(){return[dc]},attrTest:function(a){return this._super.attrTest.call(this,a)&&(Math.abs(a.scale-1)>this.options.threshold||this.state&hc)},emit:function(a){if(1!==a.scale){var b=a.scale<1?"in":"out";a.additionalEvent=this.options.event+b}this._super.emit.call(this,a)}}),i(_,U,{defaults:{event:"press",pointers:1,time:251,threshold:9},getTouchAction:function(){return[bc]},process:function(a){var b=this.options,c=a.pointers.length===b.pointers,d=a.distance<b.threshold,f=a.deltaTime>b.time;if(this._input=a,!d||!c||a.eventType&(Cb|Db)&&!f)this.reset();else if(a.eventType&Ab)this.reset(),this._timer=e(function(){this.state=kc,this.tryEmit()},b.time,this);else if(a.eventType&Cb)return kc;return mc},reset:function(){clearTimeout(this._timer)},emit:function(a){this.state===kc&&(a&&a.eventType&Cb?this.manager.emit(this.options.event+"up",a):(this._input.timeStamp=nb(),this.manager.emit(this.options.event,this._input)))}}),i(ab,Y,{defaults:{event:"rotate",threshold:0,pointers:2},getTouchAction:function(){return[dc]},attrTest:function(a){return this._super.attrTest.call(this,a)&&(Math.abs(a.rotation)>this.options.threshold||this.state&hc)}}),i(bb,Y,{defaults:{event:"swipe",threshold:10,velocity:.3,direction:Jb|Kb,pointers:1},getTouchAction:function(){return Z.prototype.getTouchAction.call(this)},attrTest:function(a){var b,c=this.options.direction;return c&(Jb|Kb)?b=a.overallVelocity:c&Jb?b=a.overallVelocityX:c&Kb&&(b=a.overallVelocityY),this._super.attrTest.call(this,a)&&c&a.offsetDirection&&a.distance>this.options.threshold&&a.maxPointers==this.options.pointers&&mb(b)>this.options.velocity&&a.eventType&Cb},emit:function(a){var b=W(a.offsetDirection);b&&this.manager.emit(this.options.event+b,a),this.manager.emit(this.options.event,a)}}),i(cb,U,{defaults:{event:"tap",pointers:1,taps:1,interval:300,time:250,threshold:9,posThreshold:10},getTouchAction:function(){return[cc]},process:function(a){var b=this.options,c=a.pointers.length===b.pointers,d=a.distance<b.threshold,f=a.deltaTime<b.time;if(this.reset(),a.eventType&Ab&&0===this.count)return this.failTimeout();if(d&&f&&c){if(a.eventType!=Cb)return this.failTimeout();var g=this.pTime?a.timeStamp-this.pTime<b.interval:!0,h=!this.pCenter||H(this.pCenter,a.center)<b.posThreshold;this.pTime=a.timeStamp,this.pCenter=a.center,h&&g?this.count+=1:this.count=1,this._input=a;var i=this.count%b.taps;if(0===i)return this.hasRequireFailures()?(this._timer=e(function(){this.state=kc,this.tryEmit()},b.interval,this),hc):kc}return mc},failTimeout:function(){return this._timer=e(function(){this.state=mc},this.options.interval,this),mc},reset:function(){clearTimeout(this._timer)},emit:function(){this.state==kc&&(this._input.tapCount=this.count,this.manager.emit(this.options.event,this._input))}}),db.VERSION="2.0.6",db.defaults={domEvents:!1,touchAction:ac,enable:!0,inputTarget:null,inputClass:null,preset:[[ab,{enable:!1}],[$,{enable:!1},["rotate"]],[bb,{direction:Jb}],[Z,{direction:Jb},["swipe"]],[cb],[cb,{event:"doubletap",taps:2},["tap"]],[_]],cssProps:{userSelect:"none",touchSelect:"none",touchCallout:"none",contentZooming:"none",userDrag:"none",tapHighlightColor:"rgba(0,0,0,0)"}};var nc=1,oc=2;eb.prototype={set:function(a){return hb(this.options,a),a.touchAction&&this.touchAction.update(),a.inputTarget&&(this.input.destroy(),this.input.target=a.inputTarget,this.input.init()),this},stop:function(a){this.session.stopped=a?oc:nc},recognize:function(a){var b=this.session;if(!b.stopped){this.touchAction.preventDefaults(a);var c,d=this.recognizers,e=b.curRecognizer;(!e||e&&e.state&kc)&&(e=b.curRecognizer=null);for(var f=0;f<d.length;)c=d[f],b.stopped===oc||e&&c!=e&&!c.canRecognizeWith(e)?c.reset():c.recognize(a),!e&&c.state&(hc|ic|jc)&&(e=b.curRecognizer=c),f++}},get:function(a){if(a instanceof U)return a;for(var b=this.recognizers,c=0;c<b.length;c++)if(b[c].options.event==a)return b[c];return null},add:function(a){if(f(a,"add",this))return this;var b=this.get(a.options.event);return b&&this.remove(b),this.recognizers.push(a),a.manager=this,this.touchAction.update(),a},remove:function(a){if(f(a,"remove",this))return this;if(a=this.get(a)){var b=this.recognizers,c=r(b,a);-1!==c&&(b.splice(c,1),this.touchAction.update())}return this},on:function(a,b){var c=this.handlers;return g(q(a),function(a){c[a]=c[a]||[],c[a].push(b)}),this},off:function(a,b){var c=this.handlers;return g(q(a),function(a){b?c[a]&&c[a].splice(r(c[a],b),1):delete c[a]}),this},emit:function(a,b){this.options.domEvents&&gb(a,b);var c=this.handlers[a]&&this.handlers[a].slice();if(c&&c.length){b.type=a,b.preventDefault=function(){b.srcEvent.preventDefault()};for(var d=0;d<c.length;)c[d](b),d++}},destroy:function(){this.element&&fb(this,!1),this.handlers={},this.session={},this.input.destroy(),this.element=null}},hb(db,{INPUT_START:Ab,INPUT_MOVE:Bb,INPUT_END:Cb,INPUT_CANCEL:Db,STATE_POSSIBLE:gc,STATE_BEGAN:hc,STATE_CHANGED:ic,STATE_ENDED:jc,STATE_RECOGNIZED:kc,STATE_CANCELLED:lc,STATE_FAILED:mc,DIRECTION_NONE:Eb,DIRECTION_LEFT:Fb,DIRECTION_RIGHT:Gb,DIRECTION_UP:Hb,DIRECTION_DOWN:Ib,DIRECTION_HORIZONTAL:Jb,DIRECTION_VERTICAL:Kb,DIRECTION_ALL:Lb,Manager:eb,Input:x,TouchAction:S,TouchInput:P,MouseInput:L,PointerEventInput:M,TouchMouseInput:R,SingleTouchInput:N,Recognizer:U,AttrRecognizer:Y,Tap:cb,Pan:Z,Swipe:bb,Pinch:$,Rotate:ab,Press:_,on:m,off:n,each:g,merge:pb,extend:ob,assign:hb,inherit:i,bindFn:j,prefixed:u});var pc="undefined"!=typeof a?a:"undefined"!=typeof self?self:{};pc.Hammer=db,"function"==typeof define&&define.amd?define(function(){return db}):"undefined"!=typeof module&&module.exports?module.exports=db:a[c]=db}(window,document,"Hammer");';
            renderer.doc.head.appendChild(hammer);
            var script = renderer.doc.createElement('script');
            script.text = "var swiper = new Hammer(document);swiper.on('swipeleft', function() {parent.ePubViewer.actions.nextPage();});swiper.on('swiperight', function() {parent.ePubViewer.actions.nextPage();});";
            renderer.doc.head.appendChild(script);
        }
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
    window.clearTimeout(ePubViewerLoadError);
    document.body.parentElement.classList.remove("load-error");
};
ePubViewer.init();
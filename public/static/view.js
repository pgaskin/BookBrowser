var v = document.querySelector(".current-view");
var d = "cards";
var lb = document.querySelector(".view-list");
var cb = document.querySelector(".view-cards");

var getCurrentView = function() {
    var c = "";
    if (v.classList.contains("list")) {
        c = "list";
    } if (v.classList.contains("cards")) {
        c = "cards";
    } else {
        c = d.toString();
    }
};

var setCurrentView = function(c) {
    v.classList.remove("list");
    v.classList.remove("cards");
    v.classList.add(c);

    if (cb) {
        if (c == "cards") {
            cb.classList.add("active");
        } else {
            cb.classList.remove("active");
        }
    }

    if (lb) {
        if (c == "list") {
            lb.classList.add("active");
        } else {
            lb.classList.remove("active");
        }
    }
};

lb.addEventListener("click", function() {
    setCurrentView("list");
});

cb.addEventListener("click", function() {
    setCurrentView("cards");
});

setCurrentView(d);
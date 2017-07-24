if (document.body.className.indexOf("no-contains-view")>-1) {
    console.log("No view detected");
} else {
    var cvbtn = document.querySelector(".view-buttons .cards");
    var lvbtn = document.querySelector(".view-buttons .list");
    var vlst = document.querySelector(".view");

    var listView = function () {
        cvbtn.classList.remove("active");
        lvbtn.classList.add("active");
        vlst.classList.remove("cards");
        vlst.classList.add("list");
        localStorage.setItem("view", "list");
    }

    var cardsView = function () {
        lvbtn.classList.remove("active");
        cvbtn.classList.add("active");
        vlst.classList.remove("list");
        vlst.classList.add("cards");
        localStorage.setItem("view", "cards");
    }

    var restoreView = function () {
        if (localStorage in window) {
            var v = localStorage.getItem("view");
            if (v !== null) {
                if (v == "list") {
                    listView();
                } else {
                    cardsView();
                }
            } else {
                cardsView();
            }
        } else {
            cardsView();
        }
    }

    lvbtn.addEventListener("click", listView);
    cvbtn.addEventListener("click", cardsView);

    restoreView();
}
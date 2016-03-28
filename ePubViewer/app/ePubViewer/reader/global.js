var $ = function(sel) {
    return document.querySelectorAll(sel)[0];
};
Array.prototype.contains = function(obj) {
    var i = this.length;
    while (i--) {
        if (this[i] === obj) {
            return true;
        }
    }
    return false;
}
Array.prototype.toggle = function (v) {
    var i = this.indexOf(v);
    if (i === -1) {
        this.push(v);
    } else {
        this.splice(i,1);
    }
    return this;
}
Array.prototype.remove = function (v) {
    var i = this.indexOf(v);
    if (i === -1) {
    } else {
        this.splice(i,1);
    }
    return this;
}
Array.prototype.clean = function(deleteValue) {
  for (var i = 0; i < this.length; i++) {
    if (this[i] == deleteValue) {         
      this.splice(i, 1);
      i--;
    }
  }
  return this;
};

function getUrlVars() {
var vars = {};
var parts = window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
vars[key] = value;
});
return vars;
}

var toggleMenu = function() {
    $('#reader').classList.toggle('menu');
}

var setMenu = function(state) {
    if (state) {
        $('#reader').classList.add('menu');
    } else {
        $('#reader').classList.remove('menu')
    }
}

var book = ePub();
var lastError = null;
var curpos = null;
var paginationloaded = false;

var handleError = function(err) {
    lastError = err;
    alert('Error: ' + (lastError.message || lastError.stack || 'unknown error'));
    console.log(lastError[0]);
}

var getID = function() {
    return document.title;
}

var getBookmarks = function() {
    var bkms = [];
    var tmp = localStorage.getItem(getID() + '::bookmarks') || '';
    bkms = tmp.split(',');
    return bkms.clean('');
}

var setBookmarks = function(arx) {
    return localStorage.setItem(getID() + '::bookmarks', arx.join(','));
}

var getBookmark = function(position) {
    return getBookmarks().contains(position);
}

var setBookmark = function(position, bookmarked) {
    var bkms = getBookmarks();
    if (bookmarked) {
        if (! bkms.contains(position)) {
            bkms.push(position)
            setBookmarks(bkms);
        }
    } else {
        setBookmarks(bkms.remove(position));
    }
}

var toggleBookmark = function(position) {
    var bkms = getBookmarks();
    setBookmarks(bkms.toggle(position));
}

var updateBookmarksList = function() {
    console.log('Updating bookmarks list');
    var bkms = getBookmarks();
    $('#bookmarkspanel').innerText = '';
    for (var i = 0; i < bkms.length; i++) {
        var it = document.createElement('a');
        it.classList.add('item');
        it.innerText = 'Bookmark ' + i;
        it.href = 'javascript:void(0);';
        it.setAttribute('data-cfi', bkms[i]);
        it.addEventListener('click', function(e) {
            console.log('Going to bookmark');
            book.goto(e.target.getAttribute('data-cfi'));
            setMenu(false);
        });
        $('#bookmarkspanel').appendChild(it);
    }
}

var updateFooter = function() {
    if (paginationloaded) {
        $('#currentpage').innerText = 'Page ' + book.pagination.pageFromCfi(curpos) + ' of ' + book.pagination.totalPages + ' - ' + Math.round(book.pagination.percentageFromCfi(curpos) * 100) + '% read';
    }
};

var updateBookmarkIcon = function() {
    console.log('Updating bookmark icon');
    if (getBookmark(curpos)) {
        $('#bookmarkbutton').innerText = 'bookmark';
    } else {
        $('#bookmarkbutton').innerText = 'bookmark_border';
    }
    updateBookmarksList();
}

var p = book.open(decodeURIComponent(getUrlVars()['url']) || '').then(function() {
    book.renderTo('bookcontent');
}).then(function() {
    $('#backbutton').addEventListener('click', function() {console.log('Previous page');book.prevPage();});
    $('#nextbutton').addEventListener('click', function() {console.log('Next page');book.nextPage();});
    $('#bookmarkspanelbutton').addEventListener('click', function() {console.log('Bookmarks panel');$('#menu').className = 'bookmarks';});
    $('#tocpanelbutton').addEventListener('click', function() {console.log('TOC panel');$('#menu').className = 'toc';});
    $('#menubutton').addEventListener('click', toggleMenu);
    $('#bookmarkbutton').addEventListener('click', function() {console.log('Toggling bookmark status');toggleBookmark(curpos);updateBookmarkIcon();});
    $('#inner').addEventListener('click', function() {setMenu(false);});
    book.getToc().then(function(toc) {
        $('#tocpanel').innerText = '';
        for (var i = 0; i < toc.length; i++) {
            var it = document.createElement('a');
            it.classList.add('item');
            it.innerText = toc[i].label;
            it.href = 'javascript:void(0);';
            it.setAttribute('data-cfi', toc[i].cfi);
            it.addEventListener('click', function(e) {
                console.log('Going to TOC entry');
                book.goto(e.target.getAttribute('data-cfi'));
                setMenu(false);
            });
            $('#tocpanel').appendChild(it);
        }
    });
    window.setInterval(function(){$('#inner').getElementsByTagName('iframe')[0].contentDocument.body.parentNode.addEventListener('click', function() {setMenu(false);});}, 10);
    book.generatePagination().then(function() {paginationloaded = true;console.log('Pagination loaded');updateFooter();});
    
    return book.getMetadata();
}).then(function(meta) {
    document.title = meta.bookTitle + ' - ' + meta.creator;
    $('#meta').innerText = document.title;
    try {
        book.goto(localStorage.getItem(getID()));
    } catch (e) {
        console.log("Could not restore book pos");
    }
    
    book.on('renderer:locationChanged', function(location) {
        console.log("Updating book pos");
        curpos = location;
        localStorage.setItem(getID(), location);
        updateBookmarkIcon();
        updateFooter();
    });
});



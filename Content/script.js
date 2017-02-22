if (typeof bookList == "undefined") {
	alert("books.js not found. Please run indexer.py then createjson.py");
}
var Books = bookList;
var IdMappings = {};
var Lists = {
  series: {},
  authors: {}
};


var ProcessEntries = function () {
  var HashIt = function (txt) {
    var jss = new jsSHA("SHA-512", "TEXT");
    jss.update(txt);
    return jss.getHash("HEX");
  };

  for (var i = 0; i < Books.length; i++) {
    Books[i].bookId = HashIt(Books[i].title + Books[i].author + Books[i].series + Books[i].seriesIndex).substring(0, 8);
    IdMappings[Books[i].bookId] = Books[i].title || "";

    Books[i].authorId = HashIt(Books[i].author || "").substring(0, 8);
    IdMappings[Books[i].authorId] = Books[i].author || "";
    Lists.authors[Books[i].authorId] = Lists.authors[Books[i].authorId] || {
      name: Books[i].author,
      count: 0
    };
    Lists.authors[Books[i].authorId].count += 1;

    Books[i].seriesId = HashIt(Books[i].series || "").substring(0, 8);
    IdMappings[Books[i].seriesId] = Books[i].series || "";
    if (Books[i].series) {
      Lists.series[Books[i].seriesId] = Lists.series[Books[i].seriesId] || {
        name: Books[i].series,
        count: 0
      };
      Lists.series[Books[i].seriesId].count += 1;
    }

    try {
      var epubdownload = "";
      for (var j = 0; j < Books[i].downloads.length; j++) {
        if (Books[i].downloads[j].type == "epub") {
          Books[i].readLink = "ePubViewer/index.html#../" + Books[i].downloads[j].link;
        }
      }
    } catch (e) {}
  }
};
ProcessEntries();

BookList = {
  props: ["books", "urlFiltered", "title", "IdMappings", "searching"],
  data: function () {
    return {
      searchText: null
    };
  },
  template: '<div><div class="pagetitle" v-if="urlFiltered">{{ IdMappings[$route.params.filterText] }}</div><div class="pagetitle" v-else>{{ title }}</div><div class="searchbar" v-if="searching == true"><input type="text" name="search" id="search" class="search" placeholder="Search..." v-model="searchText" /><button name="searchbutton" id="searchbutton" class="searchbutton"><i class="fa fa-search"></i></button></div><br><div class="booklist"><template v-for="book of books"><div class="book" v-bind:key="book.bookId" v-if="((typeof(urlFiltered) !== \'undefined\' && urlFiltered) ? (book[urlFiltered] == $route.params.filterText) : ((typeof(searching) !== \'undefined\' && searching == true) ? ((book.title || \'\').toLowerCase().indexOf((searchText || \'\').toLowerCase()) > -1 && true || (book.author || \'\').toLowerCase().indexOf((searchText || \'\').toLowerCase()) > -1 || (book.series || \'\').toLowerCase().indexOf((searchText || \'\').toLowerCase())) > -1 : true))"><router-link v-bind:to="\'/book/\' + book.bookId" class="book-cover-container"><img class="book-cover" v-bind:src="book.coverURL" /></router-link><div class="book-info"><router-link v-bind:to="\'/book/\' + book.bookId" class="book-info-title">{{ book.title }}</router-link><router-link v-bind:to="\'/author/\' + book.authorId" class="book-info-author">{{ book.author }}</router-link></div></div></template></div></div>'
};

BookView = {
  props: ["books"],
  template: '<div class="bookview"><template v-for="book of books"><div class="book" v-bind:key="book.bookId" v-if="book.bookId == $route.params.bookId"><router-link v-bind:to="\'/book/\' + book.bookId" class="book-cover-container"><img class="book-cover" v-bind:src="book.coverURL" /></router-link><div class="book-info"><div v-bind:to="\'/book/\' + book.bookId" class="book-info-title">{{ book.title }}</div><router-link v-bind:to="\'/author/\' + book.authorId" class="book-info-author">{{ book.author }}</router-link><br><br><div class="book-info-series" v-if="book.series">Series: <router-link  v-bind:to="\'/series/\' + book.seriesId">{{ book.series }}</router-link> - {{ book.seriesIndex }}</div><br><div class="book-info-description"><b>Description:</b> <div v-html="book.description"></div></div><br><br><div class="book-info-download"><a class="download-link" v-bind:href="download.link" v-for="download in book.downloads"><i class="fa fa-download"></i>{{ download.type.toUpperCase() }}</a></div><br><div class="book-info-download"><a class="download-link" v-bind:href="book.readLink" v-if="book.readLink"><i class="fa fa-book"></i>Read Online</a></div></div></div></div></div></div>'
};

List = {
  props: ["lists", "title", "listName", "linkPrefix"],
  template: '<div><div class="pagetitle">{{ title }}</div><ul class="list"><li v-for="(item,id) in lists[listName]"><router-link v-bind:to="linkPrefix + id">{{ item.name }}</router-link><span class="count-badge">{{ item.count }}</span></li></ul></div>'
};


Router = new VueRouter({
  routes: [{
    props: {
      books: Books,
      title: "Books",
      IdMappings: IdMappings,
      searching: false
    },
    path: "/books",
    component: BookList,
    title: "Books"
  }, {
    props: {
      books: Books
    },
    path: "/book/:bookId",
    component: BookView
  }, {
    props: {
      books: Books,
      urlFiltered: "authorId",
      IdMappings: IdMappings,
      searching: false
    },
    path: "/author/:filterText",
    component: BookList,
  }, {
    props: {
      lists: Lists,
      title: "Authors",
      listName: "authors",
      linkPrefix: "/author/"
    },
    path: "/authors",
    component: List,
    title: "Authors"
  }, {
    props: {
      books: Books,
      urlFiltered: "seriesId",
      IdMappings: IdMappings,
      searching: false
    },
    path: "/series/:filterText",
    component: BookList
  }, {
    props: {
      lists: Lists,
      title: "Series",
      listName: "series",
      linkPrefix: "/series/"
    },
    path: "/series",
    component: List,
    title: "Series"
  }, {
    props: {
      books: Books,
      IdMappings: IdMappings,
      searching: true,
      title: "Search Results"
    },
    path: "/search",
    component: BookList,
    title: "Search"
  }, {
    path: "/random",
    redirect: function (orig) {
      return "/book/" + Books[Math.floor(Math.random() * Books.length)].bookId;
    },
    title: "Random"
  }, {
    path: "/",
    redirect: function (orig) {
      return "/books/";
    },
    title: "Books"
  }]
});

Router.beforeEach(function (to, from, next) {
  var newTitle = "Books";
  if (to.params.filterText) {
    newTitle = IdMappings[to.params.filterText];
  }
  if (to.meta.title) {
    newTitle = to.meta.title;
  }
  document.title = newTitle;
  next();
});

BookBrowser = new Vue({
  router: Router
}).$mount('#BookBrowser');

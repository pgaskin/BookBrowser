#!/usr/bin/env python2.7
import os

sdir = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Content");
if not os.path.exists(sdir):
    os.makedirs(sdir);

bdir = os.path.join(sdir, "Books");
if not os.path.exists(bdir):
    os.makedirs(bdir);

books = [];
for subdir, dirs, files in os.walk(bdir):
    for file in files:
        fname = os.path.join(subdir, file);
        if fname.endswith(".epub"):
            books.append(fname);

bookinfos = [];
for bfile in books:
    bookinfo = {"filename": None, "coverurl": None, "title": None, "author": None, "series": None, "seriesindex": None};

    bookinfo["filename"] = bfile.replace(sdir, "");
    bookinfo["filename"] = bookinfo["filename"].lstrip(os.sep);

    if os.path.isfile(bfile + ".png"):
        bookinfo["coverurl"] = bfile.replace(sdir, "")  + ".png";
    elif os.path.isfile(bfile + ".jpg"):
        bookinfo["coverurl"] = bfile.replace(sdir, "") + ".jpg";
    elif os.path.isfile(bfile + ".jpeg"):
        bookinfo["coverurl"] = bfile.replace(sdir, "") + ".jpeg";
    else:
        bookinfo["coverurl"] = "nocover.jpg".replace(sdir, "");
    bookinfo["coverurl"] = bookinfo["coverurl"].lstrip(os.sep);

    if os.path.isfile(bfile + ".txt"):
        with open(bfile + ".txt") as f:
            for line in f:
                try:
                    sp = line.split("=", 1);
                    sp[1] = sp[1].rstrip("\n");
                    bookinfo[sp[0]] = sp[1];
                except:
                    pass;
                
    else:
        try:
            bookinfo["title"] = bfile.rsplit(os.sep, 1)[1];
        except:
            bookinfo["title"] = "undefined"

        bookinfo["author"] = "";

    bookinfos.append(bookinfo);

booktemplate = '''{notfirst}{
    "downloads": [{
         type: "epub",
         link: "{filename}"
    }],
    "coverURL": "{coverurl}",
    "title": "{title}",
    "author": "{author}",
    "description": "{description}",
    "series": "{series}",
    "seriesIndex": "{seriesindex}"
}'''

json = '''['''

first = True;
for bookinfo in bookinfos:
    gbooktemplate = booktemplate;
    notfirst = "" if first else ", ";
    gbooktemplate = gbooktemplate.replace("{notfirst}", notfirst);
    for k in ["filename", "coverurl", "title", "author", "description", "series", "seriesindex"]:
        try: 
            gbooktemplate = gbooktemplate.replace("{" + k + "}", (bookinfo[k] or "").replace("\\", "/").replace('"', '\\"').replace('\n', ' ').replace('\r', ''));
        except:
            gbooktemplate = gbooktemplate.replace("{" + k + "}", ("").replace("\\", "/").replace('"', '\\"').replace('\n', ' ').replace('\r', ''));
    json = json + gbooktemplate;
    first = False;

json = json + ''']'''

jsonf = open(os.path.join(sdir, "books.json"), "w");
jsonf.write(json);
jsonf.close();

js = "bookList = " + json + ";"
jsf = open(os.path.join(sdir, "books.js"), "w");
jsf.write(js);
jsf.close();

print json;

    

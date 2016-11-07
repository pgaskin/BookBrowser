#!/usr/bin/env python2.7
import os, zipfile, lxml
import elementtree.ElementTree as etree

sdir = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Content");
if not os.path.exists(sdir):
    os.makedirs(sdir);

bdir = os.path.join(sdir, "Books");
if not os.path.exists(bdir):
    os.makedirs(bdir);

errorlist = [];
iserrs = False;
def errorprint(t):
    errorlist.append(t);
    iserrs = True;

print("Using dir: " + bdir);

books = [];
for subdir, dirs, files in os.walk(bdir):
    for file in files:
        fname = os.path.join(subdir, file);
        if fname.endswith(".epub"):
            books.append(fname);

for bfile in books:
    hascovercache = os.path.isfile(bfile + ".png") or os.path.isfile(bfile + ".jpg") or os.path.isfile(bfile + ".jpeg");
    hasmetacache = os.path.isfile(bfile + ".txt");

    zf = None;
    contentopf = None;
    contentopfroot = None;
    contentopfpath = None;
    metadataelement = None;
    manifestelement = None;

    if not hascovercache or not hasmetacache:
        try:
            zf = zipfile.ZipFile(bfile);
            containerxml = zf.read("META-INF/container.xml")
            root = etree.fromstring(containerxml);
            contentopf = zf.read(root[0][0].get("full-path"));
            contentopfpath = root[0][0].get("full-path");
            contentopfroot = etree.fromstring(contentopf);
            metadataelement = [e for e in contentopfroot if e.tag == "{http://www.idpf.org/2007/opf}metadata"][0];
            manifestelement = [e for e in contentopfroot if e.tag == "{http://www.idpf.org/2007/opf}manifest"][0];
        except Exception, e:
            errorprint("Error processing book " + bfile + ": " + str(e));
            continue;

    if hascovercache:
        print("Already processed cover for " + bfile + ". Delete the cover image to recache it.");
    else: 
        print("Processing cover for" + bfile);
        try:
            coverurl = [e for e in manifestelement if e.get("id") == [f for f in metadataelement if f.get("name") == "cover"][0].get("content")][0].get("href");
            cp = contentopfpath.rsplit(os.sep, 1)[0];
            if cp == contentopfpath:
                cp = "";
            coverurl = os.path.join(cp, coverurl);
            with open(bfile + "." + coverurl.rsplit(".", 1)[1], 'wb') as f:
                w = zf.read(coverurl);
                f.write(w);
            
        except Exception, e:
            errorprint("Error processing cover for" + bfile + ": " + str(e));

    if hasmetacache:
        print("Already processed metadata for " + bfile + ". Delete the metadata file to recache it.");
    else: 
        print("Processing metadata for" + bfile);
        try:
            fi = open(bfile + ".txt", 'wb');
            nl = "\n";
            try: 
                btitle = [f for f in metadataelement if f.tag == "{http://purl.org/dc/elements/1.1/}title"][0].text;
                fi.write("title=" + btitle + nl);
            except:
                errorprint("Error getting title for book " + bfile);
            
            try:
                bauthor = [f for f in metadataelement if f.tag == "{http://purl.org/dc/elements/1.1/}creator"][0].text;
                fi.write("author=" + bauthor + nl);
            except:
                errorprint("Error getting author for book " + bfile);
            fi.close();
        except Exception, e:
            errorprint("Error getting metadata for book " + bfile + ": " + str(e));



print("");
print("");
print("");
print("===============")
print("Done");
print("Errors:");
for err in errorlist:
    print(err);

if iserrs == False:
    print "none";

import SimpleHTTPServer
import SocketServer
import os

sdir = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Content");
if not os.path.exists(sdir):
    os.makedirs(sdir);

bdir = os.path.join(sdir, "Books");
if not os.path.exists(bdir):
    os.makedirs(bdir);

os.chdir(sdir);

PORT = 8000;

class SimpleHTTPRequestHandlerWithCORS(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_my_headers()

        SimpleHTTPServer.SimpleHTTPRequestHandler.end_headers(self)

    def send_my_headers(self):
        self.send_header("Access-Control-Allow-Origin", "*")

Handler = SimpleHTTPRequestHandlerWithCORS;

httpd = SocketServer.TCPServer(("", PORT), Handler);

print "serving at port", PORT;
httpd.serve_forever();

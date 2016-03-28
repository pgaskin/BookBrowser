'use strict';
const electron = require('electron');
const app = electron.app;
const BrowserWindow = electron.BrowserWindow;
let mainWindow;

function createWindow() {
    mainWindow = new BrowserWindow({
        width: 600,
        height: 600,
        'min-width': 350,
        'min-height': 400,
        nodeIntegration: true
    });
    //mainWindow.setMenu(null); // Remove this line to show the default menubar
    //mainWindow.toggleDevTools();
    mainWindow.loadURL('file://' + __dirname + '/main.html'); // If your app is hosted online, then replace this with the url
    mainWindow.on('closed', function() {
        mainWindow = null;
    });
}
app.on('ready', createWindow);
// MAC OSX FIXES
app.on('window-all-closed', function() {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});
app.on('activate', function() {
    if (mainWindow === null) {
        createWindow();
    }
});

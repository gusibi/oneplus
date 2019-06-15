import os
import time
from PyQt5.QtCore import QUrl
from PyQt5 import QtCore, QtGui, QtWidgets, QtWebEngineWidgets


class PageShotter(QtWebEngineWidgets.QWebEngineView):
    def __init__(self, url, outfile, width=576, height=474, parent=None):
        super(PageShotter, self).__init__(parent)
        self.loadFinished.connect(self.save)
        self.loadProgress.connect(self.loadProgressHandler)
        self.url = url
        self.width = width
        self.height = height
        self.outfile = outfile
        self.setAttribute(QtCore.Qt.WA_DontShowOnScreen, True)
        self.setAttribute(QtCore.Qt.WA_DeleteOnClose, True)
        self.show()
        settings = QtWebEngineWidgets.QWebEngineSettings.globalSettings()
        for attr in (
            QtWebEngineWidgets.QWebEngineSettings.PluginsEnabled,
            QtWebEngineWidgets.QWebEngineSettings.ScreenCaptureEnabled,
        ):
            settings.setAttribute(attr, True)

    @QtCore.pyqtSlot(int)
    def loadProgressHandler(self, prog):
        print(time.time(), ":load progress", prog)

    def shot(self):
        print("size: ", self.size())
        # if self.size().isNull():
        #     self.resize(640, 500)
        self.resize(self.width, self.height)
        print(self.url)
        # self.load(QUrl(self.url))
        self.load(QUrl.fromLocalFile(self.url))
        print(">>>>>>>>>>>>")
        # self.run_js()

    def run_js(self):
        script = """
            var scroll = function (dHeight) {
            var t = document.documentElement.scrollTop
            var h = document.documentElement.scrollHeight
            dHeight = dHeight || 0
            var current = t + dHeight
            if (current > h) {
                window.scrollTo(0, document.documentElement.clientHeight)
              } else {
                window.scrollTo(0, current)
              }
            }
        """
        command = script + "\n scroll({})".format(self.height())
        self.page().runJavaScript(command)

    @QtCore.pyqtSlot(bool)
    def save(self, finished):
        print("finished:", finished)
        # time.sleep(2)
        if finished:
            size = self.contentsRect()
            print(u"width：%d，hight：%d" % (size.width(), size.height()))
            img = QtGui.QImage(size.width(), size.height(), QtGui.QImage.Format_ARGB32)
            self.image = img
            painter = QtGui.QPainter(img)
            self.render(painter)
            print("painter start ")
            # time.sleep(2)
            painter.end()
            print("painter end ")
            print(time.time(), ": show html")
            filename = self.outfile
            if img.save(filename):
                # time.sleep(2)
                print("save end")
                # time.sleep(2)
                filepath = os.path.join(os.path.dirname(__file__), filename)
                print(u"success：%s" % filepath)
            else:
                print(u"fail")

        else:
            print("Error")
        self.close()


if __name__ == "__main__":
    import sys

    # 54
    app = QtWidgets.QApplication(sys.argv)
    shotter = PageShotter("http://127.0.0.1:8080/render/html", "code.png")
    # shotter = PageShotter("/Users/gs/Desktop/high.html", "high.png")
    shotter = PageShotter(
        "/Users/gs/Desktop/high1.html", "high1.png", height=30 + 22 * 1
    )
    shotter = PageShotter(
        "/Users/gs/Desktop/high2.html", "high2.png", height=30 + 22 * 2
    )
    shotter = PageShotter(
        "/Users/gs/Desktop/high3.html", "high3.png", height=30 + 22 * 3
    )
    shotter = PageShotter(
        "/Users/gs/Desktop/high5.html", "high5.png", height=30 + 22 * 5
    )
    shotter = PageShotter(
        "/Users/gs/Desktop/high10.html", "high10.png", height=30 + 21.2 * 10
    )
    shotter = PageShotter(
        "/Users/gs/Desktop/high20.html", "high20.png", height=30 + 21.2 * 20
    )
    shotter = PageShotter(
        "/Users/gs/Desktop/high50.html", "high50.png", height=30 + 21.1 * 50
    )
    shotter = PageShotter(
        "/Users/gs/Desktop/high100.html", "high100.png", height=30 + 21.1 * 100
    )
    # shotter = PageShotter("/Users/gs/Desktop/325-han.html", "325.png")
    # shotter = PageShotter("http://www.zaih.com", "zaih.png")
    # shotter = PageShotter("http://127.0.0.1:8080/render/html", "code.png")
    # shotter = PageShotter("http://www.baidu.com", "baidu.png")
    # shotter = PageShotter(
    #     "http://service-g5235zgh-1254035985.ap-beijing.apigateway.myqcloud.com/test/render/89392e0cdefbc3f9bac9cdddd5154f1e",
    #     "hello.png",
    # )
    shotter.shot()
    app.exec()


import sys
import PyQt5 as pyqt
from PyQt5 import QtGui as qtgui, QtCore as qtc, QtWidgets as qtw
from PyQt5.QtGui import QColor


# Main inherits from QtGui.QWindow class
class Main(qtw.QMainWindow):

    textArea = None
    menuBar = None
    sideBar = None

    def __init__(self, parent=None):
        # init superclass 
        super().__init__(parent)
        self.initUI()

    def initUI(self):
        # set main window size
        self.setFixedSize(900, 600)
        self.textArea = qtw.QTextEdit(self)
        self.loadText("../go/test.txt")
        self.setCentralWidget(self.textArea)

    def loadText(self, file):
        document = qtgui.QTextDocument(self.textArea)
        content = ''
        with open(file) as f:
            content = f.read()
            f.close()
        document.setPlainText(content)

        self.textArea.setReadOnly(True)
        self.textArea.setTextColor(QColor("black"))
        self.textArea.setFixedSize(600, 600)
        self.textArea.setDocument(document)

        
def main():
    app = qtw.QApplication(sys.argv)
    mainWin = Main()
    mainWin.show()
    sys.exit(app.exec_())


if __name__ == '__main__':
    main()
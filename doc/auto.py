# -*- coding: utf-8 -*-
import sys
import os
import logging
import traceback
import mock

sys.path.insert(0, os.path.abspath('.'))

ModulePath = "modules"
if os.path.isfile("./logs"):
    os.remove("logs")
logging.basicConfig(filename="logs", level=logging.DEBUG)

sys.setrecursionlimit(2000)

ListAdded = []
repToIndex = []
RepertoireToSearchOn = []

RepertoireToSearchOn.append("../sources/webcore/opt/webcore/libexec")
RepertoireToSearchOn.append("../sources/amqp2engines/opt/amqp2engines/engines")
RepertoireToSearchOn.append("../sources/canolibs/lib/canolibs")
RepertoireToSearchOn.append("../sources/canolibs/lib/canolibs/unittest")
RepertoireToSearchOn.append("../sources/pyperfstore2/pyperfstore2")
RepertoireToSearchOn.append("../sources/celery-libs/lib/celerylibs")
RepertoireToSearchOn.append("../sources/collectd-libs/opt/collectd-libs")
RepertoireToSearchOn.append("../sources/apscheduler-libs/lib/apschedulerlibs")
RepertoireToSearchOn.append("../sources/snmp2amqp/opt/snmp2amqp")
RepertoireToSearchOn.append("../sources/gelf2amqp/opt/gelf2amqp")
#RepertoireToSearchOn.append("../sources/ccli/opt/ccli")
RepertoireToSearchOn.append("../sources/ccli/opt/ccli/libexec")
RepertoireToSearchOn.append("../sources/ics2amqp/opt/ics2amqp")
RepertoireToSearchOn.append("../sources/wkhtmltopdf-libs/lib/wkhtmltopdf")
RepertoireToSearchOn.append(
    "../sources/amqp2engines/opt/amqp2engines/unittest")
RepertoireToSearchOn.append("../sources/amqp2engines/opt/amqp2engines")
RepertoireToSearchOn.append("../sources/pyperfstore2/test")
RepertoireToSearchOn.append("../sources/canotools/opt/canotools")

#Crash sphinx in a weird way
#RepertoireToSearchOn.append("../sources/pyperfstore2")

#RepertoireToSearchOn.append("/home/mine/mock/canopsis/sources/canotools/opt/canotools/")
#sys.modules["utils.py"] = mock.MagicMock()


def AddAll(path):
    for root, dirs, files in os.walk(path):
        for f in files:
            if f[-3:] == ".py":
                RepertoireToSearchOn.append(os.path.join(root))
                break


def printlog(TopPrint, start=""):
    #Refl = " %s '%s' ligne : %s : %s" % (start, inspect.stack()[2][4][0].strip(), inspect.stack()[1][2], TopPrint )
    #print Refl
    #logging.info(Refl)
    pass


def RewriteFile(FileName, result):
    w = open(FileName, "w")
    w.write(result)
    w.close()


def ReplaceImport(FileName, Import, Pat):
    printlog("ModifieFile : " + FileName + " " + Pat[4:])
    try:
        f = open(FileName, "r")

        Contenu = f.read().splitlines()
        result = ""
        TempPat = Pat.replace(" ", "")

        for line in Contenu:
            TempLine = line.replace(" ", "")

            if TempLine == TempPat:
                printlog(
                    line + " remplacé par : " + Import + " dans " + FileName)
                line = "import " + Import

            result += "\n" + line

        f.close()

        RewriteFile(FileName, result)

    except IOError:
        printlog(FileName + " Doesn't exist!")


def AddImport(FileName, Import):
    printlog("AddImport : " + Import + " , to  : " + FileName)
    try:
        f = open(FileName, "r")

        Contenu = f.read().splitlines()
        result = ""
        CHECK = True

        for line in Contenu:
            if CHECK and line.find("import") != -1:
                CHECK = False
                line = line + "\nimport " + Import

            result += "\n" + line

        f.close()
        RewriteFile(FileName, result)

    except IOError:
        printlog(FileName + " Doesn't exist!")


def CommentAnnotation(FileName):
   # printlog(FileName)
    try:
        f = open(FileName, "r")

        Contenu = f.read().splitlines()
        result = ""
        start = True

        for line in Contenu:
            if len(line) > 1 and line[0] == "@":
                line = "#" + line
                #line = line[1:]
                printlog(line)

            if(start is True):
                result += line
                start = False

            else:
                result += "\n" + line

        f.close()
        RewriteFile(FileName, result)

    except IOError:
        printlog(FileName + " Doesn't exist!")


def CheckFileName(FilePath, FileName):
    if FileName.find("-") != -1:
        NewName = FileName.replace("-", "")
        os.system(
            "mv " + os.path.join(FilePath, FileName) + " " + os.path.join(
                FilePath, NewName))


def AutoApi():
    os.system("rm -r " + ModulePath)
    os.system("mkdir " + ModulePath)
    for rep in RepertoireToSearchOn:
        for files in os.listdir(rep):
            if files[-3:] == ".py":
                CommentAnnotation(os.path.join(rep, files))
                CheckFileName(rep, files)

        repName = ModulePath + "/" + rep.split("/")[-1]
        if repName in set(repToIndex):
            suffix = "2"
            printlog(
                repName + " already in so changed to :" + repName + suffix)
            repName = repName + suffix

        repToIndex.append(repName)

        parent = ""
        for part in rep.split("/")[:-1]:
            parent += part + "/"

        printlog("parent = " + parent)
        sys.path.insert(0, parent)

        os.system("mkdir " + repName)
        sys.path.insert(0, rep)

        os.system("sphinx-apidoc -o " + repName + "/ " + rep + " -f")


def IndexMake():
    indexName = "index.rst"
    f = open(indexName, "r")
    try:
        if True:
                Contenu = f.read().splitlines()
                result = ""
                can = True
               # result = Contenu.replace("from . ","")
                for line in Contenu:
                    if can is False:
                        can = True if (line.find("Indices and tables") != -1) \
                            else False

                    if line.find(":maxdepth") != -1:
                        can = False
                        result += line + "\n \n"
                        identattion = line.split(":maxdepth")[0]

                        for name in repToIndex:
                            result += identattion + name + "/modules" + "\n"

                        result += "\n"

                    elif can is True:
                        result += line + "\n"

                f.close()
                result = result[:-1]
                RewriteFile(indexName, result)

    except IOError:
        printlog(FileName + " Doesn't exist!")


def ChecklLastError(LastError, NewError):
    Continue = True
    printlog(NewError)

    if(LastError[0] == NewError):
        printlog("LastError : " + LastError[0] + " et error = " + NewError)

        printlog("Drop on " + NewError)
        compteur = 10
        Continue = False

    LastError[0] = NewError

    return Continue


def ResolveError():
    for rep in RepertoireToSearchOn:
        printlog("Searching on "+rep)
        for root, dirs, files in os.walk(rep):
            for name in files:
                if name[-3:] == ".py":
                    printlog("Checking : " + name)
                    LastError = [""]
                    compteur = 0
                    NameModuleToFind = ""

                    while compteur < 10:
                        try:
                            compteur += 1
                            __import__(name)

                        except ImportError as e:
                            if ChecklLastError(LastError, e.message) is False:
                                break

                            Lines = traceback.format_exc().splitlines()

                            try:
                                tist = Lines[-2].split(' ')[5]
                            except IndexError as e:
                                printlog(e.message + " " + str(Lines))
                                break
                            #printlog (str(Lines[-2].split(' ')))

                            AllPart = tist.split(".")
                            if len(AllPart) > 1:
                                sys.modules[tist] = mock.MagicMock()
                                printlog("mocked = " + tist)
                                #from apscheduler.jobstores.mongodb_store import MongoDBJobStore
                                #ImportError: No module named jobstores.mongodb_store
                                moduleMissing = Lines[-2].split(' ')[-1]
                                if moduleMissing != "*":
                                    tempPath = Lines[-3].split('"')
                                    if len(tempPath) > 1:
                                        fullPath = tempPath[1]
                                        moduleMissing = Lines[-2].split(' ')[-1]
                                        moduleMissing = moduleMissing.replace("'","")
                                        ReplaceImport(os.path.join(fullPath), moduleMissing, Lines[-2])

                            for part in AllPart:
                                #printlog ("to import = "+part)
                                if part != ("(name)"):
                                        for subpart in part.split(","):
                                            printlog("add mock for " + subpart)
                                            ListAdded.append(subpart)
                                            sys.modules[subpart] = mock.MagicMock()

                        except (OSError, IOError) as e:
                            Lines = traceback.format_exc().splitlines()
                            INfoPath = Lines[-1]
                            printlog("OsError : " + INfoPath)
                            if(ChecklLastError(LastError, INfoPath) is False):
                                break

                            if INfoPath.find("[Errno 2]"):
                                printlog(traceback.format_exc())

                                WrongPath = Lines[-1].split("'")[-2]
                                printlog("WrongPath = " + WrongPath)
                                DirsList = WrongPath.split("/")
                                Last = DirsList[-1]
                                DirsList = DirsList[1:-1]
                                printlog(str(DirsList))
                                path = "/"

                                for dirToCreate in DirsList:
                                    path += dirToCreate + "/"
                                    if dirToCreate != "":
                                        printlog("creating : " + path)
                                        os.system("mkdir " + path)

                                path += Last

                                if(Last.find(".") != -1):
                                    os.system("touch " + path)

                                else:
                                    os.system("mkdir " + path)

                        except NameError as e:
                            #  NameError: name 'relativedelta' is not defined
                            Lines = traceback.format_exc().splitlines()
                            if(ChecklLastError(LastError, e.message) is False):
                                break

                            tempPath = Lines[-3].split('"')
                            if len(tempPath) > 1:
                                fullPath = tempPath[1]
                                moduleMissing = Lines[-1].split(' ')[-4]
                                moduleMissing = moduleMissing.replace("'", "")
                                AddImport(
                                    os.path.join(fullPath), moduleMissing)

                        except Exception:
                            printlog("Error commune")
                            Lines = traceback.format_exc().splitlines()
                            if(ChecklLastError(LastError, Lines[-1]) is False):
                                break

    printlog("added : ")
    for added in ListAdded:
        printlog(added)


def AutoDoc2():
    #AddAll("../sources/")
    AutoApi()
    IndexMake()
    ResolveError()

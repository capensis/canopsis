#!/usr/bin/env python

import sys, os, fnmatch
import shutil

## Configurations

install_path = "/opt/canopsis/"
webui_path = os.path.join(install_path, "var/www/canopsis")

cps_filename = "canopsis.js"
cps_min_filename = "canopsis.min.js"
cps_gz_filename = "canopsis.min.js.gz"

cps_filepath = os.path.join(webui_path, cps_filename)
cps_min_filepath = os.path.join(webui_path, cps_min_filename)
cps_gz_filepath = os.path.join(webui_path, cps_gz_filename)

debug = False

exclude_files = [
]

paths = [
	os.path.join(webui_path, "app/lib/locale.js"),
	os.path.join(webui_path, "app/lib/global.js"),
	os.path.join(webui_path, "app/lib/locale.js"),
	os.path.join(webui_path, "app/lib/log.js"),
	os.path.join(webui_path, "app/lib/tools.js"),
	os.path.join(webui_path, "auth.js"),
	os.path.join(webui_path, "app.js"),
	os.path.join(webui_path, "app/lib/renderers.js"),
	os.path.join(webui_path, "app/lib/form/cfield.js"),
	os.path.join(webui_path, "app/lib/form/field/cdate.js"),
	os.path.join(webui_path, "app/lib/view/cgrid.js"),
	os.path.join(webui_path, "app/lib/view/cgrid_state.js"),
	os.path.join(webui_path, "app/view/Tabs/Content.js"),
	os.path.join(webui_path, "app/model"),
	os.path.join(webui_path, "app/lib/store"),
	os.path.join(webui_path, "app/store"),
	os.path.join(webui_path, "app/lib/menu"),
	os.path.join(webui_path, "app/lib/form"),
	os.path.join(webui_path, "app/lib/view"),
	os.path.join(webui_path, "app/view/Tabs/Content.js"),
	os.path.join(webui_path, "app/view"),
	os.path.join(webui_path, "app/lib"),
	os.path.join(webui_path, "app/lib/controller"),
	os.path.join(webui_path, "app/controller"),
	os.path.join(webui_path, "widgets/line_graph/line_graph.js"),
	os.path.join(webui_path, "widgets/"),
	os.path.join(webui_path, "app/controller/Widgets.js"),
	os.path.join(webui_path, "app/view/Viewport.js")
]

appended_files = []

## Functions

def locate(pattern, root=os.curdir):
    '''Locate all files matching supplied filename pattern in and below
    supplied root directory.'''
    for path, dirs, files in os.walk(os.path.abspath(root)):
        for filename in fnmatch.filter(files, pattern):
            yield os.path.join(path, filename)

def append_file(file_path):
	if debug:
		min_file.write("\n/* %s */\n" % file_path)
		min_file.write("console.log('   -> %s')\n" % file_path)
	shutil.copyfileobj(open(file_path, 'r'), min_file)
	appended_files.append(file_path)

def concact_files(wpath):
	if os.path.isfile(wpath):
		file_path = wpath
		if file_path not in exclude_files and file_path not in appended_files:
			print "  + %s" % file_path
			append_file(file_path)
	else:
		for file_path in locate("*.js", wpath):
			if file_path not in exclude_files and file_path not in appended_files:
				print "  + %s" % file_path
				append_file(file_path)

## Main

print "Remove old files"
for path in [cps_filepath, cps_min_filepath, cps_gz_filepath]:
	if os.path.exists(path):
		print " + %s" % path
		os.remove(path)


print "Open '%s' in write mode" % cps_filename
min_file = open(cps_filepath, "w")

if debug:
	min_file.write("console.log('Start canopsis.min.js')\n")

print " + Append files"
if debug:
	min_file.write("console.log(' + Load normal files')\n")
for path in paths:
	concact_files(path)

print "Close '%s'" % cps_filepath
if debug:
	min_file.write("console.log('End canopsis.min.js')\n")

print "%s files appended in '%s'" % (len(appended_files), cps_filename)

min_file.close()


print "Minimify '%s'" % cps_filename
if not os.path.exists('sources/externals/rjsmin-1.0.5/rjsmin.py'):
	if not os.path.exists('sources/externals/rjsmin-1.0.5.tar.gz'):
		print " + Error: 'rjsmin-1.0.5.tar.gz' not found in external dir"
		sys.exit(1)
	else:
		print "Extract 'rjsmin-1.0.5.tar.gz' ..."
		os.system('cd sources/externals && tar xfz rjsmin-1.0.5.tar.gz')
		if not os.path.exists('sources/externals/rjsmin-1.0.5/rjsmin.py'):
			print " + Error: Impossible to extract"
			sys.exit(1)

		print " + Done"

os.system("cat %s | sources/externals/rjsmin-1.0.5/rjsmin.py > %s" % (cps_filepath, cps_min_filepath))
if not os.path.exists(cps_min_filepath):
	print " + Error: Impossible to minimify"
	sys.exit(1)

print " + Done"

print "Compress '%s'" % cps_min_filename
os.system("gzip -c -9 %s > %s" % (cps_min_filepath, cps_gz_filepath))
if not os.path.exists(cps_gz_filepath):
	print " + Error: Impossible to compress"
	sys.exit(1)

print " + Done"
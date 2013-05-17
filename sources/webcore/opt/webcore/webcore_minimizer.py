#!/usr/bin/env python

import sys, os, fnmatch
import shutil
import getpass

## Configurations

if getpass.getuser() != 'canopsis':
	install_path = "/opt/canopsis/"
else:
	install_path = os.path.expanduser('~')
	
webui_path = os.path.join(install_path, "var/www/canopsis")

cps_filename = "canopsis.js"
cps_min_filename = "canopsis.min.js"
cps_gz_filename = "canopsis.min.js.gz"

cps_filepath = os.path.join(webui_path, cps_filename)
cps_min_filepath = os.path.join(webui_path, cps_min_filename)
cps_gz_filepath = os.path.join(webui_path, cps_gz_filename)

debug = False

exclude_files = [
	os.path.join(webui_path, "app/lib/locale.js"),
	os.path.join(webui_path, "app/lib/global.js"),
	os.path.join(webui_path, "app/lib/global_options.js")
]

paths = [
	os.path.join(webui_path, "app/lib/log.js"),
	os.path.join(webui_path, "app/lib/tools.js"),
	os.path.join(webui_path, "auth.js"),
	os.path.join(webui_path, "app.js"),
	os.path.join(webui_path, "app/lib/renderers.js"),
	os.path.join(webui_path, "app/lib/form/cfield.js"),
	os.path.join(webui_path, "app/lib/form/field/cdate.js"),
	os.path.join(webui_path, "app/lib/form/field/cdatePicker.js"),
	os.path.join(webui_path, "app/lib/view/ccard.js"),
	os.path.join(webui_path, "app/lib/view/cgrid.js"),
	os.path.join(webui_path, "app/lib/view/cgrid_state.js"),
	os.path.join(webui_path, "app/view/Tabs/Content.js"),
	os.path.join(webui_path, "app/model"),
	os.path.join(webui_path, "app/lib/store"),
	os.path.join(webui_path, "app/store"),
	os.path.join(webui_path, "app/lib/menu"),
	os.path.join(webui_path, "app/lib/form"),
	os.path.join(webui_path, "app/lib/view/cpopup.js"),
	os.path.join(webui_path, "app/lib/view"),
	os.path.join(webui_path, "app/view/Tabs/Content.js"),
	os.path.join(webui_path, "app/view"),
	os.path.join(webui_path, "app/lib"),
	os.path.join(webui_path, "app/lib/controller"),
	os.path.join(webui_path, "app/controller"),
	os.path.join(webui_path, "widgets/line_graph/line_graph.js"),
	os.path.join(webui_path, "widgets/"),
	os.path.join(webui_path, "../widgets/"),
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

def append_file(file, file_path):
	file.write("\n/* %s */\n" % file_path)
	if debug:
		file.write("console.log('   -> %s')\n" % file_path)
	shutil.copyfileobj(open(file_path, 'r'), file)
	appended_files.append(file_path)

def exclude_locales(wpath):
	for file_path in locate("*lang*.js", wpath):
		exclude_files.append(file_path)

def concact_files(file, wpath):
	if os.path.isfile(wpath):
		file_path = wpath
		if file_path not in exclude_files and file_path not in appended_files:
			print "  + %s" % file_path
			append_file(file, file_path)
	else:
		for file_path in locate("*.js", wpath):
			if file_path not in exclude_files and file_path not in appended_files:
				print "  + %s" % file_path
				append_file(file, file_path)

## Main

print "Remove old files"
for path in [cps_filepath, cps_min_filepath, cps_gz_filepath]:
	if os.path.exists(path):
		print " + %s" % path
		os.remove(path)

print " + Exclude locales file"
exclude_locales(webui_path)

print "Open '%s' in write mode" % cps_filename
file = open(cps_filepath, "w")

if debug:
	file.write("console.log('Start canopsis.min.js')\n")

print " + Append files"
if debug:
	file.write("console.log(' + Load normal files')\n")
for path in paths:
	concact_files(file, path)

print "Close '%s'" % cps_filepath
if debug:
	file.write("console.log('End canopsis.min.js')\n")

file.close()
print "%s files appended in '%s'" % (len(appended_files), cps_filename)

print "Minimify '%s' to '%s'" % (cps_filename, cps_min_filename)
if not os.path.exists('%s/bin/uglifyjs' % install_path):
	print " + Error: 'uglifyjs'"
	sys.exit(1)

os.system("cd %s && bin/node bin/uglifyjs %s > %s" % (install_path, cps_filepath, cps_min_filepath))

print " + Done"

print "Compress '%s' to '%s'" % (cps_min_filename, cps_gz_filename)
os.system("gzip -c -9 %s > %s" % (cps_min_filepath, cps_gz_filepath))
if not os.path.exists(cps_gz_filepath):
	print " + Error: Impossible to compress"
	sys.exit(1)

print " + Done"

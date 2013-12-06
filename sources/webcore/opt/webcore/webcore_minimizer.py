#!/usr/bin/env python2.7
# -*- coding: utf-8 -*-

import subprocess
import fnmatch
import gzip
import csv
import os


class Minimizer(object):
	"""
		Minimize JavaScript files.
	"""

	def __init__(self, *args, **kwargs):
		super(Minimizer, self).__init__(*args, **kwargs)

		self.rootdir = '/opt/canopsis'
		self.canodir = os.path.join(self.rootdir, 'var', 'www', 'canopsis')

		self.canopsisjs      = os.path.join(self.canodir, 'canopsis.js')
		self.canopsisminjs   = os.path.join(self.canodir, 'canopsis.min.js')
		self.canopsisminjsgz = os.path.join(self.canodir, 'canopsis.min.js.gz')

		self.nodejspath = os.path.join(self.rootdir, 'bin', 'node')
		self.uglifypath = os.path.join(self.rootdir, 'bin', 'uglifyjs')

		self.exclude = [
			os.path.join(self.canodir, 'resources'),
			os.path.join(self.canodir, 'themes'),

			os.path.join(self.canodir, 'canopsis.js'),
			os.path.join(self.canodir, 'canopsis.min.js'),

			os.path.join(self.canodir, 'app', 'lib', 'locale.js'),
			os.path.join(self.canodir, 'app', 'lib', 'global.js'),
			os.path.join(self.canodir, 'app', 'lib', 'global_options.js')
		]

		self.files = []
		self.added = []

	def jsFinder(self, path, dirname, filenames):
		if dirname in self.exclude:
			return

		for filename in fnmatch.filter(filenames, '*.js'):
			fullpath = os.path.join(dirname, filename)

			if fullpath not in self.exclude and not fnmatch.fnmatch(filename, '*lang*.js'):
				self.files.append(fullpath)

	def get_deps(self, filename):
		files = []

		with open(filename, 'r') as f:
			first_line = f.readline()

			# First line may contains dependencies in CSV formatted list
			if first_line.startswith('//need:'):
				requirements = first_line[7:]
				parsed = csv.reader([requirements]).next()

				# Add all required files if not excluded or already added
				for required in parsed:
					fullpath = os.path.join(self.canodir, required)

					if fullpath not in self.exclude and fullpath not in files and fullpath not in self.added:
						files += self.get_deps(fullpath)

			# Now dependencies are in the list, just add this file now
			files.append(filename)

		# Be sure every file is unique
		found = []

		for f in files:
			if f not in found:
				found.append(f)

		return found

	def minify(self):
		# Find all JS files
		os.path.walk(self.canodir, self.jsFinder, None)

		# Order files : lib, model, store, view, widgets, controller, general
		
		orderers = [
			{
				'path': os.path.join(self.canodir, 'app', 'lib'),
				'files': []
			},
			{
				'path': os.path.join(self.canodir, 'app', 'model'),
				'files': []
			},
			{
				'path': os.path.join(self.canodir, 'app', 'store'),
				'files': []
			},
			{
				'path': os.path.join(self.canodir, 'app', 'view'),
				'files': []
			},
			{
				'path': os.path.join(self.canodir, 'widgets'),
				'files': []
			},
			{
				'path': os.path.join(self.canodir, 'app', 'controller'),
				'files': []
			},
			{
				'path': self.canodir,
				'files': []
			},
		]

		for filename in self.files:
			for orderer in orderers:
				if filename.startswith(orderer['path']):
					orderer['files'].append(filename)
					break

		self.files = []

		for orderer in orderers:
			self.files += orderer['files']

		# Handle dependencies

		for filename in self.files:
			if filename in self.exclude or filename in self.added:
				continue

			# Open file and read the first line
			self.added += self.get_deps(filename)

		# Concatenate all files
		print 'Open {0} for writing'.format(self.canopsisjs)

		with open(self.canopsisjs, 'w') as canopsisjs:
			for filename in self.added:
				canopsisjs.write('\n/* file:{0} */\n'.format(filename))

				print '- ', filename

				with open(filename, 'r') as f:
					line = f.readline()

					while line:
						canopsisjs.write(line)

						line = f.readline()

		print '{0} written'.format(self.canopsisjs)

		# Minify canopsis.js
		print 'Open {0} for writting'.format(self.canopsisminjs)

		with open(self.canopsisminjs, 'w') as canopsisminjs:
			subprocess.call(
				[self.nodejspath, self.uglifypath, self.canopsisjs],
				stdout=canopsisminjs
			)

		print '{0} written'.format(self.canopsisminjs)

		# Compress in GZ format
		print 'Open {0} for writting'.format(self.canopsisminjsgz)

		with open(self.canopsisminjs, 'rb') as canopsisminjs:
			with gzip.open(self.canopsisminjsgz, 'wb') as canopsisminjsgz:
				canopsisminjsgz.writelines(canopsisminjs)

		print '{0} written'.format(self.canopsisminjsgz)


if __name__ == "__main__":
	minimizer = Minimizer()
	minimizer.minify()
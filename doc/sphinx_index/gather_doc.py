# -*- coding: utf-8 -*-

import os
import glob
import subprocess
import shutil

# Sphinx index contains some conf and index pages that do not
# directly depend on a version canopsis. Those piece of documentation
# and files that are required for sphinx-build process
SPHINX_INDEX = {'repo': 'canopsis', 'branch': 'NRPUIV2', 'folder': 'sphinx_index'}

# What you want to include or not
VERSIONS = [{'repo': 'canopsis', 'branch': 'release-ficus', 'folder': 'Ficus'},
            {'repo': 'canopsis', 'branch': 'NRPUIV2', 'folder': 'Sakura'},
           ]

# Same thing about connectors
CONNECTORS = [{'repo': 'mod-canopsis', 'branch': 'master', 'name': 'shinken'},
             ]

def cmd(*args):
  """Execute a shell command"""
  p = subprocess.Popen(args, stdout=subprocess.PIPE)
  return p.communicate()[0]

def mvrm(source, destination):
  """Recursively mv files from source to destination and rm source"""
  for path in glob.glob(source + "/*"):
    shutil.move(path, destination)
  os.rmdir(source)

if os.path.exists('doc'):
  print 'A folder named \'doc\' already exists : aborting procedure'
  quit()
else:
  print 'Creating folder \'doc\''
  os.makedirs('doc')
  os.chdir('doc')

# SPHINX_INDEX
print 'Downloading sphinx index from branch ' + SPHINX_INDEX['branch'] + ' (' + SPHINX_INDEX['repo'] + ') in \'doc\''
print cmd('svn', 'export', 'https://github.com/capensis/' + SPHINX_INDEX['repo']  + '/branches/' + SPHINX_INDEX['branch'] + '/doc/' + SPHINX_INDEX['folder'])
mvrm(SPHINX_INDEX['folder'], '.')

# VERSIONS
with open('index.rst', 'r+') as root_index, open('new_index.rst', 'w') as new_index: # We need to add each version in the main index
  for line in root_index:
    new_index.write(line)
    if 'overview' in line: # This is where titles need to be added
      for version in VERSIONS: # Download each version and add it in the index
        print 'Downloading documentation from branch ' + version['branch'] + ' (' + version['repo'] + ') in \'doc/' + version['folder'] + '\''
        print cmd('svn', 'export', 'https://github.com/capensis/' + version['repo'] + '/branches/' + version['branch']  + '/doc/' + version['folder'])
        new_index.write('   ' + version['folder'] + '/index\n')
      new_index.write('   connectors/index\n')
os.rename('new_index.rst', 'index.rst')

# CONNECTORS
os.chdir('connectors')
with open('index.rst', 'a') as c_index:
  for connector in CONNECTORS:
    print 'Downloading documentation for connector in repo ' + connector['repo'] + ' (' + connector['branch'] + ') in \'doc/connectors\''
    print cmd('svn', 'export', 'https://github.com/capensis/' + connector['repo'] + '/branches/' + connector['branch']  + '/doc/')
    for rst in glob.glob('doc/*.rst'):
      file_name = rst[4:-4]
      capitalized_file_name = file_name.capitalize()
      index_name = capitalized_file_name.replace('-', ' ').replace('_', ' ')
      c_index.write('   ' + index_name + '\n')
      shutil.move(rst, '.')
    if os.path.exists('doc/img'):
      mvrm('doc/img', '../_static/images/connectors')
    shutil.rmtree('doc')

os.chdir('..')
print 'Building sphinx documentation'
print cmd('sphinx-build', '-b', 'html', '.', '_build')
print 'Process completed'
print 'Documentation is available in \'doc/_build\''

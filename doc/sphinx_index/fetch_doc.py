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
VERSIONS = [{'repo': 'canopsis', 'branch': 'release-ficus', 'folder': 'ficus'},
            {'repo': 'canopsis', 'branch': 'NRPUIV2', 'folder': 'sakura'},
           ]

# Same thing about connectors
CONNECTORS = [{'repo': 'mod-canopsis', 'branch': 'master', 'name': 'shinken'},
             ]

# html ligne to insert titles in the main index
HTML_INDEX_TITLE = '<li class="toctree-l1"><a class="reference internal" href="{0}/index.html">{1}</a></li>\n'

def cmd(*args):
  """Execute a shell command"""
  p = subprocess.Popen(args, stdout=subprocess.PIPE)
  return p.communicate()[0]

def mvrm(source, destination):
  """Recursively mv files from source to destination and rm source"""
  for path in glob.glob(source + '/*'):
    shutil.move(path, destination)
  shutil.rmtree(source)

def mvrmbuild(source):
  """Similar to mvrm to extract the built doc to the correct location"""
  for element in glob.glob(source + '/*'):
    if element != source + '/_build':
      if os.path.isdir(element):
        shutil.rmtree(element)
      else:
        os.remove(element)
  mvrm(source + '/_build', source)


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
os.remove('fetch_doc.py')

# VERSIONS
# We open index and a new_index. index content will be written in new_index, so as we can write title lines
# for each version where we want. At the end of this bloc we ecrase index by new_index
with open('index.html', 'r+') as root_index, open('new_index.html', 'w') as new_index:
  for line in root_index:
    new_index.write(line)
    if 'toctree-l1 current' in line: # This is where titles need to be added
      for version in VERSIONS: # Download each version and build it
        print 'Downloading documentation from branch ' + version['branch'] + ' (' + version['repo'] + ') in \'doc/' + version['folder'] + '\''
        print cmd('svn', 'export', 'https://github.com/capensis/' + version['repo'] + '/branches/' + version['branch']  + '/doc/' + version['folder'])
        print cmd('sphinx-build', '-b', 'html', version['folder'], version['folder'] + '/_build')
        new_index.write(HTML_INDEX_TITLE.format(version['folder'], version['folder'].capitalize()))
        mvrmbuild(version['folder']) # Just keep the html files in _build

      if CONNECTORS:
        new_index.write(HTML_INDEX_TITLE.format('connectors', 'Connectors'))
os.rename('new_index.html', 'index.html')

## CONNECTORS
if not CONNECTORS:
  shutil.rmtree('connectors')
else:
  os.chdir('connectors')
  with open('index.rst', 'a') as c_index:
    for connector in CONNECTORS:
      print 'Downloading documentation for connector in repo ' + connector['repo'] + ' (' + connector['branch'] + ') in \'doc/connectors\''
      print cmd('svn', 'export', 'https://github.com/capensis/' + connector['repo'] + '/branches/' + connector['branch']  + '/doc/')
      for rst in glob.glob('doc/*.rst'):
        index_entry = rst[4:-4] # doc/my-connec_tor.rst --> my-connec_tor
        c_index.write('   ' + index_entry + '\n') # Indexes in the connector index with rst syntax
        shutil.move(rst, '.')
      if os.path.exists('doc/img'):
        mvrm('doc/img', '_static/images')
      shutil.rmtree('doc') # cleaning
  print cmd('sphinx-build', '-b', 'html', '.', '_build')
  mvrmbuild('.')

print 'Process completed'
print 'Documentation is available in \'doc\''

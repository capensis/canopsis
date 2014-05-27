import os


def tasks(path):
	path = os.path.expanduser(path)
	dirList = os.listdir(path)
	list_tasks = []
	for mfile in dirList:
		ext = mfile.split(".")[1]
		name = mfile.split(".")[0]
		if name != "." and ext == "py" and name != '__init__':
			list_tasks.append(name)
	return tuple(list_tasks)

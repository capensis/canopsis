import string
import socket

if __name__ == '__main__':
	import sys

	iface = '192.168.3.56'
	port = 162
	community = 'public'
	
	# Create trap server
	sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	sock.bind((iface, port))

	print 'Listening for SNMP traps on ' + str((iface, port))

	# Receive SNMP traps, process them, print details to stdout
	while 1:
		(message, peer) = sock.recvfrom(65536)

		print message

#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from socket import socket, AF_INET, SOCK_DGRAM

if __name__ == '__main__':

    iface = '192.168.3.56'
    port = 162
    community = 'public'

    # Create trap server
    sock = socket(AF_INET, SOCK_DGRAM)
    sock.bind((iface, port))

    print('Listening for SNMP traps on ' + str((iface, port)))

    # Receive SNMP traps, process them, print details to stdout
    while 1:
        (message, peer) = sock.recvfrom(65536)

        print(message)

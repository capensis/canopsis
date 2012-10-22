notifications_oid = {
 '1.3.6.1.4.1.8072.4.0.1': 'nsNotifyStart',
 '1.3.6.1.4.1.8072.4.0.2': 'nsNotifyShutdown',
 '1.3.6.1.4.1.8072.4.0.3': 'nsNotifyRestart'
}

notifications = {
	'nsNotifyShutdown': {
                             'SEVERITY': 'CRITICAL',
			     'STATE': 'DEGRADED',
                             'TYPE': 'SNMP Agent Shutdown'},
	'nsNotifyStart': {
                             'SEVERITY': 'INFORMATIONAL',
                             'STATE': 'OPERATIONAL',
                             'TYPE': 'SNMP Agent start'},
	'nsNotifyRestart': {
                             'SEVERITY': 'INFORMATIONAL',
                             'STATE': 'OPERATIONAL',
                             'TYPE': 'SNMP Agent restart'},

}

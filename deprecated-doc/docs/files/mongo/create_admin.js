// Add super user
conn = new Mongo();
db = conn.getDB('admin');

if(db.getUser('admin') == null) {
    db.createUser(
        {
            'user': 'admin',
            'pwd': 'admin',
            'roles': [
                {
                    'role': 'readWriteAnyDatabase',
                    'db': 'admin'
                },
                {
                    'role': 'userAdminAnyDatabase',
                    'db': 'admin'
                },
                {
                    'role': 'dbAdminAnyDatabase',
                    'db': 'admin'
                },
                {
                    'role': 'root',
                    'db': 'admin'
                }
            ]
        }
    );
}
else {
    db.updateUser(
        'admin',
        {
            'pwd': 'admin',
            'roles': [
                {
                    'role': 'readWriteAnyDatabase',
                    'db': 'admin'
                },
                {
                    'role': 'userAdminAnyDatabase',
                    'db': 'admin'
                },
                {
                    'role': 'dbAdminAnyDatabase',
                    'db': 'admin'
                },
                {
                    'role': 'root',
                    'db': 'admin'
                }
            ]
        }
    );
}

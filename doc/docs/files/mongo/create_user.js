// Add common user
conn = new Mongo();
db = conn.getDB('admin');
db.auth('admin', 'admin');

db = conn.getDB('canopsis');

if(db.getUser('cpsmongo') == null) {
    db.createUser(
        {
            'user': 'cpsmongo',
            'pwd': 'canopsis',
            'roles': [
                {
                    'role': 'dbOwner',
                    'db': 'canopsis'
                }
            ]
        }
    );
}
else {
    db.updateUser(
        'cpsmongo',
        {
            'pwd': 'canopsis',
            'roles': [
                {
                    'role': 'dbOwner',
                    'db': 'canopsis'
                }
            ]
        }
    );
}

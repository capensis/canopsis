// Add common user
conn = new Mongo();
db = conn.getDB('admin');
db.auth('root', 'root');

canopsis_database = 'canopsis'
canopsis_user = 'cpsmongo'
canopsis_pwd = 'canopsis'

db = conn.getDB('canopsis');

if(db.getUser(canopsis_user) == null) {
    db.createUser(
        {
            'user': canopsis_user,
            'pwd': canopsis_pwd,
            'roles': [
                {
                    'role': 'dbOwner',
                    'db': canopsis_database
                }
            ]
        }
    );
}
else {
    db.updateUser(
        canopsis_user,
        {
            'pwd': canopsis_pwd,
            'roles': [
                {
                    'role': 'dbOwner',
                    'db': canopsis_database
                }
            ]
        }
    );
}

// Initiate replicaset and wait for primary election
//
// Enabling replicaset mode is a requisite for Canopsis engines to be able to
// use transactions in MongoDB.
//
// Yes, it is a single-node replicaset, which sounds funny but it is sufficient
// for a lab environment like here.
//
// Read: https://www.mongodb.com/community/forums/t/why-replica-set-is-mandatory-for-transactions-in-mongodb/9533/2

try {
    rs.initiate({ _id : "rs0", members: [{ _id: 0, host: "mongodb:27017" }]});
} catch (err) {
    if (err.code == 23) {  // 23: AlreadyInitialized
        print('Replicaset was already configured: OK');
    } else {
        throw err;
    }
}

let primaryWaitRetries = 0;
let isMaster = db.adminCommand({isMaster: 1}).ismaster;
while (!isMaster && primaryWaitRetries < 60) {
    print('Node is not primary yet.');
    primaryWaitRetries++;
    sleep(500); // 500 ms
    isMaster = db.adminCommand({isMaster: 1}).ismaster;
}
if (isMaster) {
    print('Node is primary. Init job complete!');
} else {
    print('Timeout while waiting for primary election (~ 30 s). Aborting.');
    /* Passing an arg to MongoDB's quit() function seems to allow for exit
     * code management, as intended here.
     * However this is absolutely undocumented by:
     * https://www.mongodb.com/docs/v5.0/reference/method/quit/
     * Let's hope they don't remove this secret feature!
     */
    quit(1);
}

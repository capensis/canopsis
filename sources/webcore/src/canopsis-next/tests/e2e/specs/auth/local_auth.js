module.exports = {
  'Canopsis Authentication' : function (client) {
    client
      .canopsis_auth(client, client.globals.canopsisEnv.url, client.globals.canopsisEnv.username, client.globals.canopsisEnv.password)
      .end()
  }
};


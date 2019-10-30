function getBaseUrl(path = '') {
  return (process.env.VUE_DEV_SERVER_URL + path).replace(/([^:]\/)\/+/g, '$1');
}

module.exports.getBaseUrl = getBaseUrl;

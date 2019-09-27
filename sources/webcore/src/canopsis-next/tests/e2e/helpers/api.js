const qs = require('qs');
const axios = require('axios');
const faker = require('faker');

const { API_ROUTES, DEFAULT_LOCALE } = require('@/config');
const { generateUser, prepareUserByData } = require('@/helpers/entities');

const { CREDENTIALS } = require('../constants');

const { getBaseUrl } = require('./url');

const request = axios.create({
  baseURL: getBaseUrl('/api'),
  withCredentials: true,
});

function auth(username, password) {
  return request.post(API_ROUTES.auth, {
    username,
    password,

    json_response: true,
  });
}

function authAsAdmin() {
  return auth(CREDENTIALS.admin.username, CREDENTIALS.admin.password);
}

async function createUser(user) {
  const { headers } = await authAsAdmin();

  const userData = prepareUserByData(user);

  const result = await request.post(API_ROUTES.user.create, qs.stringify({ user: JSON.stringify(userData) }), {
    headers: {
      'content-type': 'application/x-www-form-urlencoded',
      Cookie: headers['set-cookie'],
    },
  });

  result.data._id = user._id;
  result.data.password = user.password;

  return result;
}

function createAdminUser(defaultUser) {
  const user = { ...generateUser(), ...defaultUser };

  if (!defaultUser) {
    user._id = faker.internet.userName();
    user.password = faker.internet.password();
    user.mail = faker.internet.email();
    user.firstname = faker.name.firstName();
    user.lastname = faker.name.lastName();
    user.ui_language = DEFAULT_LOCALE;
    user.role = 'admin';
  }

  return createUser(user);
}

async function removeUser(id) {
  const { headers } = await authAsAdmin();

  return request.delete(`${API_ROUTES.user.remove}/${id}`, {
    headers: {
      Cookie: headers['set-cookie'],
    },
  });
}

module.exports.auth = auth;
module.exports.authAsAdmin = authAsAdmin;
module.exports.createUser = createUser;
module.exports.createAdminUser = createAdminUser;
module.exports.removeUser = removeUser;

const axios = require('axios');
const qs = require('qs');
const sha1 = require('sha1');
const faker = require('faker');

const { API_ROUTES } = require('@/config');
const { generateUser } = require('@/helpers/entities');

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

  const userData = { ...user };

  if (userData.password && userData.password !== '') {
    userData.shadowpasswd = sha1(userData.password);
  }

  const result = await request.post(API_ROUTES.user.create, qs.stringify({ user: JSON.stringify(userData) }), {
    headers: {
      'content-type': 'application/x-www-form-urlencoded',
      Cookie: headers['set-cookie'],
    },
  });

  result.data._id = userData._id;

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
    user.role = 'admin';
  }

  return createUser(user);
}

async function removeUser(id) {
  const { headers } = await authAsAdmin();

  return request.delete(`${API_ROUTES.user.remove}/${id}`, {
    Cookie: headers['set-cookie'],
  });
}

module.exports.auth = auth;
module.exports.authAsAdmin = authAsAdmin;
module.exports.createUser = createUser;
module.exports.createAdminUser = createAdminUser;
module.exports.removeUser = removeUser;

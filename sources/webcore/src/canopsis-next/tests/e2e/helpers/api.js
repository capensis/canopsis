const qs = require('qs');
const axios = require('axios');
const faker = require('faker');
const { omit } = require('lodash');

const { API_ROUTES, DEFAULT_LOCALE } = require('@/config');
const {
  generateViewTab,
  generateUser,
  prepareUserByData,
  generateRight,
  generateRoleRightByChecksum,
} = require('@/helpers/entities');
const { generateTemporaryView } = require('./entities');

const { CREDENTIALS } = require('../constants');
const { USERS_RIGHTS_MASKS, USERS_RIGHTS_TYPES } = require('@/constants');

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

async function createView(viewDefault, options = {}) {
  const response = await request.post(API_ROUTES.view, {
    tabs: [generateViewTab()],
    ...viewDefault,
  }, options);

  return response.data;
}

async function createViewGroup(group, options = {}) {
  const response = await request.post(API_ROUTES.viewGroup, group, options);

  return response.data;
}

async function removeViewGroup(groupId, options = {}) {
  return request.post(`${API_ROUTES.viewGroup}/${groupId}`, options);
}

async function getRoleList(role = 'admin', options = {}) {
  const response = await request.get(`${API_ROUTES.role.list}/${role}`, options);

  return response.data.data;
}

async function addRightsForView(viewId, options = {}) {
  const right = {
    ...generateRight(),
    _id: viewId,
    type: USERS_RIGHTS_TYPES.rw,
    desc: `Rights on view: ${viewId}`,
  };

  const response = await request.post(API_ROUTES.action, right, options);

  return response.data.data;
}

function removeRightsForView(viewId, options = {}) {
  return request.post(`${API_ROUTES.action}/${viewId}`, options);
}

async function createUserRole(role, rights, options = {}) {
  return request.post(API_ROUTES.role.create, {
    role: {
      ...role,
      rights,
    },
  }, options);
}

async function createWidgetView(viewDefault) {
  const generatedView = generateTemporaryView();

  const response = await authAsAdmin();

  const options = {
    headers: {
      Cookie: response.headers['set-cookie'],
    },
  };

  const { _id: groupId } = await createViewGroup({
    name: generatedView.group,
  }, options);

  const { _id: viewId } = await createView({
    ...generatedView,
    ...viewDefault,
    group_id: groupId,
  }, options);

  const [right] = await addRightsForView(viewId, options);

  const [role] = await getRoleList('admin', options);

  const checksum = USERS_RIGHTS_MASKS.read + USERS_RIGHTS_MASKS.update + USERS_RIGHTS_MASKS.delete;
  await createUserRole(
    role,
    {
      ...role.rights,
      [right._id]: generateRoleRightByChecksum(checksum),
    },
    options,
  );

  return {
    viewId,
    groupId,
  };
}

async function removeWidgetView(viewId) {
  const response = await authAsAdmin();

  const options = {
    headers: {
      Cookie: response.headers['set-cookie'],
    },
  };

  await request.delete(`${API_ROUTES.view}/${viewId}`, options);

  const roles = await getRoleList('admin', options);

  await removeRightsForView(viewId, options);

  await Promise.all(roles.map(role => createUserRole(role, omit(role.rights, [viewId]), options)));
}

module.exports.auth = auth;
module.exports.authAsAdmin = authAsAdmin;
module.exports.createUser = createUser;
module.exports.createAdminUser = createAdminUser;
module.exports.removeUser = removeUser;
module.exports.createWidgetView = createWidgetView;
module.exports.removeWidgetView = removeWidgetView;
module.exports.removeViewGroup = removeViewGroup;

const qs = require('qs');
const axios = require('axios');
const faker = require('faker');
const { omit, merge } = require('lodash');

const { API_ROUTES, DEFAULT_LOCALE } = require('@/config');
const {
  generateViewTab,
  generateUser,
  prepareUserByData,
  generateRight,
  generateRoleRightByChecksum,
  generateWidgetByType,
  generateViewRow,
} = require('@/helpers/entities');
const { generateTemporaryView } = require('./entities');
const { queueFunction, onNextQueueFunction } = require('./nightwatch-child-process');

const { CREDENTIALS } = require('../constants');
const { USERS_RIGHTS_MASKS, USERS_RIGHTS_TYPES } = require('@/constants');

const { getBaseUrl } = require('./url');

const request = axios.create({
  baseURL: getBaseUrl('/backend'),
  withCredentials: true,
});

/**
 * Login through request
 * @param {string} username
 * @param {string} password
 * @returns {Promise}
 */
function auth(username, password) {
  return request.post(API_ROUTES.auth, {
    username,
    password,

    json_response: true,
  });
}

/**
 * Login as root
 * @returns {Promise}
 */
function authAsAdmin() {
  return auth(CREDENTIALS.admin.username, CREDENTIALS.admin.password);
}

/**
 * Post a user with random data
 * @param {Object} [user]
 * @returns {Object}
 */
async function createUser(user) {
  const { headers } = await authAsAdmin();

  const fakeUser = {
    ...generateUser(),
    _id: faker.internet.userName(),
    password: faker.internet.password(),
    mail: faker.internet.email(),
    firstname: faker.name.firstName(),
    lastname: faker.name.lastName(),
    ui_language: DEFAULT_LOCALE,
    ...user,
  };
  const userData = prepareUserByData(fakeUser);

  const result = await request.post(API_ROUTES.user.create, qs.stringify({ user: JSON.stringify(userData) }), {
    headers: {
      'content-type': 'application/x-www-form-urlencoded',
      Cookie: headers['set-cookie'],
    },
  });

  result.data._id = fakeUser._id;
  result.data.password = fakeUser.password;

  return result;
}

/**
 * Post a user with random data and admin role
 * @param {Object} [defaultUser]
 * @returns {Object}
 */
function createAdminUser(defaultUser) {
  return createUser({
    ...defaultUser,
    role: 'admin',
  });
}

/**
 * Delete user
 * @param {String} id
 * @returns {Promise}
 */
async function removeUser(id) {
  const { headers } = await authAsAdmin();

  return request.delete(`${API_ROUTES.user.remove}/${id}`, {
    headers: {
      Cookie: headers['set-cookie'],
    },
  });
}

/**
 * Post view for tests
 * @param {Object} viewDefault
 * @param {Object} [options]
 * @returns {Promise}
 */
async function createView(viewDefault, options = {}) {
  const response = await request.post(API_ROUTES.view, {
    tabs: [generateViewTab()],
    ...viewDefault,
  }, options);

  return response.data;
}

/**
 * Query for get view
 * @param {String} viewId
 * @param {Object} [options]
 * @returns {Promise}
 */
async function getView(viewId, options = {}) {
  const response = await request.get(`${API_ROUTES.view}/${viewId}`, options);

  return response.data;
}

/**
 * Query for update view
 * @param {String} viewId
 * @param {Object} viewDefault
 * @param {Object} [options]
 * @returns {Promise<any>}
 */
async function updateView(viewId, viewDefault, options = {}) {
  const response = await request.put(`${API_ROUTES.view}/${viewId}`, viewDefault, options);

  return response.data;
}

/**
 * Api query for create widget
 * @param {String} viewId
 * @param {Object} widget
 * @returns {Promise}
 */
async function createWidgetForView(viewId, { row, type, ...widget }) {
  const { headers } = await authAsAdmin();

  const options = {
    headers: {
      Cookie: headers['set-cookie'],
    },
  };

  const viewData = await getView(viewId, options);

  const generatedWidget = merge(generateWidgetByType(type), widget);

  const tab = {
    ...generateViewTab(),
    rows: [{
      ...generateViewRow(),
      ...row,
      widgets: [generatedWidget],
    }],
  };

  await updateView(viewId, {
    ...viewData,
    tabs: [tab],
  }, options);

  return {
    viewId,
    widgetId: generatedWidget._id,
  };
}

/**
 * Post a group for view
 * @param {Object} group
 * @param {Object} [options]
 * @returns {Object}
 */
async function createViewGroup(group, options = {}) {
  const response = await request.post(API_ROUTES.view.group, group, options);

  return response.data;
}

/**
 * Delete view group
 * @param {Number} groupId
 * @param {Object} [options]
 * @returns {Promise}
 */
async function removeViewGroup(groupId, options = {}) {
  return request.delete(`${API_ROUTES.view.group}/${groupId}`, options);
}

/**
 * Get role list for role status
 * @param {String} role
 * @param {Object} [options]
 * @returns {Promise}
 */
async function getRoleList(role = 'admin', options = {}) {
  const response = await request.get(`${API_ROUTES.role.list}/${role}`, options);

  return response.data.data;
}

/**
 * Post view rights
 * @param {Number} viewId
 * @param {Object} [options]
 * @returns {Object}
 */
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

/**
 * Delete view rights
 * @param {Number} viewId
 * @param {Object} [options]
 * @returns {Promise}
 */
function removeRightsForView(viewId, options = {}) {
  return request.delete(`${API_ROUTES.action}/${viewId}`, options);
}

/**
 * Post a user role with new rights
 * @param {String} role
 * @param {Object} rights
 * @param {Object} [options]
 * @returns {Promise}
 */
async function createUserRole(role, rights, options = {}) {
  return request.post(API_ROUTES.role.create, {
    role: {
      ...role,
      rights,
    },
  }, options);
}

/**
 * Create a view and add rights
 * @param {Object} [viewDefault]
 * @returns {Object}
 */
async function createWidgetView(viewDefault) {
  await onNextQueueFunction();

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

  await queueFunction(async () => {
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
  });

  return {
    viewId,
    groupId,
  };
}

/**
 * Remove a view and delete rights
 * @param {Number} viewId
 * @param {Number} groupId
 */
async function removeWidgetView(viewId, groupId) {
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
  await removeViewGroup(groupId, options);
}

module.exports.auth = auth;
module.exports.authAsAdmin = authAsAdmin;
module.exports.createUser = createUser;
module.exports.createAdminUser = createAdminUser;
module.exports.removeUser = removeUser;
module.exports.createWidgetView = createWidgetView;
module.exports.createWidgetForView = createWidgetForView;
module.exports.removeWidgetView = removeWidgetView;

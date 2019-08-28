const RANDOM = Math.random().toString(36).substring(6);

module.exports.CREDENTIALS = {
  admin: {
    username: 'root',
    password: 'root',
  },
};

module.exports.USERS = {
  text: 'user',
  counts: 5,
  create: {
    prefix: 'create',
  },
  edit: {
    prefix: 'edit',
  },
  mass: {
    prefix: 'mass',
  },
};

module.exports.ROLES = {
  text: 'role',
  counts: 3,
};

module.exports.NAVIGATION = {
  groups: {
    text: 'view',
    counts: 2,
    create: {
      prefix: 'create',
    },
    copy: {
      prefix: 'copy',
    },
    edit: {
      prefix: 'edit',
    },
  },
};

module.exports.USER = {
  username: 'test-username',
  firstname: 'test-firstname',
  lastname: 'test-lastname',
  email: 'test-email@example.com',
  password: 'test-password',
  role: 1,
  language: 2,
  navigationType: 1,
};

module.exports.VIEW = {
  name: `test-name-${RANDOM}`,
  title: `test-title-${RANDOM}`,
  description: `test-description-${RANDOM}`,
  group: `test-group-${RANDOM}`,
};

module.exports.SERVICE_WEATHER_WIDGET_MODAL_TYPES = {
  moreInfo: 1,
  alarmList: 2,
  both: 3,
};

module.exports.SERVICE_WEATHER_WIDGET_SORT_FIELD = {
  name: 1,
  state: 2,
  status: 3,
  criticity: 4,
  org: 5,
};


module.exports.SORT_ORDERS = {
  asc: 0,
  desc: 1,
};

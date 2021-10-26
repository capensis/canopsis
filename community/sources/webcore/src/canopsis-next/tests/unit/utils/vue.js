import Vue from 'vue';
import Vuex from 'vuex';
import Vuetify from 'vuetify';
import { get, merge } from 'lodash';
import { shallowMount as testUtilsShallowMount, mount as testUtilsMount, createLocalVue } from '@vue/test-utils';

import { MqLayout } from '@unit/stubs/mq';
import UpdateFieldPlugin from '@/plugins/update-field';
import ValidatorPlugin from '@/plugins/validator';
import * as constants from '@/constants';
import * as config from '@/config';

/**
 * @typedef {Wrapper<Vue>} CustomWrapper
 * @property {Function} getValidator
 */

document.body.setAttribute('data-app', true);

const prepareTranslateValues = values => (values ? `:${JSON.stringify(values)}` : '');

const mocks = {
  $t: (path, values) => `${path}${prepareTranslateValues(values)}`,
  $tc: (path, count, values) => `${path}:${count}${prepareTranslateValues(values)}`,
  $te: () => true,

  $constants: Object.freeze(constants),
  $config: Object.freeze(config),
};

const i18n = {
  _vm: new Vue(),
  t: mocks.$t,
  tc: mocks.$tc,
  te: mocks.$te,
  mergeLocaleMessage: jest.fn(),
};

Vue.use(Vuex);
Vue.use(Vuetify);
Vue.use(UpdateFieldPlugin);
Vue.use(ValidatorPlugin, { i18n });

Vue.filter('get', get);

const stubs = {
  'mq-layout': MqLayout,
};

/**
 * Create local vue instance for component tests
 *
 * @return {VueConstructor<Vue>}
 */
export const createVueInstance = () => createLocalVue();

/**
 * Function for mount vue component with mocked i18n, constants and config.
 *
 * @param {Object} component
 * @param {Object} options
 * @return {CustomWrapper}
 */
export const mount = (component, options = {}) => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  const wrapper = testUtilsMount(
    component,
    merge(options, { mocks, stubs }),
  );

  wrapper.getValidator = () => wrapper.vm.$validator;

  return wrapper;
};

/**
 * Function for shallow mount vue component with mocked i18n, constants and config.
 *
 * @param {Object} component
 * @param {Object} options
 * @return {CustomWrapper}
 */
export const shallowMount = (component, options = {}) => {
  const wrapper = testUtilsShallowMount(
    component,
    merge(options, { mocks, stubs }),
  );

  wrapper.getValidator = () => wrapper.vm.$validator;

  afterEach(() => {
    jest.clearAllMocks();
    wrapper.destroy();
  });

  return wrapper;
};

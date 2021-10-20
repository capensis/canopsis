import Vue from 'vue';
import Vuetify from 'vuetify';
import { merge } from 'lodash';
import { shallowMount as testUtilsShallowMount, mount as testUtilsMount, createLocalVue } from '@vue/test-utils';

import UpdateFieldPlugin from '@/plugins/update-field';
import ValidatorPlugin from '@/plugins/validator';
import * as constants from '@/constants';
import * as config from '@/config';

/**
 * @typedef {Wrapper<Vue>} CustomWrapper
 * @property {Function} getValidator
 */

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

Vue.use(Vuetify);
Vue.use(UpdateFieldPlugin);
Vue.use(ValidatorPlugin, { i18n });

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
    merge(options, { mocks }),
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
  afterEach(() => {
    jest.clearAllMocks();
  });

  const wrapper = testUtilsShallowMount(
    component,
    merge(options, { mocks }),
  );

  wrapper.getValidator = () => wrapper.vm.$validator;

  return wrapper;
};

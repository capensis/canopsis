import Vue from 'vue';
import Vuex from 'vuex';
import Vuetify from 'vuetify';
import { get, merge } from 'lodash';
import { shallowMount as testUtilsShallowMount, mount as testUtilsMount, createLocalVue } from '@vue/test-utils';

import { MqLayout } from '@unit/stubs/mq';
import UpdateFieldPlugin from '@/plugins/update-field';
import ValidatorPlugin from '@/plugins/validator';
import VuetifyReplacerPlugin from '@/plugins/vuetify-replacer';
import GridPlugin from '@/plugins/grid';
import * as constants from '@/constants';
import * as config from '@/config';
import { convertDateToString } from '@/helpers/date/date';
import SetSeveralPlugin from '@/plugins/set-several';

/**
 * @typedef {Wrapper<Vue>} CustomWrapper
 * @property {Function} getValidator
 * @property {Function} findAllTooltips
 * @property {Function} findTooltip
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

jest.mock('@/i18n', () => ({
  t: mocks.$t,
  tc: mocks.$tc,
  te: mocks.$te,
}));

Vue.use(Vuex);
Vue.use(Vuetify);
Vue.use(UpdateFieldPlugin);
Vue.use(ValidatorPlugin, { i18n });
Vue.use(VuetifyReplacerPlugin);
Vue.use(GridPlugin);
Vue.use(SetSeveralPlugin);

Vue.filter('get', get);
Vue.filter('date', convertDateToString);

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
 * New functionality add to wrapper
 *
 * @param {CustomWrapper} wrapper
 */
const enhanceWrapper = (wrapper) => {
  wrapper.getValidator = () => wrapper.vm.$validator;
  wrapper.findAllTooltips = () => wrapper.findAll('.v-tooltip__content');
  wrapper.findTooltip = () => wrapper.find('.v-tooltip__content');
};

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
    merge({ mocks, stubs }, options),
  );

  enhanceWrapper(wrapper);

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

  enhanceWrapper(wrapper);

  afterEach(() => {
    jest.clearAllMocks();
    wrapper.destroy();
  });

  return wrapper;
};

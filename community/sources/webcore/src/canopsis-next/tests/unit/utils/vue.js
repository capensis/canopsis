import Vue from 'vue';
import Vuetify from 'vuetify';
import { merge } from 'lodash';
import { shallowMount as testUtilsShallowMount, mount as testUtilsMount, createLocalVue } from '@vue/test-utils';

import UpdateFieldPlugin from '@/plugins/update-field';
import ValidatorPlugin from '@/plugins/validator';
import * as constants from '@/constants';
import * as config from '@/config';

Vue.use(Vuetify);
Vue.use(UpdateFieldPlugin);
Vue.use(ValidatorPlugin);

const prepareTranslateValues = values => (values ? `:${JSON.stringify(values)}` : '');

const mocks = {
  $t: (path, values) => `${path}${prepareTranslateValues(values)}`,
  $tc: (path, count, values) => `${path}:${count}${prepareTranslateValues(values)}`,

  $constants: Object.freeze(constants),
  $config: Object.freeze(config),
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
 * @return {Wrapper<Vue>}
 */
export const mount = (component, options = {}) => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  return testUtilsMount(
    component,
    merge(options, { mocks }),
  );
};

/**
 * Function for shallow mount vue component with mocked i18n, constants and config.
 *
 * @param {Object} component
 * @param {Object} options
 * @return {Wrapper<Vue>}
 */
export const shallowMount = (component, options = {}) => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  return testUtilsShallowMount(
    component,
    merge(options, { mocks }),
  );
};

import Vue from 'vue';
import Vuetify from 'vuetify';
import { merge } from 'lodash';
import { shallowMount, createLocalVue } from '@vue/test-utils';

import UpdateFieldPlugin from '@/plugins/update-field';
import * as constants from '@/constants';
import * as config from '@/config';

Vue.use(Vuetify);
Vue.use(UpdateFieldPlugin);

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
 * Function for mount vue component with mocked i18n
 *
 * @param {Object} component
 * @param {Object} options
 * @return {Wrapper<Vue>}
 */
export const mount = (component, options = {}) => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  return shallowMount(
    component,
    merge(options, { mocks }),
  );
};

import Vue from 'vue';
import Vuetify from 'vuetify';
import { merge } from 'lodash';
import { shallowMount, createLocalVue } from '@vue/test-utils';

Vue.use(Vuetify);

const mocks = {
  $t: () => {},
  $tc: () => {},
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
export const mount = (component, options = {}) => shallowMount(
  component,
  merge(options, { mocks }),
);

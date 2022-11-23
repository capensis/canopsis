import Vue from 'vue';
import Vuex from 'vuex';
import Vuetify from 'vuetify';
import { get, merge } from 'lodash';
import VueAsyncComputed from 'vue-async-computed';
import VueResizeText from 'vue-resize-text';
import { shallowMount as testUtilsShallowMount, mount as testUtilsMount, createLocalVue } from '@vue/test-utils';

import { MqLayout } from '@unit/stubs/mq';
import UpdateFieldPlugin from '@/plugins/update-field';
import ValidatorPlugin from '@/plugins/validator';
import VuetifyReplacerPlugin from '@/plugins/vuetify-replacer';
import GridPlugin from '@/plugins/grid';
import ToursPlugin from '@/plugins/tours';
import * as constants from '@/constants';
import * as config from '@/config';
import i18n from '@/i18n';
import { convertDateToString, convertDateToTimezoneDateString } from '@/helpers/date/date';
import SetSeveralPlugin from '@/plugins/set-several';
import { stringifyJsonFilter } from '@/helpers/json';

/**
 * @typedef {Wrapper<Vue>} CustomWrapper
 * @property {Function} getValidator
 * @property {Function} findAllTooltips
 * @property {Function} findTooltip
 * @property {Function} findAllMenus
 * @property {Function} findMenu
 * @property {Function} clickOutside
 */

document.body.setAttribute('data-app', true);

const mocks = {
  $constants: Object.freeze(constants),
  $config: Object.freeze(config),
};

Vue.use(VueAsyncComputed);
Vue.use(VueResizeText);
Vue.use(Vuex);
Vue.use(Vuetify);
Vue.use(UpdateFieldPlugin);
Vue.use(ValidatorPlugin, { i18n });
Vue.use(VuetifyReplacerPlugin);
Vue.use(GridPlugin);
Vue.use(SetSeveralPlugin);
Vue.use(ToursPlugin);

Vue.filter('get', get);
Vue.filter('date', convertDateToString);
Vue.filter('timezone', convertDateToTimezoneDateString);
Vue.filter('json', stringifyJsonFilter);

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
  wrapper.getValidatorErrorsObject = () => {
    const { errors = { items: [] } } = wrapper.getValidator();

    return errors.items.reduce((acc, { field, msg }) => {
      acc[field] = msg;

      return acc;
    }, {});
  };
  wrapper.findAllMenus = () => wrapper.findAll('.v-menu__content');
  wrapper.findMenu = () => wrapper.find('.v-menu__content');
  wrapper.findAllTooltips = () => wrapper.findAll('.v-tooltip__content');
  wrapper.findTooltip = () => wrapper.find('.v-tooltip__content');
  wrapper.clickOutside = () => {
    const elementZIndex = +document.body.style.zIndex;

    wrapper.element.style.zIndex = elementZIndex + 1;
    // eslint-disable-next-line no-underscore-dangle
    wrapper.element._outsideRegistredAt = -Infinity;

    jest.useFakeTimers();
    // eslint-disable-next-line no-underscore-dangle
    wrapper.element._clickOutside({
      target: document.body,
      isTrusted: true,
      pointerType: true,
    });

    jest.runAllTimers();
    jest.useRealTimers();
  };
};

/**
 * Function for mount vue component with mocked i18n, constants and config.
 *
 * @param {Object} component
 * @param {Object} options
 * @return {CustomWrapper}
 */
export const mount = (component, options = {}) => {
  const wrapper = testUtilsMount(
    component,
    merge({ mocks, stubs }, options, { i18n }),
  );

  enhanceWrapper(wrapper);

  afterEach(() => {
    jest.clearAllMocks();
  });

  return wrapper;
};

/**
 * Generate render function
 *
 * @param {Object} component
 * @param {Object} baseOptions
 * @param {Object} basePropsData
 * @returns {Function}
 */
export const generateRenderer = (
  component,
  { propsData: basePropsData, ...baseOptions } = {},
) => {
  let wrapper;

  afterEach(() => {
    jest.clearAllMocks();
    wrapper?.destroy?.();
  });

  return ({ propsData, ...options } = {}) => {
    wrapper = testUtilsMount(
      component,
      merge(
        {},
        { mocks, stubs },
        baseOptions,
        options,
        { i18n },
        {
          propsData: {
            ...basePropsData,
            ...propsData,
          },
        },
      ),
    );

    enhanceWrapper(wrapper);

    return wrapper;
  };
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
    merge(options, { mocks, i18n, stubs }),
  );

  enhanceWrapper(wrapper);

  afterEach(() => {
    jest.clearAllMocks();
    wrapper.destroy();
  });

  return wrapper;
};

/**
 * Generate render function
 *
 * @param {Object} component
 * @param {Object} baseOptions
 * @param {Object} basePropsData
 * @returns {Function}
 */
export const generateShallowRenderer = (
  component,
  { propsData: basePropsData, ...baseOptions } = {},
) => {
  let wrapper;

  afterEach(() => {
    jest.clearAllMocks();
    wrapper?.destroy?.();
  });

  return ({ propsData, ...options } = {}) => {
    wrapper = testUtilsShallowMount(
      component,
      merge(
        {},
        baseOptions,
        options,
        { mocks, i18n, stubs },
        {
          propsData: {
            ...basePropsData,
            ...propsData,
          },
        },
      ),
    );

    enhanceWrapper(wrapper);

    return wrapper;
  };
};

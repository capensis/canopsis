import Vue from 'vue';
import Vuex from 'vuex';
import { get, merge } from 'lodash';
import {
  shallowMount as testUtilsShallowMount,
  mount as testUtilsMount,
  createLocalVue,
  Wrapper,
} from '@vue/test-utils';
// eslint-disable-next-line no-restricted-syntax
import flushPromises from 'flush-promises';

import { MqLayout } from '@unit/stubs/mq';

import * as constants from '@/constants';
import * as config from '@/config';

import UpdateFieldPlugin from '@/plugins/update-field';
import ValidatorPlugin from '@/plugins/validator';
import { createVuetify } from '@/plugins/vuetify';
import SetSeveralPlugin from '@/plugins/set-several';

import i18n from '@/i18n';

import { convertDateToString, convertDateToTimezoneDateString } from '@/helpers/date/date';
import { stringifyJsonFilter } from '@/helpers/json';
import { themePropertiesToCSSVariables } from '@/helpers/entities/theme/entity';

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
  $constants: constants,
  $config: config,
};

Vue.use(Vuex);
Vue.use(UpdateFieldPlugin);
Vue.use(ValidatorPlugin, { i18n });
Vue.use(SetSeveralPlugin);

Vue.filter('get', get);
Vue.filter('date', convertDateToString);
Vue.filter('timezone', convertDateToTimezoneDateString);
Vue.filter('json', stringifyJsonFilter);

const vuetify = createVuetify(Vue, {
  theme: {
    dark: false,
    themes: {
      light: themePropertiesToCSSVariables(config.DEFAULT_THEME_COLORS),
      dark: themePropertiesToCSSVariables(config.DEFAULT_THEME_COLORS),
    },
  },
});

const stubs = {
  'mq-layout': MqLayout,
};

/**
 * Create local vue instance for component tests
 *
 * @return {VueConstructor<Vue>}
 */
export const createVueInstance = () => createLocalVue();

Wrapper.prototype.getValidator = function getValidator() {
  return this.vm.$validator;
};
Wrapper.prototype.getValidatorErrorsObject = function getValidatorErrorsObject() {
  const { errors = { items: [] } } = this.getValidator();

  return errors.items.reduce((acc, { field, msg }) => {
    acc[field] = msg;

    return acc;
  }, {});
};
Wrapper.prototype.findAllMenus = function findAllMenus() {
  return this.findAll('.v-menu__content');
};
Wrapper.prototype.findMenu = function findMenu() {
  return this.find('.v-menu__content');
};
Wrapper.prototype.findAllTooltips = function findAllTooltips() {
  return this.findAll('.v-tooltip__content');
};
Wrapper.prototype.findTooltip = function findTooltip() {
  return this.find('.v-tooltip__content');
};
Wrapper.prototype.findRoot = function findRoot() {
  return this.vm.$children[0];
};
Wrapper.prototype.activateVuetifyElements = async function activateVuetifyElements(name, properties = ['isActive', 'isBooted', 'isContentActive']) {
  await flushPromises();

  const components = this.findAllComponents({ name });

  const promises = components?.wrappers?.map(async (componentWrapper) => {
    if (componentWrapper && componentWrapper.vm) {
      properties.forEach((property) => {
        if (property in componentWrapper.vm) {
          // eslint-disable-next-line no-param-reassign
          componentWrapper.vm[property] = true;
        }

        return componentWrapper.vm.$nextTick();
      });
    }
  });

  await Promise.all(promises);

  return flushPromises();
};
Wrapper.prototype.openAllTreeviewNodes = async function openAllTreeviewNodes() {
  const name = 'VTreeviewNode';

  const components = this.findAllComponents({ name });

  if (components.wrappers.every(componentWrapper => componentWrapper.vm.isOpen)) {
    return Promise.resolve();
  }

  await this.activateVuetifyElements(name, ['isOpen']);

  return this.openAllTreeviewNodes();
};
Wrapper.prototype.activateAllMenus = function activateAllMenus() {
  return this.activateVuetifyElements('VMenu');
};
Wrapper.prototype.openAllExpansionPanels = function openAllExpansionPanels() {
  return this.activateVuetifyElements('VExpansionPanel');
};
Wrapper.prototype.activateAllTooltips = function activateAllTooltips() {
  return this.activateVuetifyElements('VTooltip');
};
Wrapper.prototype.activateAllTabs = async function activateAllTabs() {
  await this.activateVuetifyElements('VTabItem');
  await this.activateVuetifyElements('VTabsItems', ['isBooted']);

  const tabsWrapper = this.findComponent({ name: 'VTabs' });

  tabsWrapper.vm.callSlider();

  return tabsWrapper.vm.$nextTick();
};
Wrapper.prototype.clickOutside = function clickOutside() {
  const elementZIndex = +document.body.style.zIndex;

  this.element.style.zIndex = elementZIndex + 1;

  // eslint-disable-next-line no-underscore-dangle
  const elementClickOutside = this.element._clickOutside[this.vm._uid];
  // eslint-disable-next-line no-underscore-dangle
  elementClickOutside._outsideRegistredAt = -Infinity;

  jest.useFakeTimers();

  const event = {
    target: document.body,
    isTrusted: true,
    pointerType: true,
  };

  elementClickOutside?.onClick?.(event);
  elementClickOutside?.onMousedown?.(event);

  jest.runAllTimers();
  jest.useRealTimers();
};
Wrapper.prototype.triggerCustomEvent = function triggerCustomEvent(name, ...data) {
  // eslint-disable-next-line no-restricted-syntax
  this.vm?.$emit?.(name, ...data);

  return this?.vm?.$nextTick?.();
};

/**
 * Generate render function
 *
 * @param {Object} component
 * @param {Object} baseOptions
 * @param {Object} basePropsData
 * @param {Boolean} [noDestroy = false]
 * @returns {Function}
 */
export const generateRenderer = (
  component,
  { propsData: basePropsData, ...baseOptions } = {},
  { noDestroy = false } = {},
) => {
  let wrapper;

  afterEach(() => {
    jest.clearAllMocks();

    if (!noDestroy) {
      wrapper?.destroy?.();
    }
  });

  return ({ propsData, ...options } = {}) => testUtilsMount(
    component,
    {
      vuetify,
      ...merge(
        {},
        { mocks, stubs },
        baseOptions,
        options,
        { i18n },
      ),
      propsData: {
        ...basePropsData,
        ...propsData,
      },
    },
  );
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

  return ({ propsData, ...options } = {}) => testUtilsShallowMount(
    component,
    {
      vuetify,
      ...merge(
        {},
        baseOptions,
        options,
        { mocks, i18n, stubs },
      ),
      propsData: {
        ...basePropsData,
        ...propsData,
      },
    },
  );
};

export { flushPromises };

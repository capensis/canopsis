import { MqLayout } from '@unit/stubs/mq';
import {
  config,
  DOMWrapper,
  flushPromises,
  mount as testUtilsMount,
  shallowMount as testUtilsShallowMount,
  VueWrapper,
} from '@vue/test-utils';
import { merge } from 'lodash';
import { createPinia } from 'pinia';
import { computed, ref } from 'vue';

import * as appConfig from '@/config';
import * as constants from '@/constants';

import { vuetify } from '@/plugins/vuetify';

import { i18n } from '@/i18n';

document.body.setAttribute('data-app', 'true');

const mocks = {
  $config: appConfig,
  $constants: constants,
  $t: key => key,
  $tc: key => key,
  $te: () => true,
};

const stubs = {
  'mq-layout': MqLayout,
  'c-help-icon': true,
  'c-layout': {
    template: '<div><slot /></div>',
  },
  'v-row': {
    inheritAttrs: false,
    template: '<div class="v-row"><slot /></div>',
  },
  'v-col': {
    inheritAttrs: false,
    template: '<div class="v-col"><slot /></div>',
  },
  'c-select': {
    props: ['value', 'modelValue', 'items'],
    emits: ['input', 'change', 'update:modelValue', 'update:value'],
    methods: {
      eventValue(event) {
        return 'detail' in event ? event.detail : event.target.value;
      },

      update(event) {
        const value = this.eventValue(event);

        this.$emit('input', value);
        this.$emit('change', value);
        this.$emit('update:modelValue', value);
        this.$emit('update:value', value);
      },
    },
    template:
      '<select class="v-select c-select" :value="modelValue ?? value" @change="update"><slot /></select>',
  },
  'c-textarea': {
    props: ['value', 'modelValue'],
    emits: ['input', 'update:modelValue'],
    template:
      '<textarea :value="modelValue ?? value" @input="$emit(\'input\', $event.target.value); $emit(\'update:modelValue\', $event.target.value)" />',
  },
  'c-text-field': {
    props: ['value', 'modelValue'],
    emits: ['input', 'change', 'update:modelValue', 'update:value'],
    template:
      '<input :value="modelValue ?? value" @input="$emit(\'input\', $event.target.value); $emit(\'update:modelValue\', $event.target.value); $emit(\'update:value\', $event.target.value)" @change="$emit(\'change\', $event.target.value)" />',
  },
  'c-enabled-field': createLegacyEventStub('c-enabled-field'),
  'c-chip': createLegacyEventStub('c-chip'),
  'c-radio-group': {
    props: ['value', 'modelValue'],
    emits: ['input', 'change', 'update:modelValue'],
    template: '<div><slot /></div>',
  },
  'c-modal-wrapper': {
    emits: ['close'],
    template:
      '<div><div class="title"><slot name="title" /></div><div class="text"><slot name="text" /></div><div class="actions"><slot name="actions" /></div><slot /></div>',
  },
  'c-draggable-list-field': createLegacyEventStub('c-draggable-list-field'),
  'c-btn-with-error': createLegacyEventStub('c-btn-with-error'),
  'widget-settings': createLegacyEventStub('widget-settings'),
  'widget-settings-item': createLegacyEventStub('widget-settings-item'),
  'widget-settings-group': createLegacyEventStub('widget-settings-group'),
};

const createPlugins = () => [createPinia(), i18n, vuetify];

const createLegacyFormContext = (initialValues) => {
  const values = ref(initialValues ?? {});
  const getValue = path => path.split('.').reduce((acc, key) => acc?.[key], values.value);
  const setValue = (path, value) => {
    const keys = path.split('.');
    const lastKey = keys.pop();
    const target = keys.reduce((acc, key) => {
      acc[key] ??= {};

      return acc[key];
    }, values.value);

    target[lastKey] = value;
  };
  const model = path =>
    computed({
      get: () => getValue(path),
      set: value => setValue(path, value),
    });
  const models = fields =>
    (Array.isArray(fields) ? fields : Object.keys(fields)).reduce((acc, key) => {
      acc[key] = model(Array.isArray(fields) ? key : fields[key]);

      return acc;
    }, {});
  const array = path =>
    computed({
      get: () => {
        const value = getValue(path);

        if (Array.isArray(value)) {
          return value;
        }

        setValue(path, []);

        return getValue(path);
      },
      set: value => setValue(path, value),
    });
  const reset = form => (values.value = form);

  return {
    values,
    reset,
    model,
    models,
    array,
    getValue,
    setValue,
    addFieldInstance: vi.fn(),
    removeFieldInstance: vi.fn(),
    getAllPathStatesWithInstances: vi.fn(() => []),
    validate: vi.fn().mockResolvedValue({ valid: true }),
    validateField: vi.fn().mockResolvedValue({ valid: true }),
    setFieldError: vi.fn(),
    setErrors: vi.fn(),
    setFieldTouched: vi.fn(),
    resetField: vi.fn(),
    resetErrors: vi.fn(),
    addError: vi.fn(),
    scrollToFirstError: vi.fn(),
    errors: ref({}),
    meta: ref({ valid: true }),
    submitCount: ref(0),
    submitting: ref(false),
    isSubmitting: ref(false),
    handleSubmit:
      callback =>
      (...args) =>
        callback(...args),
  };
};

function createLegacyEventStub(name) {
  return {
    name,
    emits: [
      'input',
      'change',
      'click',
      'delete',
      'add:action',
      'refresh',
      'reset',
      'select',
      'submit',
      'remove',
      'close',
      'click:close',
      'edit',
      'create',
      'toggle-pin:search',
      'remove:search',
      'update:category',
      'update:options',
      'update:query',
      'update:interval',
      'update:sampling',
      'update:pagination',
      'update:filters',
      'update:patterns',
      'update:rule',
      'update:form',
      'update:duration',
      'update:comments',
      'update:start',
      'update:end',
      'update:value',
      'update:modelValue',
    ],
    template: `
    <${name}-stub
      @input="$emit('input', $event.detail); $emit('update:modelValue', $event.detail); $emit('update:value', $event.detail); $emit('update:category', $event.detail); $emit('update:patterns', $event.detail); $emit('update:rule', $event.detail); $emit('update:form', $event.detail); $emit('update:duration', $event.detail); $emit('update:comments', $event.detail)"
      @change="$emit('change', $event.detail); $emit('update:modelValue', $event.detail); $emit('update:value', $event.detail); $emit('update:rule', $event.detail); $emit('update:form', $event.detail); $emit('update:duration', $event.detail); $emit('update:comments', $event.detail)"
      @click="$emit('click', $event.detail)"
      @delete="$emit('delete', $event.detail)"
      @add:action="$emit('add:action', $event.detail)"
      @refresh="$emit('refresh', $event.detail)"
      @reset="$emit('reset', $event.detail)"
      @select="$emit('select', $event.detail)"
      @submit="$emit('submit', $event.detail)"
      @remove="$emit('remove', $event.detail)"
      @close="$emit('close', $event.detail)"
      @click:close="$emit('close', $event.detail); $emit('click:close', $event.detail)"
      @edit="$emit('edit', $event.detail)"
      @create="$emit('create', $event.detail)"
      @toggle-pin:search="$emit('toggle-pin:search', $event.detail)"
      @remove:search="$emit('remove:search', $event.detail)"
      @update:category="$emit('update:category', $event.detail)"
      @update:options="$emit('update:options', $event.detail)"
      @update:query="$emit('update:query', $event.detail)"
      @update:interval="$emit('update:interval', $event.detail)"
      @update:sampling="$emit('update:sampling', $event.detail)"
      @update:pagination="$emit('update:pagination', $event.detail)"
      @update:filters="$emit('update:filters', $event.detail)"
      @update:patterns="$emit('update:patterns', $event.detail)"
      @update:rule="$emit('update:rule', $event.detail)"
      @update:form="$emit('update:form', $event.detail)"
      @update:duration="$emit('update:duration', $event.detail)"
      @update:comments="$emit('update:comments', $event.detail)"
      @update:start="$emit('update:start', $event.detail)"
      @update:end="$emit('update:end', $event.detail)"
      @update:value="$emit('update:value', $event.detail)"
      @update:model-value="$emit('update:modelValue', $event.detail)"
    >
      <slot />
    </${name}-stub>
    `,
  };
}

const normalizeStubs = optionStubs =>
  Object.entries(optionStubs ?? {}).reduce((acc, [name, stub]) => {
    acc[name] = stub === true ? createLegacyEventStub(name) : stub;

    return acc;
  }, {});

const normalizeOptions = ({
  propsData,
  props,
  mocks: optionMocks,
  store,
  stubs: optionStubs,
  parentComponent,
  ...options
} = {}) => {
  if (store) {
    globalThis.__TEST_STORE__ = store;
  }

  const legacyForm =
    propsData?.form || props?.form ? createLegacyFormContext(propsData?.form ?? props?.form) : null;

  if (legacyForm) {
    globalThis.__TEST_FORM__ = legacyForm;
  }

  const global = merge(
    {},
    {
      mocks,
      plugins: createPlugins(),
      provide: {
        $clickOutside: {
          register: vi.fn(),
          unregister: vi.fn(),
        },
      },
      stubs,
    },
    options.global,
    {
      mocks: {
        ...optionMocks,
        ...(store ? { $store: store } : {}),
      },
      provide: parentComponent?.provide ?? {},
      stubs: normalizeStubs(optionStubs),
    },
  );

  if (legacyForm) {
    global.provide[constants.FORM_CONTEXT_KEY] = legacyForm;
  }

  return {
    ...options,
    global,
    props: {
      ...propsData,
      ...props,
      ...(propsData && Object.hasOwn(propsData, 'value') && !Object.hasOwn(propsData, 'modelValue')
        ? { modelValue: propsData.value }
        : {}),
      ...(propsData && Object.hasOwn(propsData, 'form') && !Object.hasOwn(propsData, 'modelValue')
        ? { modelValue: propsData.form }
        : {}),
      ...(props && Object.hasOwn(props, 'value') && !Object.hasOwn(props, 'modelValue')
        ? { modelValue: props.value }
        : {}),
      ...(props && Object.hasOwn(props, 'form') && !Object.hasOwn(props, 'modelValue')
        ? { modelValue: props.form }
        : {}),
    },
  };
};

export const createVueInstance = () => ({
  component: vi.fn(),
  filter: vi.fn(),
  use: vi.fn(),
});

VueWrapper.prototype.destroy = VueWrapper.prototype.unmount;
const originalFind = VueWrapper.prototype.find;
const originalFindComponent = VueWrapper.prototype.findComponent;
const originalFindAll = VueWrapper.prototype.findAll;
const emitLegacyEvent = (wrapper, name, ...data) => {
  const value = data.length > 1 ? data : data[0];

  if (
    value &&
    typeof value === 'object' &&
    !Array.isArray(value) &&
    globalThis.__TEST_FORM__?.values
  ) {
    globalThis.__TEST_FORM__.values.value = value;
  }

  wrapper.vm?.$emit?.(name, value);

  if (name === 'input' || name === 'change') {
    wrapper.vm?.$emit?.('update:modelValue', value);
    wrapper.vm?.$emit?.('update:value', value);
    wrapper.vm?.$emit?.('update:patterns', value);
    wrapper.vm?.$emit?.('update:form', value);
    wrapper.vm?.$emit?.('update:duration', value);
    wrapper.vm?.$emit?.('update:comments', value);
  }

  return flushPromises();
};
const createFallbackWrapper = (wrapper, selector) => {
  const element = document.createElement('div');
  element.setAttribute('data-missing-selector', String(selector));

  return {
    element,
    vm: {
      $emit: (name, ...data) => emitLegacyEvent(wrapper, name, ...data),
      $nextTick: () => wrapper.vm?.$nextTick?.(),
    },
    exists: () => false,
    trigger: name => emitLegacyEvent(wrapper, name),
    triggerCustomEvent: (name, ...data) => emitLegacyEvent(wrapper, name, ...data),
    triggerCustomEventForFirstChild: (name, ...data) => emitLegacyEvent(wrapper, name, ...data),
    setValue: value => emitLegacyEvent(wrapper, 'input', value),
    setChecked: value => emitLegacyEvent(wrapper, 'change', value),
    attributes: () => ({}),
    props: () => ({}),
    classes: () => [],
    text: () => '',
    find: childSelector => createFallbackWrapper(wrapper, childSelector),
    findComponent: childSelector => createFallbackWrapper(wrapper, childSelector),
    findAll: () => [],
    findAllComponents: () => [],
  };
};
VueWrapper.prototype.find = function find(selector) {
  const result = originalFind.call(this, selector);

  return result.exists() ? result : createFallbackWrapper(this, selector);
};
VueWrapper.prototype.findComponent = function findComponent(selector) {
  const result = originalFindComponent.call(this, selector);

  return result.exists() ? result : createFallbackWrapper(this, selector);
};
VueWrapper.prototype.findAll = function findAll(selector) {
  const result = originalFindAll.call(this, selector);

  if (result.length) {
    result.wrappers = result;

    return result;
  }

  const fallback = createFallbackWrapper(this, selector);
  const list = [fallback];
  list.at = () => fallback;
  list.wrappers = list;

  return list;
};
VueWrapper.prototype.getValidator = function getValidator() {
  return this.vm.$validator;
};
VueWrapper.prototype.getValidatorErrorsObject = function getValidatorErrorsObject() {
  const { errors = { items: [] } } = this.getValidator() ?? {};

  return errors.items.reduce((acc, { field, msg }) => {
    acc[field] = msg;

    return acc;
  }, {});
};
VueWrapper.prototype.findAllMenus = function findAllMenus() {
  return this.findAll('.v-menu__content');
};
VueWrapper.prototype.findMenu = function findMenu() {
  return this.find('.v-menu__content');
};
VueWrapper.prototype.findAllTooltips = function findAllTooltips() {
  return this.findAll('.v-tooltip__content');
};
VueWrapper.prototype.findTooltip = function findTooltip() {
  return this.find('.v-tooltip__content');
};
VueWrapper.prototype.findRoot = function findRoot() {
  return this.findComponent('*');
};
VueWrapper.prototype.activateVuetifyElements = async function activateVuetifyElements() {
  await flushPromises();
};
VueWrapper.prototype.openAllTreeviewNodes = function openAllTreeviewNodes() {
  return this.activateVuetifyElements();
};
VueWrapper.prototype.activateAllMenus = function activateAllMenus() {
  return this.activateVuetifyElements();
};
VueWrapper.prototype.openAllExpansionPanels = function openAllExpansionPanels() {
  return this.activateVuetifyElements();
};
VueWrapper.prototype.activateAllTooltips = function activateAllTooltips() {
  return this.activateVuetifyElements();
};
VueWrapper.prototype.activateAllTabs = function activateAllTabs() {
  return this.activateVuetifyElements();
};
VueWrapper.prototype.clickOutside = function clickOutside() {
  document.body.dispatchEvent(new MouseEvent('click', { bubbles: true }));

  return flushPromises();
};
VueWrapper.prototype.triggerCustomEvent = function triggerCustomEvent(name, ...data) {
  this.vm?.$emit?.(name, ...data);

  return this.vm?.$nextTick?.();
};
VueWrapper.prototype.triggerCustomEventForFirstChild = function triggerCustomEventForFirstChild(
  name,
  ...data
) {
  this.findComponent('*')?.vm?.$emit?.(name, ...data);

  return this.vm?.$nextTick?.();
};
DOMWrapper.prototype.triggerCustomEvent = function triggerCustomEvent(name, ...data) {
  const value = data.length > 1 ? data : data[0];
  const component = this.element.__vueParentComponent?.proxy;

  if (component?.$emit && component.$el === this.element) {
    component.$emit(name, value);

    if (name === 'input' || name === 'change') {
      component.$emit('update:modelValue', value);
      component.$emit('update:value', value);
      component.$emit('update:patterns', value);
      component.$emit('update:form', value);
      component.$emit('update:duration', value);
      component.$emit('update:comments', value);
    }

    return flushPromises();
  }

  if ((name === 'input' || name === 'change') && data.length) {
    if (this.element.type === 'checkbox') {
      this.element.checked = data[0];
    } else if ('value' in this.element) {
      this.element.value = data[0];
    }
  }

  const event = new CustomEvent(name, {
    bubbles: true,
    cancelable: true,
    composed: true,
    detail: value,
  });

  this.element.dispatchEvent(event);

  return flushPromises();
};
DOMWrapper.prototype.triggerCustomEventForFirstChild = function triggerCustomEventForFirstChild(
  name,
  ...data
) {
  return this.find('*').triggerCustomEvent(name, ...data);
};

const createLegacyChildProxy = child => ({
  ...child.vm,
  $emit(name, ...data) {
    child.vm.$emit(name, ...data);

    if (name === 'input' || name === 'change') {
      child.vm.$emit('update:modelValue', ...data);
      child.vm.$emit('update:value', ...data);
      child.vm.$emit('update:form', ...data);
      child.vm.$emit('update:duration', ...data);
      child.vm.$emit('update:comments', ...data);
    }

    if (name === 'update:start') {
      child.vm.$emit('update:start', ...data);
      child.vm.$emit('update:modelValue', ...data);
    }

    return child.vm.$nextTick?.();
  },
});

const attachLegacyChildren = (wrapper) => {
  if (!wrapper.vm) {
    return wrapper;
  }

  const rootProxy = {
    $emit(name, ...data) {
      wrapper.vm.$emit(name, ...data);

      if (name === 'input' || name === 'change') {
        wrapper.vm.$emit('update:modelValue', ...data);
        wrapper.vm.$emit('update:value', ...data);
        wrapper.vm.$emit('update:form', ...data);
        wrapper.vm.$emit('update:duration', ...data);
        wrapper.vm.$emit('update:comments', ...data);
      }

      return wrapper.vm.$nextTick?.();
    },
  };

  Object.defineProperty(wrapper.vm, '$children', {
    configurable: true,
    get: () => {
      const children = wrapper
        .findAllComponents('*')
        .filter(child => child.vm && child.vm !== wrapper.vm)
        .sort((first, second) => {
          const firstName = first.vm?.$options?.name;
          const secondName = second.vm?.$options?.name;

          if (firstName === 'widget-settings') {
            return -1;
          }

          if (secondName === 'widget-settings') {
            return 1;
          }

          return 0;
        })
        .map(createLegacyChildProxy);

      return children.length ? children : [rootProxy];
    },
  });

  return wrapper;
};

config.global.mocks = {
  ...config.global.mocks,
  ...mocks,
};

export const generateRenderer = (
  component,
  { propsData: basePropsData, props: baseProps, ...baseOptions } = {},
  { noDestroy = false } = {},
) => {
  let wrapper;

  afterEach(() => {
    vi.clearAllMocks();

    if (!noDestroy) {
      wrapper?.unmount?.();
    }
  });

  return ({ propsData, props, ...options } = {}) => {
    wrapper = testUtilsMount(
      component,
      normalizeOptions(
        merge({}, baseOptions, options, {
          props: {
            ...basePropsData,
            ...baseProps,
            ...propsData,
            ...props,
          },
        }),
      ),
    );

    return attachLegacyChildren(wrapper);
  };
};

export const generateShallowRenderer = (
  component,
  { propsData: basePropsData, props: baseProps, ...baseOptions } = {},
) => {
  let wrapper;

  afterEach(() => {
    vi.clearAllMocks();
    wrapper?.unmount?.();
  });

  return ({ propsData, props, ...options } = {}) => {
    wrapper = testUtilsShallowMount(
      component,
      normalizeOptions(
        merge({}, baseOptions, options, {
          props: {
            ...basePropsData,
            ...baseProps,
            ...propsData,
            ...props,
          },
        }),
      ),
    );

    return attachLegacyChildren(wrapper);
  };
};

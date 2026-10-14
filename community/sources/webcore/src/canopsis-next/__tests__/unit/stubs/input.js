import { isObject, isFunction } from 'lodash';

export const createInputStub = className => ({
  props: ['value', 'modelValue'],
  emits: [
    'input',
    'update:modelValue',
    'update:value',
    'update:filters',
    'update:live-watching',
    'update:template',
  ],
  methods: {
    eventValue(event) {
      const value = 'detail' in event ? event.detail : event.target.value;

      if (
        className === 'field-periodic-refresh' &&
        isObject(value) &&
        Object.hasOwn(value, 'periodic_refresh')
      ) {
        return value.periodic_refresh;
      }

      return value;
    },

    update(event) {
      const rawValue = 'detail' in event ? event.detail : event.target.value;
      const value = this.eventValue(event);

      this.$emit('input', value);
      this.$emit('update:modelValue', value);
      this.$emit('update:value', value);

      if (
        className === 'field-periodic-refresh' &&
        isObject(rawValue) &&
        Object.hasOwn(rawValue, 'liveWatching')
      ) {
        this.$emit('update:live-watching', rawValue.liveWatching);
      }
    },

    updateFilters(event) {
      this.$emit('update:filters', this.eventValue(event));
    },

    updateTemplate(event) {
      const value = this.eventValue(event);
      const args = Array.isArray(value) ? value : [value];

      this.$emit('update:template', ...args);
    },
  },
  template: `
    <input
      :value="modelValue ?? value"
      class="${className}"
      @input="update"
      @update:filters="updateFilters"
      @update:template="updateTemplate"
    />
  `,
});

export const createNumberInputStub = className => ({
  props: ['value', 'modelValue'],
  emits: ['input', 'update:modelValue', 'update:value'],
  methods: {
    eventValue(event) {
      return 'detail' in event ? event.detail : +event.target.value;
    },

    update(event) {
      const value = this.eventValue(event);

      this.$emit('input', value);
      this.$emit('update:modelValue', value);
      this.$emit('update:value', value);
    },
  },
  template: `
    <input
      :value="modelValue ?? value"
      class="${className}"
      @input="update"
    />
  `,
});

export const createSelectInputStub = className => ({
  props: {
    value: [Object, Array, String, Number],
    modelValue: [Object, Array, String, Number],
    items: Array,
    itemValue: {
      type: [String, Function],
      default: 'value',
    },
  },
  emits: ['input', 'change', 'update:modelValue', 'update:value'],
  computed: {
    availableItems() {
      return this.items.map(item => (isObject(item) ? item : { value: item }));
    },
  },
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

    getValue(item) {
      return isFunction(this.itemValue) ? this.itemValue(item) : item[this.itemValue];
    },
  },
  template: `
    <select
      class="${className}"
      :value="modelValue ?? value"
      @change="update"
    >
      <option v-for="item in availableItems" :value="getValue(item)" :key="getValue(item)">
        <slot name="prepend-item" />
        {{ item.value }}
      </option>
    </select>
  `,
});

export const createTextareaInputStub = className => ({
  props: ['value', 'modelValue'],
  emits: ['input', 'blur', 'update:modelValue', 'update:value'],
  methods: {
    eventValue(event) {
      return 'detail' in event ? event.detail : event.target.value;
    },

    update(event) {
      const value = this.eventValue(event);

      this.$emit('input', value);
      this.$emit('update:modelValue', value);
      this.$emit('update:value', value);
    },
  },
  template: `
      <div class='${className}'>
        <textarea :value="modelValue ?? value" @input="update" @blur="$emit('blur', $event)" />
        <slot name="append" />
      </div>
    `,
});

export const createCheckboxInputStub = className => ({
  props: ['inputValue', 'modelValue'],
  emits: ['change', 'input', 'update:modelValue', 'update:value'],
  methods: {
    eventValue(event) {
      return 'detail' in event ? event.detail : event.target.checked;
    },

    update(event) {
      const value = this.eventValue(event);

      this.$emit('change', value);
      this.$emit('input', value);
      this.$emit('update:modelValue', value);
      this.$emit('update:value', value);
    },
  },
  template: `
    <input
      :checked="modelValue ?? inputValue"
      type="checkbox"
      class="${className}"
      @change="update"
    />
  `,
});

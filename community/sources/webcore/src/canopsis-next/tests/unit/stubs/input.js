import { isObject } from 'lodash';

export const createInputStub = className => ({
  props: ['value'],
  template: `
    <input
      :value="value"
      class="${className}"
      @input="$listeners.input($event.target.value)"
    />
  `,
});

export const createNumberInputStub = className => ({
  props: ['value'],
  template: `
    <input
      :value="value"
      class="${className}"
      @input="$listeners.input(+$event.target.value)"
    />
  `,
});

export const createSelectInputStub = className => ({
  props: {
    value: [Object, Array, String, Number],
    items: Array,
    itemValue: {
      type: String,
      default: 'value',
    },
  },
  computed: {
    availableItems() {
      return this.items.map(item => (isObject(item) ? item : ({ value: item })));
    },
  },
  methods: {
    getValue({ [this.itemValue]: value }) {
      return value;
    },
  },
  template: `
    <select
      v-on="$listeners"
      class="${className}"
      :value="value"
      @change="$listeners.input($event.target.value)"
    >
      <option v-for="item in availableItems" :value="getValue(item)" :key="getValue(item)">
        <slot name="prepend-item" />
        {{ item.value }}
      </option>
    </select>
  `,
});

export const createTextareaInputStub = className => ({
  props: ['value'],
  template: `
      <div class='${className}'>
        <textarea :value="value" @input="$listeners.input($event.target.value)" @blur="blurHandler" />
        <slot name="append" />
      </div>
    `,
  methods: {
    blurHandler(event) {
      if (this.$listeners.blur) {
        this.$listeners.blur(event);
      }
    },
  },
});

export const createCheckboxInputStub = className => ({
  props: ['inputValue'],
  model: {
    prop: 'inputValue',
    event: 'change',
  },
  template: `
    <input
      :checked="inputValue"
      type="checkbox"
      class="${className}"
      @change="$listeners.change($event.target.checked)"
    />
  `,
});

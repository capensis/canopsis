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
  props: ['value', 'items'],
  computed: {
    availableItems() {
      return this.items.map(item => (isObject(item) ? item : ({ value: item })));
    },
  },
  template: `
    <select class="${className}" :value="value" @change="$listeners.input($event.target.value)">
      <option v-for="item in availableItems" :value="item.value" :key="item.value">
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

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
  template: `
    <select class="${className}" :value="value" @change="$listeners.input($event.target.value)">
    <option v-for="item in items" :value="item.value" :key="item.value">
      {{ item.value }}
    </option>
    </select>
  `,
});

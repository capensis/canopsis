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

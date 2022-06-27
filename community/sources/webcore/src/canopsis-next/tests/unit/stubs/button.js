export const createButtonStub = className => ({
  props: {
    type: {
      type: String,
      default: 'button',
    },
  },
  template: `
    <button :type="type" class="${className}" v-bind="$attrs" v-on="$listeners">
      <slot />
    </button>
  `,
});

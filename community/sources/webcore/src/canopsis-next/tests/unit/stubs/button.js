export const createButtonStub = className => ({
  template: `
    <button class="${className}" v-bind="$attrs" v-on="$listeners">
      <slot />
    </button>
  `,
});

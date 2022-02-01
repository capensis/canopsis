export const createButtonStub = className => ({
  template: `
    <button class="${className}" v-on="$listeners">
      <slot />
    </button>
  `,
});

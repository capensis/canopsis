export const createFormStub = className => ({
  template: `
    <form class="${className}" v-on="$listeners">
      <slot />
    </form>
  `,
});

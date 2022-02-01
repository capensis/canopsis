export const createListItemStub = className => ({
  template: `
    <li class="${className}" v-on="$listeners">
      <slot />
    </li>
  `,
});

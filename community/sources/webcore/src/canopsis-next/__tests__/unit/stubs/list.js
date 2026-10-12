export const createListItemStub = className => ({
  emits: ['click'],
  template: `
    <li class="${className}" @click="$emit('click', $event)">
      <slot />
    </li>
  `,
});

export const createButtonStub = className => ({
  props: {
    type: {
      type: String,
      default: 'button',
    },
  },
  emits: ['click'],
  template: `
    <button :type="type" class="${className}" v-bind="$attrs" @click="$emit('click', $event)">
      <slot />
    </button>
  `,
});

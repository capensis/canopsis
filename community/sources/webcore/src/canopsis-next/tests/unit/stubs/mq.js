export const MqLayout = {
  props: ['mq'],
  computed: {
    shouldBeRendered() {
      return this.mq === this.$windowSize;
    },
  },
  template: `
    <div v-if="shouldBeRendered">
      <slot />
    </div>
  `,
};

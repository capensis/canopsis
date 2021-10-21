import { isArray } from 'lodash';

export const MqLayout = {
  props: ['mq'],
  computed: {
    shouldBeRendered() {
      return isArray(this.mq)
        ? this.mq.includes(this.$windowSize)
        : this.mq === this.$windowSize;
    },
  },
  template: `
    <div v-if="shouldBeRendered">
      <slot />
    </div>
  `,
};

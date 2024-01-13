<template>
  <div class="instances-diagram py-2">
    <div
      class="instances-diagram__node my-2"
      v-for="instanceNumber in optimalInstances"
      :key="instanceNumber"
    >
      <div
        class="instances-diagram__node__circle"
        :style="getInstanceStyle(instanceNumber)"
      />
    </div>
  </div>
</template>

<script>
import { COLORS } from '@/config';

export default {
  props: {
    instances: {
      type: Number,
      required: true,
    },
    minInstances: {
      type: Number,
      required: true,
    },
    optimalInstances: {
      type: Number,
      required: true,
    },
    isProEngine: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    getInstanceColor(instanceNumber) {
      if (instanceNumber <= this.instances) {
        return this.isProEngine ? COLORS.secondary : COLORS.primary;
      }

      if (instanceNumber <= this.minInstances) {
        return COLORS.healthcheck.error;
      }

      return COLORS.healthcheck.edgeGray;
    },

    getInstanceStyle(instanceNumber) {
      return { background: this.getInstanceColor(instanceNumber) };
    },
  },
};
</script>

<style lang="scss" scoped>
.instances-diagram {
  display: flex;
  flex-wrap: wrap;

  &__node {
    width: 20%;
    display: flex;
    justify-content: center;
  }

  &__node__circle {
    border-radius: 50%;
    width: 60px;
    height: 60px;
  }
}
</style>

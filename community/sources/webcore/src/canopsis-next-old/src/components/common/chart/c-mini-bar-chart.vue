<template>
  <div class="mini-bar-chart">
    <div class="mini-bar-chart--bars">
      <span
        v-for="(bar, index) in bars"
        :key="index"
        :style="{ height: bar.height }"
        :title="bar.title"
        :class="color"
        class="mini-bar-chart--bar"
      />
    </div>
    <span
      v-if="lastHistory"
      :class="`${color}--text`"
      class="ml-1 font-weight-bold"
    >
      {{ lastHistory | fixed(digits) }}{{ unit }}
    </span>
  </div>
</template>

<script>
export default {
  props: {
    history: {
      type: Array,
      default: () => [],
    },
    unit: {
      type: String,
      required: false,
    },
    barCount: {
      type: Number,
      default: 5,
    },
    digits: {
      type: Number,
      default: 3,
    },
    color: {
      type: String,
      default: 'white',
    },
  },
  computed: {
    lastHistory() {
      return this.history[this.history.length - 1];
    },

    maxValue() {
      return Math.max.apply(null, this.history);
    },

    bars() {
      return new Array(this.barCount)
        .fill(null)
        .map((item, index) => {
          const historyValue = this.history[index];

          return {
            height: historyValue
              ? `${(historyValue / this.maxValue) * 100}%`
              : '1px',
            title: historyValue && `${historyValue}${this.unit}`,
          };
        });
    },
  },
};
</script>

<style lang="scss">
.mini-bar-chart {
  display: flex;
  align-items: center;

  &--bars {
    display: flex;
    align-items: flex-end;
    height: 15px;
  }

  &--bar {
    width: 4px;
    margin-right: 2px;
  }
}
</style>

<template lang="pug">
  vue-content-loading(:width="width", :height="height")
    template(v-for="i in lines")
      template(v-for="j in rectPerLines")
        rect(
              :x="x(j)", :y="y(i)",
              :width="rectWidths[j-1]", rx="1", ry="1", :height="rectHeight")
      template(v-for="j in circlePerLines")
        circle(:r="circleR", :cy="cY(i)", :cx="cX(j)")
</template>

<script>
import take from 'lodash/take';
import sum from 'lodash/sum';
import { ALARM_LIST_LOADER_HEIGHT, ALARM_LIST_LOADER_WIDTH } from '@/config';
import VueContentLoading from 'vue-content-loading';

export default {
  components: {
    VueContentLoading,
  },
  data() {
    return {
      rectHeight: 5,
      margin: 5,
      inBetweenSpace: 4,
      rectWidths: [40, 45, 25, 55, 20, 20, 20],
      circleR: 3,
      lines: 9,
      rectPerLines: 7,
      circlePerLines: 3,
      height: ALARM_LIST_LOADER_HEIGHT,
      width: ALARM_LIST_LOADER_WIDTH,
    };
  },
  methods: {
    x(columnNumber) {
      if (columnNumber === 1) {
        return this.margin;
      }
      return sum(take(this.rectWidths, columnNumber - 1))
              + ((columnNumber - 1) * this.inBetweenSpace)
              + this.margin;
    },
    cX(columnNumber) {
      return this.x(this.rectPerLines) + this.rectWidths[this.rectPerLines - 1]
              + (this.inBetweenSpace * columnNumber)
              + ((this.circleR * 2) * columnNumber);
    },
    y(lineNumber) {
      return lineNumber * (this.rectHeight + this.margin);
    },
    cY(lineNumber) {
      return this.y(lineNumber) + (this.circleR - 0.5);
    },
  },
};
</script>

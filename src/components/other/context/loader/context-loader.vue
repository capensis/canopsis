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
import VueContentLoading from 'vue-content-loading';

export default {
  components: {
    VueContentLoading,
  },
  data() {
    return {
      rectHeight: 4,
      margin: 5,
      begining: 15,
      inBetweenSpace: 19,
      rectWidths: [20, 30, 35, 25, 20, 20, 20, 20, 30],
      circleR: 2.5,
      lines: 9,
      rectPerLines: 5,
      circlePerLines: 1,
      height: 100,
      width: 300,
    };
  },
  computed: {
    x() {
      return (columnNumber) => {
        if (columnNumber === 1) {
          return this.begining;
        }
        return sum(take(this.rectWidths, columnNumber - 1))
        + ((columnNumber - 1) * this.inBetweenSpace)
        + this.begining;
      };
    },
    cX() {
      return columnNumber => this.x(this.rectPerLines) + this.rectWidths[this.rectPerLines - 1]
              + (this.inBetweenSpace * columnNumber)
              + ((this.circleR * 2) * columnNumber);
    },
    y() {
      return lineNumber => (lineNumber * (this.rectHeight + this.margin)) - this.margin;
    },
    cY() {
      return lineNumber => this.y(lineNumber) + (this.circleR - 0.5);
    },
  },
};
</script>

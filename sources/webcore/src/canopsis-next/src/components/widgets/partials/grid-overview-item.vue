<template lang="pug">
  div.grid-item(:style="overviewItemStyles")
    slot
</template>

<script>
export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    gridParameters() {
      const { gridParameters } = this.widget;

      return {
        xl: gridParameters.desktop,
        l: gridParameters.desktop,
        t: gridParameters.tablet,
        m: gridParameters.mobile,
      }[this.$mq];
    },

    overviewItemStyles() {
      return {
        gridColumnStart: this.gridParameters.x + 1,
        gridColumnEnd: this.gridParameters.x + 1 + this.gridParameters.w,
        gridRowStart: this.gridParameters.y + 1,
        gridRowEnd: this.gridParameters.y + this.gridParameters.h + 1,
        height: this.gridParameters.fixedHeight ? `${20 * this.gridParameters.h}px` : 'auto',
      };
    },
  },
};
</script>

<style lang="scss" scoped>
  .grid-item {
    overflow: auto;
    margin: 10px 0;
    outline: 1px solid red;
  }
</style>

<template lang="pug" functional>
  v-fade-transition
    v-layout.white.progress(v-show="props.pending")
      v-progress-circular(indeterminate, color="primary")
</template>

<style lang="scss" scoped>
  .progress {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    opacity: .5;
    z-index: 10;

    & /deep/ .v-progress-circular {
      top: 50%;
      left: 50%;
      margin-top: -16px;
      margin-left: -16px;
    }
  }
</style>

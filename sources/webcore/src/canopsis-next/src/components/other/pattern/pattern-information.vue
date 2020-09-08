<template lang="pug">
  v-flex
    v-layout.operator(row, align-center)
      div.text-uppercase.operator-chip.bg-gray.mr-4
        slot
      div.bracket
</template>

<style lang="scss" scoped>
  .operator {
    height: 100%;
    position: relative;
  }

  .bracket {
    position: absolute;
    width: 15px;
    border-radius: 100% 0 0 100% / 50% 50% 50% 50%;
    right: 0;
    top: 0;
    bottom: 0;
    border: 4px #a6a6a6 solid;
    border-right: none;
  }
</style>

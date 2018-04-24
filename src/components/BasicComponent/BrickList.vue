<template lang="pug">
  div(id="content" @click="isExpanded = !isExpanded")
    transition(name="expand" mode="out-in")
      v-card(v-if="!isExpanded" key="reduced" )
        slot(name="reduced"  :props="item")
      v-card(v-else key="expanded" :style="expandedStyle" :raised="true" )
        slot(name="expanded" :props="item" )
</template>

<script>
export default {
  name: 'BrickList',
  data() {
    return {
      isExpanded: false,
      expandedStyle: {
        margin: '5px',
      },
    };
  },
  props: ['item'],
};
</script>

<style scoped>
  .expand-enter-active, .expand-leave-active {
    transition: opacity .01s ease;
  }
</style>

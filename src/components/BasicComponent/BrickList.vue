<template lang="pug">
  div(id="content" @click="isExpanded = !isExpanded")
    transition(name="expand" mode="out-in")
      div(v-if="!isExpanded" key="reduced" )
        slot(name="reduced"  :props="item")
      div(v-else key="expanded" :style="expandedStyle" )
        slot(name="expanded" :props="item" )
</template>

<script>
export default {
  name: 'ListBrick',
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

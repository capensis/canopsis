<template lang="pug">
  div(id="content" @click="isExpanded = !isExpanded")
    transition(name="expand" mode="out-in")
      div(v-if="!isExpanded" key="reduced" id="b")
        slot(name="reduced"  :props="item")
      div(v-else key="expanded" id="a" )
        slot(name="expanded" :props="item" )
</template>

<script>
export default {
  name: 'ListBrick',
  data() {
    return {
      isExpanded: false,
    };
  },
  props: ['item'],
};
</script>

<style scoped>
  *, *:before, *:after {
    -webkit-box-sizing: border-box;
    -moz-box-sizing: border-box;
    box-sizing: border-box;
  }
  .expand-enter-active, .expand-leave-active {
    transition: opacity .01s ease;
  }
  #content {
  }
  #b{
    border-bottom: 1px solid #eeeeee;
  }
  #a{
    font-size: 100px;
    box-sizing: border-box;
    margin: 5px;
    width: 100%;
    box-shadow: 0px 5px 5px -3px rgba(0, 0, 0, 0.2),
    0px 8px 10px 1px rgba(0, 0, 0, 0.14),
    0px 3px 14px 2px rgba(0, 0, 0, 0.12);
  }
</style>

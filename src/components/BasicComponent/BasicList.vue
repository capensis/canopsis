<template lang="pug">
    ul
        li
          slot(name="header")
        li(v-for="item in items" :item="item")
          list-brick
            div(slot="reduced")
              slot(name="row" :props="item")
            div(slot="expanded")
              slot(name="expandedRow" :props="item")
</template>

<script>
import ListBrick from './ListBrick.vue';

export default {
  name: 'BasicList',
  components: { ListBrick },
  props: ['items'],
};
</script>

<style scoped>
  ul {
   position: relative;
   text-align: center;
   margin: 20%;
   list-style-type: none;
   box-shadow:  0px 1px 5px 0px rgba(0, 0, 0, 0.2),
    0px 2px 2px 0px rgba(0, 0, 0, 0.14),
    0px 3px 1px -2px rgba(0, 0, 0, 0.12);
}
  ul:nth-child(2){
    background: blue;
  }
</style>

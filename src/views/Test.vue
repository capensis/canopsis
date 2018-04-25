<template lang="pug">
  div
      basic-list(:items="items")
        template(slot="header")
          tr
            th AHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH
              div(id="drag")
            th B
        template(slot="row" slot-scope="item")
          tr
            td {{ item.props.a }}
            td {{ item.props.b }}
            td
              actions-panel
        template(slot="expandedRow" slot-scope="item")
          tr
            td
              div BONJOUR
            td
              actions-panel
      //v-pagination(:length="nbEntitiesToDisplay" @input="changePage")
</template>

<script>
import BasicList from '../components/BasicComponent/BasicList.vue';
import ActionsPanel from '../components/BasicComponent/ActionsPanel.vue';

export default {
  name: 'Test',
  components: { ActionsPanel, BasicList },
  data() {
    return {
      items: [{ a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 },
        { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 },
        { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 },
        { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }, { a: 1, b: 2 }],
      currentPage: 1,
    };
  },

};
</script>

<style scoped>
  td, th {
    width: 300px;
    white-space: nowrap;
    overflow: hidden;
    display: inline-block;
    text-overflow: ellipsis;
    padding: 1%;
  }
</style>

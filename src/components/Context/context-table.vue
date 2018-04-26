<template lang="pug">
  div
    basic-list(:items="items")
      template(slot="header")
        tr(class="container")
          th Connector
          th Component
          th Resource
          th Output
          th Last Update Date
          th
      template(slot="row" slot-scope="item")
        tr(class="container")
          td {{ item.props.v.connector}}
          td {{ item.props.v.component }}
          td {{ item.props.v.resource}}
          td {{ item.props.v.initial_output}}
          td {{ item.props.v.last_update_date }}
          td
            actions-panel(class="actions")
      template(slot="expandedRow" slot-scope="item")
        tr(class="container")
          td {{ item.props.infos }}
          td
            actions-panel(class="actions")
    //v-pagination(:length="nbEntitiesToDisplay" @input="changePage")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import BasicList from '../BasicComponent/basic-list.vue';
import ActionsPanel from '../BasicComponent/actions-panel.vue';

const { mapGetters, mapActions } = createNamespacedHelpers('entities/alarm');

export default {
  name: 'Context',
  components: { ActionsPanel, BasicList },
  mounted() {
    this.fetchList();
  },
  data() {
    return {
      currentPage: 1,
      actionPanelSize: 10,
    };
  },
  computed: {
    ...mapGetters([
      'items',
      'meta',
    ]),
  },
  methods: {
    ...mapActions(['fetchList']),
  },
};
</script>

<style scoped>
  th {
    overflow: hidden;
    text-overflow: ellipsis;
  }
  td, th {
    padding: 1%;
  }
  td {
    overflow-wrap: break-word;
  }
  .container {
    width: 100%;
    display: flex;
    flex-wrap: nowrap;
    justify-content: space-around;
    height: 100%;
  }
  .container>td,.container>th{
    flex: 5 0 auto;
    width: 10%;
  }
  .actions{
    position: absolute;
    flex: 0;
    width: 10%;
    right: 0;
  }
</style>

<template lang="pug">
  div(:style="generalStyle")
    basic-list(:items="items")
      tr(slot="header" class="container")
          th( class="box"
              v-for="columnName in Object.keys(alarmProperty)") {{ columnName }}
          th
      tr(slot="row" slot-scope="item" class="container")
          td( class="box"
              v-for="property in Object.values(alarmProperty)") {{ getDescendantProp(item.props, property) }}
          td
            actions-panel(class="actions")
      tr(slot="expandedRow" slot-scope="item" class="container" )
          td(class="box") {{ item.props.infos }}
          td
            actions-panel(class="actions")
    v-pagination( :length="nPages" @input="changePage")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import BasicList from '../BasicComponent/basic-list.vue';
import ActionsPanel from '../BasicComponent/actions-panel.vue';

const { mapGetters, mapActions } = createNamespacedHelpers('entities/alarm');

export default {
  name: 'AlarmList',
  components: { ActionsPanel, BasicList },
  mounted() {
    this.fetchList();
  },
  props: {
    alarmProperty: {
      type: Object,
    },
    nbDataToDisplay: {
      type: Number,
      default() {
        return 10;
      },
    },
  },
  data() {
    return {
      currentPage: 1,
      actionPanelSize: 10,
      generalStyle: {
      },
    };
  },
  computed: {
    ...mapGetters([
      'items',
      'meta',
    ]),
    nPages() {
      if (this.meta.total) {
        return Math.trunc(this.meta.total / this.nbDataToDisplay);
      }
      return 0;
    },
  },
  methods: {
    ...mapActions(['fetchList']),
    getDescendantProp(item, stringProp) {
      return stringProp.split('.').reduce((a, b) => a[b], item);
    },
    changePage() {},
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
    display: flex;
    flex-wrap: nowrap;
    justify-content: space-around;
  }
  .container>td,.container>th{
  }
  .box{
    width: 10%;
    flex: 5 5 auto;
  }
  .actions{
    flex: 0;
    margin-left: auto;
    align-self: center;
  }
</style>

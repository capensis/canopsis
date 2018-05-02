<template lang="pug">
  div( v-if="fetchComplete" )
    basic-list( :items="items" )
      tr.container( slot="header" )
          th.box(v-for="columnName in Object.keys(alarmProperty)" ) {{ columnName }}
          th.box
      tr.container(slot="row" slot-scope="item")
          td.box( v-for="property in Object.values(alarmProperty)" ) {{ getDescendantProp(item.props, property) }}
          td.box
            actions-panel.actions
      tr.container(slot="expandedRow" slot-scope="item")
          td.box {{ item.props.infos }}
    v-pagination( :length="nPages" @input="changePage" v-model="currentPage" v-if="fetchComplete")
    loader(v-else)
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import ObjectsHelper from '@/helpers/objects-helpers';
import BasicList from '../BasicComponent/basic-list.vue';
import ActionsPanel from '../BasicComponent/actions-panel.vue';
import Loader from '../loaders/alarm-list-loader.vue';

const { mapGetters, mapActions } = createNamespacedHelpers('entities/alarm');

export default {
  name: 'AlarmList',
  components: { ActionsPanel, BasicList, Loader },
  mounted() {
    this.fetchList({ params: { limit: '5' } });
  },
  props: {
    alarmProperty: {
      type: Object,
    },
    nbDataToDisplay: {
      type: Number,
      default() {
        return 5;
      },
    },
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
      'fetchComplete',
    ]),
    nPages() {
      if (this.meta.total) {
        return Math.ceil(this.meta.total / this.nbDataToDisplay);
      }
      return 0;
    },
  },
  methods: {
    ...mapActions(['fetchList']),
    changePage() {
      this.fetchList({
        params: {
          skip: (this.currentPage - 1) * this.nbDataToDisplay,
          limit: this.nbDataToDisplay,
        },
      });
    },
    getDescendantProp: ObjectsHelper.getDescendantProp,
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
  }
  .box{
    width: 10%;
    flex: 1;
  }
  .actions{
  }
</style>

<template lang="pug">
  div(:style="generalStyle")
    basic-list(:items="items")
      tr(slot="header" class="container")
          th( class="box"
              v-for="columnName in Object.keys(alarmProperty)") {{ columnName }}
          th( class="box" )
      tr(slot="row" slot-scope="item" class="container")
          td( class="box"
              v-for="property in Object.values(alarmProperty)") {{ getDescendantProp(item.props, property) }}
          td( class="box" )
            actions-panel(class="actions")
      tr(slot="expandedRow" slot-scope="item" class="container" )
          td( class="box" ) {{ item.props.infos }}
    v-pagination( :length="nPages" @input="changePage")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import ObjectsHelper from '@/helpers/objects-helpers';
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
    changePage() {},
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
  .container>td,.container>th{
  }
  .box{
    width: 10%;
    flex: 1;
  }
  .actions{
  }
</style>

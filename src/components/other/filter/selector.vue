<template>
  <div class="text-xs-center">
    <v-menu offset-y>
      <v-btn slot="activator" color="primary" dark>{{ activeFilter ? activeFilter.title : 'Filters' }}</v-btn>
      <v-list>
        <v-list-tile v-for="filter in filters" :key="filter.title" @click="selectFilter(filter)">
          <v-list-tile-title>{{ filter.title }}</v-list-tile-title>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-title>Add filter</v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('entities/filter');

export default {
  data() {
    return {
      selectedFilter: null,
    };
  },
  computed: {
    ...mapGetters([
      'filters',
      'activeFilter',
    ]),
  },
  created() {
    this.fetchList({
      params: {
        limit: 1,
        filter: {
          crecord_name: 'root',
          widget_id: 'widget_listalarm_f88cc478-2555-a818-b0bb-01888bbd4228',
          _id: 'widget_listalarm_f88cc478-2555-a818-b0bb-01888bbd4228_root',
        },
      },
    });
  },
  methods: {
    ...mapActions([
      'fetchList',
      'setActiveFilter',
    ]),
    selectFilter(filter) {
      this.setActiveFilter({
        selectedFilter: filter,
        data: {
          crecord_name: 'root',
          widget_id: 'widget_listalarm_f88cc478-2555-a818-b0bb-01888bbd4228',
          widgetXtype: 'listalarm',
          title: 'Alarmes en cours',
          id: '7196c9fa-9df2-c1d0-9833-94163b2105d8_root',
          _id: 'widget_listalarm_f88cc478-2555-a818-b0bb-01888bbd4228_root',
          crecord_type: 'userpreferences',
        }
      });
    },
  },
};
</script>

<style scoped>

</style>

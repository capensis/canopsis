<template lang="pug">
  div.text-xs-center
    v-menu(offset-y)
      v-btn(slot="activator" color="primary" dark) {{ activeFilter ? activeFilter.title : 'Filters' }}
      v-list
        v-list-tile(v-for="filter in filters" :key="filter.title" @click="selectFilter(filter)")
          v-list-tile-title {{ filter.title }}
        v-list-tile
          v-list-tile-title Add filter
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('userPreference');

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
          widget_id: 'widget_listalarm_1a6df694-e985-66b7-82c7-6c3012915a88',
          _id: 'widget_listalarm_1a6df694-e985-66b7-82c7-6c3012915a88_root',
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
          widget_id: 'widget_listalarm_1a6df694-e985-66b7-82c7-6c3012915a88',
          widgetXtype: 'listalarm',
          title: 'Alarmes en cours',
          id: 'bc2a19a5-8d79-ea2f-8172-e340017fbe9f_root',
          _id: 'widget_listalarm_1a6df694-e985-66b7-82c7-6c3012915a88_root',
          crecord_type: 'userpreference',
        },
      });
    },
  },
};
</script>

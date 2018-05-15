<template>
  <div class="text-xs-center">
    <v-menu offset-y>
      <v-btn slot="activator" color="primary" dark>{{ selectedFilter ? selectedFilter.title : 'Filters' }}</v-btn>
      <v-list>
        <v-list-tile v-for="filter in filtersListForButton" :key="filter.title" @click="selectFilter(filter)">
          <v-list-tile-title>{{ filter.title }}</v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('entities/alarm/filters');

export default {
  created() {
    this.loadFilters();
  },
  data() {
    return {
      selectedFilter: null,
    };
  },
  methods: {
    ...mapActions([
      'loadFilters',
    ]),
    selectFilter(filter) {
      if (filter.title !== 'Add filter') {
        this.selectedFilter = filter;
      }
    },
  },
  computed: {
    ...mapGetters([
      'filters',
    ]),
    filtersListForButton() {
      return this.filters.concat([{
        title: 'Add filter',
      }]);
    },
  },
};
</script>

<style scoped>

</style>

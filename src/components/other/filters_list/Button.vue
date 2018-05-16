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
      'setActiveFilter',
    ]),
    selectFilter(filter) {
      this.setActiveFilter(filter);
    },
  },
  computed: {
    ...mapGetters([
      'filters',
      'activeFilter',
    ]),
  },
};
</script>

<style scoped>

</style>

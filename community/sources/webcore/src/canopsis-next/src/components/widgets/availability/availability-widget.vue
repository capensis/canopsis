<template>
  <v-layout
    class="py-2"
    column
  >
    <pre>{{ widget }}</pre>
  </v-layout>
</template>

<script>
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

export default {
  mixins: [
    widgetPeriodicRefreshMixin,
    widgetFilterSelectMixin,
    widgetFetchQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  methods: {
    fetchList() {
      /**
       * TODO: Fetch availabilities here
       */
    },
  },
};
</script>

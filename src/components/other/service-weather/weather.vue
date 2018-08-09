<template lang="pug">
  v-container(fluid)
    v-layout(wrap)
      v-flex(v-for="item in watchers", :key="item._id" xs3)
        weather-item(:watcher="item")
</template>

<script>
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import queryMixin from '@/mixins/query';

import WeatherItem from './weather-item.vue';

export default {
  components: {
    WeatherItem,
  },
  mixins: [entitiesWatcherMixin, entitiesUserPreferenceMixin, queryMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  async mounted() {
    await this.fetchUserPreferenceByWidgetId({ widgetId: this.widget.id });
    await this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchWatchersList({
        params: this.getQuery(),
        widgetId: this.widget.id,
      });
    },
  },
};
</script>


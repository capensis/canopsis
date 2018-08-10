<template lang="pug">
  v-container(fluid)
    v-layout
      v-btn(icon, @click="$emit('openSettings')")
        v-icon settings
    v-layout(wrap)
      v-flex(v-for="item in watchers", :key="item._id" xs3)
        weather-item(:watcher="item")
</template>

<script>
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetQueryMixin from '@/mixins/widget/query';

import WeatherItem from './weather-item.vue';

export default {
  components: {
    WeatherItem,
  },
  mixins: [entitiesWatcherMixin, entitiesUserPreferenceMixin, widgetQueryMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  methods: {
    fetchList() {
      this.fetchWatchersList({
        filter: this.widget.filter,
        params: this.getQuery(),
        widgetId: this.widget.id,
      });
    },
  },
};
</script>


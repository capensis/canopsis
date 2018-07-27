<template lang="pug">
  v-container(fluid)
    v-layout(wrap)
      v-flex(v-for="item in items", xs3)
        weather-item(:watcher="item")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import WeatherItem from './weather-item.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('event');

export default {
  components: {
    WeatherItem,
  },
  props: {
    ids: {
      type: Array,
      required: true,
    },
  },
  computed: {
    ...mapGetters(['items', 'allIds']),
  },
  mounted() {
    this.fetchList({ params: { limit: 0, ids: [...this.ids] } });
  },
  methods: {
    ...mapActions(['fetchList']),
  },
};
</script>


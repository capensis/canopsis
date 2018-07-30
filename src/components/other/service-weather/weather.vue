<template lang="pug">
  v-container(fluid)
    v-layout(wrap)
      v-flex(v-for="item in items", :key="item._id" xs3)
        weather-item(:watcher="item")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import WeatherItem from './weather-item.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('watcher');

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
    this.fetchWeatherList({
      direction: 'ASC',
      limit: NaN,
      orderby: 'display_name',
      start: null,
    });
  },
  methods: {
    ...mapActions(['fetchWeatherList']),
  },
};
</script>


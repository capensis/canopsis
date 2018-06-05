<template lang="pug">
  div
    v-container(fluid)
      ul
        li(v-for="item in items")
          list-actions(:item="item")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import ListActions from '../components/list-actions.vue';

const { mapActions: alarmMapActions, mapGetters: alarmMapGetters } = createNamespacedHelpers('alarm');

export default {
  name: 'HelloWorld',
  components: {
    ListActions,
  },
  computed: {
    ...alarmMapGetters(['items']),
  },
  mounted() {
    this.fetchList({ params: { limit: 10 } });
  },
  methods: {
    ...alarmMapActions(['fetchList']),
  },
};
</script>

<style scoped>

</style>

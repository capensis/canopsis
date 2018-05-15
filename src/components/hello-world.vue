<template lang="pug">
  div
    h2 {{ msg }}
    v-container(fluid)
      ul
        li(v-for="item in items")
          list-actions(:item="item")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import ListActions from './list-actions.vue';

const { mapActions: alarmMapActions, mapGetters: alarmMapGetters } = createNamespacedHelpers('entities/alarm');

export default {
  name: 'HelloWorld',
  components: {
    ListActions,
  },
  props: {
    msg: {
      type: String,
      required: true,
    },
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

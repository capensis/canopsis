<template lang="pug">
  div
    h2 {{ msg }}
    v-container(fluid)
      list-actions
</template>

<script>
import ListActions from './list-actions.vue';

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
};
</script>

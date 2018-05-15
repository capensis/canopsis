<template lang="pug">
  div
    h2 {{ msg }}
    v-container(fluid)
      list-actions
      filters-selector
</template>

<script>
import ListActions from './list-actions.vue';
import FiltersSelector from './other/filters_list/Index.vue';

export default {
  name: 'HelloWorld',
  components: {
    ListActions,
    FiltersSelector,
  },
  props: {
    msg: {
      type: String,
      required: true,
    },
  },
};
</script>

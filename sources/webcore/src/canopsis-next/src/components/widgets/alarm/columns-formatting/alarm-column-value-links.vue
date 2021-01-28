<template lang="pug">
  v-list.links-list.transparent.pa-0(dense)
    category-links(:links="links", :label="label")
</template>

<script>
import CategoryLinks from './category-links/category-links.vue';

export default {
  components: { CategoryLinks },
  props: {
    links: {
      type: Array,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
  },
};
</script>

<template lang="pug">
  div
    entities-list(
    :widget="widget",
    )
</template>

<script>
import EntitiesList from '@/components/other/context/entities-list.vue';

export default {
  components: {
    EntitiesList,
  },
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
};
</script>


<template lang="pug">
  div
    entities-select(label="Dependencies", @update:entities="updateDependencies($event)")
</template>

<script>
import EntitiesSelect from '@/components/other/context-explorer/actions/createEntities/entitiesSelect.vue';

export default {
  components: { EntitiesSelect },
  data() {
    return {
      dependencies: [],
    };
  },
  methods: {
    updateDependencies(entities) {
      this.dependencies = entities.map(entity => entity._id);
    },
  },
};
</script>

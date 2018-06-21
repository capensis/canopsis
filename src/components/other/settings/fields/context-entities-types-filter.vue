<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.contextTypeOfEntities.title')}}
    v-container
      v-checkbox(
      v-for="(entitiesType, index) in entitiesTypes",
      v-bind="entitiesType",
      :key="`entities-type-${index}`",
      :input-value="value",
      @change="$emit('input', $event)",
      hide-details
      )
</template>

<script>
export default {
  props: {
    value: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      entitiesTypes: [
        { label: this.$t('settings.contextTypeOfEntities.fields.component'), value: 'component' },
        { label: this.$t('settings.contextTypeOfEntities.fields.connector'), value: 'connector' },
        { label: this.$t('settings.contextTypeOfEntities.fields.resource'), value: 'resource' },
        { label: this.$t('settings.contextTypeOfEntities.fields.watcher'), value: 'watcher' },
      ],
    };
  },
};
</script>

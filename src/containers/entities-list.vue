<template lang="pug">
  div
    entities-list(
    :widget="widget",
    :columns="$mq | mq(contextColumns)",
    @openSettings="$emit('openSettings', $event)"
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
  data() {
    return {
      settingsFields: [
        'title',
        'default-column-sort',
        'column-names',
        'context-entities-types-filter',
      ],
      contextColumns: {
        l: [
          {
            text: this.$t('tables.contextEntities.columns._id'),
            value: '_id',
          },
          {
            text: this.$t('tables.contextEntities.columns.type'),
            value: 'type',
          },
          {
            text: this.$t('tables.contextEntities.columns.name'),
            value: 'name',
          },
        ],
        m: [
          {
            text: this.$t('tables.contextEntities.columns.name'),
            value: 'name',
          },
        ],
      },
    };
  },
};
</script>


<template lang="pug">
  div
    entities-list(
    :properties="$mq | mq(contextProperties)",
    @openSettings="openSettings"
    )
</template>

<script>
import EntitiesList from '@/components/other/context/entities-list.vue';
import ContextSettingsFields from '@/components/other/settings/context-settings-fields.vue';

export default {
  components: {
    EntitiesList,
    ContextSettingsFields,
  },
  data() {
    return {
      settingsFields: [
        'title',
        'default-column-sort',
        'column-names',
        'context-entities-types-filter',
      ],
      contextProperties: {
        laptop: [
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
        mobile: [
          {
            text: this.$t('tables.contextEntities.columns.name'),
            value: 'name',
          },
        ],
        tablet: [
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


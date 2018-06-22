<template lang="pug">
  div
    settings-wrapper(v-model="isSettingsOpen", :title="$t('settings.titles.contextTableSettings')")
      context-settings-fields
    context-table(
    :contextProperties="$mq| mq(contextProperties)",
    @openSettings="openSettings"
    )
</template>

<script>
// COMPONENTS
import ContextTable from '@/components/other/context-explorer/context-table.vue';
import ContextSettingsFields from '@/components/other/settings/context-settings-fields.vue';
// MIXINS
import contextEntityMixin from '@/mixins/context';
import settingsMixin from '@/mixins/settings';
// OTHERS
import { PAGINATION_LIMIT } from '@/config';

export default {
  components: {
    ContextTable,
    ContextSettingsFields,
  },
  mixins: [contextEntityMixin, settingsMixin],
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
  mounted() {
    this.fetchContextEntities({
      params: {
        limit: PAGINATION_LIMIT,
      },
    });
  },
};
</script>

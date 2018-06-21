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
import ContextTable from '@/components/other/context-explorer/context-table.vue';
import ContextSettingsFields from '@/components/other/settings/context-settings-fields.vue';
import settingsMixin from '@/mixins/settings';

export default {
  components: {
    ContextTable,
    ContextSettingsFields,
  },
  mixins: [settingsMixin],
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
            text: 'name',
            value: 'name',
          },
          {
            text: 'type',
            value: 'type',
          },
        ],
        mobile: [
          {
            text: 'name',
            value: 'name',
          },
        ],
        tablet: [
          {
            text: 'name',
            value: 'name',
          },
        ],
      },
    };
  },
};
</script>

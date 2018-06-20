<template lang="pug">
  div
    context-settings(
    v-model="isSettingsOpen",
    :title="$t('settings.titles.contextTableSettings')",
    )
    context-table(
    :contextProperties="$mq| mq(contextProperties)",
    @openSettings="openSettings"
    )
</template>

<script>
import ContextTable from '@/components/other/context-explorer/context-table.vue';
import ContextSettings from '@/components/other/settings/context-settings.vue';
import settingsMixin from '@/mixins/settings';

export default {
  name: 'context',
  components: { ContextTable, ContextSettings },
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

<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      v-divider
      field-stats-groups(:groups.sync="settings.widget.parameters.groups")
    v-btn(@click="", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldTitle from '../partial/fields/title.vue';
import FieldStatsGroups from '../partial/fields/stats-groups.vue';

export default {
  components: {
    FieldTitle,
    FieldStatsGroups,
  },
  mixins: [widgetSettingsMixin],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: {
          title: widget.title,
          parameters: {
            groups: widget.groups || [],
          },
        },
      },
    };
  },
};
</script>


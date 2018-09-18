<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      v-divider
      field-duration(v-model="settings.widget.parameters.duration")
      v-divider
      field-date-select(:title="$t('settings.tstop')", v-model="settings.widget.parameters.tstop")
      v-divider
      field-stats-select
      v-divider
      field-filter-editor(v-model="settings.widget.parameters.mfilter")
      v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import { SIDE_BARS } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldTitle from '../partial/fields/title.vue';
import FieldDuration from '../partial/fields/duration.vue';
import FieldDateSelect from '../partial/fields/date-select.vue';
import FieldStatsSelect from '../partial/fields/stats-select.vue';
import FieldFilterEditor from '../partial/fields/filter-editor.vue';

export default {
  name: SIDE_BARS.statsTableSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldTitle,
    FieldDuration,
    FieldDateSelect,
    FieldStatsSelect,
    FieldFilterEditor,
  },
  mixins: [widgetSettingsMixin],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: {
          title: widget.title,
          parameters: {
            duration: widget.parameters.duration,
            tstop: widget.parameters.tstop,
            mfilter: widget.parameters.mfilter,
          },
        },
      },
    };
  },
};

</script>

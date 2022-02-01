<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-date-interval(
        v-model="settings.widget.parameters.dateInterval",
        :hidden-fields="['periodValue']"
      )
      v-divider
      field-filter-editor(
        v-model="settings.widget.parameters.mfilter",
        :hidden-fields="['title']",
        :entities-type="$constants.ENTITIES_TYPES.entity"
      )
      v-divider
      field-stats-selector(v-model="settings.widget.parameters.stats", required, with-trend, with-sorting)
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
            v-model="settings.widget.parameters.sort",
            :columns="defaultSortColumns",
            :columns-label="$t('settings.columnName')"
          )
          v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldDateInterval from '@/components/sidebars/settings/fields/stats/date-interval.vue';
import FieldStatsSelector from '@/components/sidebars/settings/fields/stats/stats-selector.vue';
import FieldFilterEditor from '@/components/sidebars/settings/fields/common/filter-editor.vue';
import FieldDefaultSortColumn from '@/components/sidebars/settings/fields/common/default-sort-column.vue';

export default {
  name: SIDE_BARS.statsTableSettings,
  components: {
    FieldTitle,
    FieldDateInterval,
    FieldStatsSelector,
    FieldFilterEditor,
    FieldDefaultSortColumn,
  },
  mixins: [widgetSettingsMixin],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: cloneDeep(widget),
      },
    };
  },
  computed: {
    defaultSortColumns() {
      return Object.keys(this.settings.widget.parameters.stats).map(key => ({ label: key, value: key }));
    },
  },
};

</script>

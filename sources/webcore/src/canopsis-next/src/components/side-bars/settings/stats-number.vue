<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-date-interval(v-model="settings.widget.parameters.dateInterval", :hiddenFields="['periodValue']")
      v-divider
      field-filter-editor(
        data-test="widgetFilterEditor",
        v-model="settings.widget.parameters.mfilter",
        :hiddenFields="['title']",
        :entitiesType="$constants.ENTITIES_TYPES.entity"
      )
      v-divider
      field-stat-selector(v-model="settings.widget.parameters.stat")
      v-divider
      v-list-group(data-test="advancedSettings")
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-stat-display-mode(v-model="settings.widget.parameters.displayMode")
          v-divider
          field-default-elements-per-page(v-model="settings.widget.parameters.limit")
          v-divider
          field-sort-order(v-model="settings.widget.parameters.sortOrder")
          v-divider
    v-btn.primary(data-test="statsNumberSubmitButton", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import widgetSettingsMixin from '@/mixins/widget/settings';
import { SIDE_BARS } from '@/constants';

import FieldTitle from './fields/common/title.vue';
import FieldDateInterval from './fields/stats/date-interval.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldStatSelector from './fields/stats/stat-selector.vue';
import FieldStatDisplayMode from './fields/stats/stat-display-mode.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldSortOrder from './fields/stats/sort-order.vue';

export default {
  name: SIDE_BARS.statsNumberSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldTitle,
    FieldDateInterval,
    FieldFilterEditor,
    FieldStatSelector,
    FieldStatDisplayMode,
    FieldDefaultElementsPerPage,
    FieldSortOrder,
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
};
</script>

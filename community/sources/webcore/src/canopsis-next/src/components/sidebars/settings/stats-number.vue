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
      field-stat-selector(v-model="settings.widget.parameters.stat")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-stat-display-mode(v-model="settings.widget.parameters.displayMode")
          v-divider
          field-default-elements-per-page(v-model="settings.widget.parameters.limit")
          v-divider
          field-sort-order(v-model="settings.widget.parameters.sortOrder")
          v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldDateInterval from '@/components/sidebars/settings/fields/stats/date-interval.vue';
import FieldFilterEditor from '@/components/sidebars/settings/fields/common/filter-editor.vue';
import FieldStatSelector from '@/components/sidebars/settings/fields/stats/stat-selector.vue';
import FieldStatDisplayMode from '@/components/sidebars/settings/fields/stats/stat-display-mode.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/settings/fields/common/default-elements-per-page.vue';
import FieldSortOrder from '@/components/sidebars/settings/fields/stats/sort-order.vue';

export default {
  name: SIDE_BARS.statsNumberSettings,
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

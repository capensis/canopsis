<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-filter-editor(
        data-test="widgetFilterEditor",
        v-model="settings.widget.parameters.mfilter",
        :hiddenFields="['title']",
        :entitiesType="$constants.ENTITIES_TYPES.entity"
      )
      v-divider
      field-text-editor(
        data-test="widgetTestTemplate",
        v-model="settings.widget.parameters.template",
        :title="$t('settings.templateEditor')"
      )
      v-divider
      v-list-group(v-if="edition === $constants.CANOPSIS_EDITION.cat", data-test="textWidgetStats")
        v-list-tile(slot="activator") {{ $t('settings.stats') }}
          .font-italic.caption.ml-1 ({{ $t('common.optional') }})
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-stats-selector(v-model="settings.widget.parameters.stats")
          v-divider
          field-date-interval(v-model="settings.widget.parameters.dateInterval")
          v-divider
    v-btn.primary(data-test="submitText", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS } from '@/constants';

import widgetSettingsMixin from '@/mixins/widget/settings';
import entitiesInfoMixin from '@/mixins/entities/info';

import FieldTitle from './fields/common/title.vue';
import FieldDateInterval from './fields/stats/date-interval.vue';
import FieldStatsSelector from './fields/stats/stats-selector.vue';
import FieldFilterEditor from './fields/common/filter-editor.vue';
import FieldTextEditor from './fields/common/text-editor.vue';

export default {
  name: SIDE_BARS.textSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldTitle,
    FieldDateInterval,
    FieldStatsSelector,
    FieldFilterEditor,
    FieldTextEditor,
  },
  mixins: [widgetSettingsMixin, entitiesInfoMixin],
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

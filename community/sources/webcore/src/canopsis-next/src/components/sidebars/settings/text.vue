<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title", :title="$t('common.title')")
      v-divider
      field-filter-editor(
        v-model="form.parameters.mfilter",
        :hidden-fields="['title']",
        :entities-type="$constants.ENTITIES_TYPES.entity"
      )
      v-divider
      field-text-editor(
        v-model="form.parameters.template",
        :title="$t('settings.templateEditor')"
      )
      v-divider
      v-list-group(v-if="isCatVersion")
        template(#activator="")
          v-list-tile {{ $t('settings.stats') }}
            div.font-italic.caption.ml-1 ({{ $t('common.optional') }})
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-stats-selector(v-model="form.parameters.stats")
          v-divider
          field-date-interval(v-model="form.parameters.dateInterval")
          v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldDateInterval from '@/components/sidebars/settings/fields/stats/date-interval.vue';
import FieldStatsSelector from '@/components/sidebars/settings/fields/stats/stats-selector.vue';
import FieldFilterEditor from '@/components/sidebars/settings/fields/common/filter-editor.vue';
import FieldTextEditor from '@/components/sidebars/settings/fields/common/text-editor.vue';

export default {
  name: SIDE_BARS.textSettings,
  components: {
    FieldTitle,
    FieldDateInterval,
    FieldStatsSelector,
    FieldFilterEditor,
    FieldTextEditor,
  },
  mixins: [widgetSettingsMixin, entitiesInfoMixin],
};
</script>

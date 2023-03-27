<template lang="pug">
  v-card.subheading
    v-card-text
      v-layout(row, wrap)
        v-flex(xs1)
          v-avatar.white--text.mr-2(color="primary", size="32") {{ actionNumber }}
          c-expand-btn(v-model="expanded")
        v-flex(xs11)
          scenario-info-item.scenario-info-type.px-2(
            :label="$t('common.type')",
            :value="action.type",
            hide-icon
          )
          v-expand-transition(mode="out-in")
            v-layout.px-2(v-if="expanded", column)
              v-tabs(
                v-model="activeTab",
                slider-color="primary",
                color="transparent",
                centered,
                fixed-tabs
              )
                v-tab {{ $t('common.general') }}
                v-tab {{ $t('scenario.tabs.pattern') }}
              v-divider
              v-tabs-items.pt-2(v-model="activeTab")
                v-tab-item
                  scenario-action-card-general-tab(:action="action")
                v-tab-item
                  c-patterns-field(
                    :value="patterns",
                    with-alarm,
                    with-entity,
                    readonly
                  )
</template>

<script>
import { OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/forms/filter';

import ScenarioInfoItem from './scenario-info-item.vue';
import ScenarioActionCardGeneralTab from './scenario-action-card-general-tab.vue';

export default {
  components: { ScenarioInfoItem, ScenarioActionCardGeneralTab },
  props: {
    action: {
      type: Object,
      required: true,
    },
    actionNumber: {
      type: [Number, String],
      default: 0,
    },
  },
  data() {
    return {
      activeTab: 0,
      expanded: false,
    };
  },
  computed: {
    patterns() {
      return filterPatternsToForm(
        this.action,
        [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
        [OLD_PATTERNS_FIELDS.alarm, OLD_PATTERNS_FIELDS.entity],
      );
    },
  },
};
</script>

<style lang="scss" scoped>
.scenario-info-type ::v-deep .v-list__tile {
  height: 30px;
}
</style>

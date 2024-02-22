<template>
  <v-card class="text-subtitle-1">
    <v-card-text>
      <v-layout wrap>
        <v-flex xs1>
          <v-avatar
            class="white--text mr-2"
            color="primary"
            size="32"
          >
            {{ actionNumber }}
          </v-avatar>
          <c-expand-btn v-model="expanded" />
        </v-flex>
        <v-flex xs11>
          <scenario-info-item
            :label="$t('common.type')"
            :value="action.type"
            class="scenario-info-type px-2"
            hide-icon
          />
          <v-expand-transition mode="out-in">
            <v-layout
              v-if="expanded"
              class="px-2"
              column
            >
              <v-tabs
                v-model="activeTab"
                slider-color="primary"
                centered
                fixed-tabs
              >
                <v-tab>{{ $t('common.general') }}</v-tab>
                <v-tab>{{ $t('scenario.tabs.pattern') }}</v-tab>
              </v-tabs>
              <v-divider />
              <v-tabs-items
                v-model="activeTab"
                class="pt-2"
              >
                <v-tab-item>
                  <scenario-action-card-general-tab :action="action" />
                </v-tab-item>
                <v-tab-item>
                  <c-patterns-field
                    :value="patterns"
                    with-alarm
                    with-entity
                    readonly
                  />
                </v-tab-item>
              </v-tabs-items>
            </v-layout>
          </v-expand-transition>
        </v-flex>
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

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
      return filterPatternsToForm(this.action, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]);
    },
  },
};
</script>

<style lang="scss" scoped>
.scenario-info-type ::v-deep .v-list-item {
  height: 30px;
}
</style>

<template>
  <v-card class="text-subtitle-1">
    <v-card-text>
      <v-row class="flex-grow-0" no-gutters>
        <v-col class="pt-2" cols="1">
          <v-avatar class="text-white mr-2" color="primary" size="32">
            {{ actionNumber }}
          </v-avatar>
          <c-expand-btn v-model="expanded" />
        </v-col>
        <v-col cols="11">
          <scenario-info-item
            :label="$t('common.type')"
            :value="action.type"
            class="scenario-info-type px-2"
            hide-icon
          />
          <v-expand-transition mode="out-in">
            <c-layout v-if="expanded" class="px-2" column>
              <v-tabs
                v-model="activeTab"
                color="primary"
                align-tabs="center"
                fixed-tabs
              >
                <v-tab :value="0">
                  {{ $t('common.general') }}
                </v-tab>
                <v-tab :value="1">
                  {{ $t('scenario.tabs.pattern') }}
                </v-tab>
              </v-tabs>
              <v-divider />
              <v-window v-model="activeTab" class="pt-2">
                <v-window-item :value="0">
                  <scenario-action-card-general-tab :action="action" />
                </v-window-item>
                <v-window-item :value="1">
                  <c-patterns-field
                    :value="patterns"
                    with-alarm
                    with-entity
                    readonly
                  />
                </v-window-item>
              </v-window>
            </c-layout>
          </v-expand-transition>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import ScenarioActionCardGeneralTab from './scenario-action-card-general-tab.vue';
import ScenarioInfoItem from './scenario-info-item.vue';

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
.scenario-info-type :deep(.v-list-item) {
  min-height: 30px;
}
</style>

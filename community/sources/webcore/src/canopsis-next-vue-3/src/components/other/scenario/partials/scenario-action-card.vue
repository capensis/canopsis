<template>
  <v-card class="text-subtitle-1">
    <v-card-text>
      <c-layout wrap>
        <div class="scenario-action-card__number pt-2">
          <v-avatar class="text-white mr-2" color="primary" size="32">
            {{ actionNumber }}
          </v-avatar>
          <c-expand-btn v-model="expanded" />
        </div>
        <div class="scenario-action-card__content">
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
fixed-tabs>
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
                  <pre class="scenario-action-card__patterns"
                    >{{ patterns }}
</pre
                  >
                </v-window-item>
              </v-window>
            </c-layout>
          </v-expand-transition>
        </div>
      </c-layout>
    </v-card-text>
  </v-card>
</template>

<script>
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
      return JSON.stringify(
        {
          alarm_pattern: this.action.alarm_pattern,
          entity_pattern: this.action.entity_pattern,
        },
        null,
        2,
      );
    },
  },
};
</script>

<style lang="scss" scoped>
.scenario-action-card__number {
  width: 8.333333%;
}

.scenario-action-card__content {
  width: 91.666667%;
}

.scenario-action-card__patterns {
  margin: 0;
  white-space: pre-wrap;
}

.scenario-info-type :deep(.v-list-item) {
  min-height: 30px;
}
</style>

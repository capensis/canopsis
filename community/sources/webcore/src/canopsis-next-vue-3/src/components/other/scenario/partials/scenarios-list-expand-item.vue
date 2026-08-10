<template>
  <c-layout class="secondary lighten-2 py-3">
    <div class="scenarios-list-expand-item__content">
      <v-card>
        <v-card-text>
          <c-layout justify="center" align="center" wrap>
            <div v-if="scenario.author" class="w-100">
              <scenario-info-item
                :label="$t('common.author')"
                :value="scenario.author.display_name"
                icon="person"
              />
            </div>
            <div class="w-100">
              <scenario-info-item
                :label="$t('common.trigger', 2)"
                :value="preparedTriggers"
                icon="bolt"
              />
            </div>
            <div v-if="hasDisableDuringPeriods" class="w-100">
              <scenario-info-item
                :label="$t('common.disableDuringPeriods')"
                :value="scenario.disable_during_periods.join(', ')"
                icon="highlight_off"
              />
            </div>
            <div v-for="(action, index) in scenario.actions" :key="index" class="mt-2 w-100">
              <scenario-action-card :action="action" :action-number="index + 1" />
            </div>
          </c-layout>
        </v-card-text>
      </v-card>
    </div>
  </c-layout>
</template>

<script>
import { isEmpty } from 'lodash';

import ScenarioActionCard from './scenario-action-card.vue';
import ScenarioInfoItem from './scenario-info-item.vue';

export default {
  components: { ScenarioInfoItem, ScenarioActionCard },

  props: {
    scenario: {
      type: Object,
      required: true,
    },
  },
  computed: {
    preparedTriggers() {
      return this.scenario.triggers
        .map(({ type, ...rest }) => {
          if (isEmpty(rest)) {
            return type;
          }

          const preparedRest = Object.entries(rest)
            .map(([key, value]) => `${key}: ${value}`)
            .join(', ');

          return `${type} (${preparedRest})`;
        })
        .join(', ');
    },

    hasDisableDuringPeriods() {
      return this.scenario.disable_during_periods?.length;
    },
  },
};
</script>

<style lang="scss" scoped>
.scenarios-list-expand-item__content {
  width: 66.666667%;
  margin-left: auto;
  margin-right: auto;
}
</style>

<template>
  <v-layout class="secondary lighten-2 py-3">
    <v-flex
      xs12
      sm8
      offset-sm2
    >
      <v-card>
        <v-card-text>
          <v-layout
            wrap
            justify-center
            align-center
          >
            <v-flex
              v-if="scenario.author"
              xs12
            >
              <scenario-info-item
                :label="$t('common.author')"
                :value="scenario.author.display_name"
                icon="person"
              />
            </v-flex>
            <v-flex xs12>
              <scenario-info-item
                :label="$tc('common.trigger', 2)"
                :value="preparedTriggers"
                icon="bolt"
              />
            </v-flex>
            <v-flex
              v-if="hasDisableDuringPeriods"
              xs12
            >
              <scenario-info-item
                :label="$t('common.disableDuringPeriods')"
                :value="scenario.disable_during_periods.join(', ')"
                icon="highlight_off"
              />
            </v-flex>
            <v-flex
              class="mt-2"
              v-for="(action, index) in scenario.actions"
              :key="index"
              xs12
            >
              <scenario-action-card
                :action="action"
                :action-number="index + 1"
              />
            </v-flex>
          </v-layout>
        </v-card-text>
      </v-card>
    </v-flex>
  </v-layout>
</template>

<script>
import { isEmpty } from 'lodash';

import ScenarioInfoItem from './scenario-info-item.vue';
import ScenarioActionCard from './scenario-action-card.vue';

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
      return this.scenario.triggers.map(({ type, ...rest }) => {
        if (isEmpty(rest)) {
          return type;
        }

        const preparedRest = Object.entries(rest)
          .map(([key, value]) => `${key}: ${value}`)
          .join(', ');

        return `${type} (${preparedRest})`;
      }).join(', ');
    },

    hasDisableDuringPeriods() {
      return this.scenario.disable_during_periods?.length;
    },
  },
};
</script>

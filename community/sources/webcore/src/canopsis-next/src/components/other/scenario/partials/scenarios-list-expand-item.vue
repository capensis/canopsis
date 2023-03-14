<template lang="pug">
  v-layout.secondary.lighten-2.py-3(row)
    v-flex(xs12, sm8, offset-sm2)
      v-card
        v-card-text
          v-layout(wrap, justify-center, align-center)
            v-flex(xs12)
              scenario-info-item(
                :label="$t('common.author')",
                :value="scenario.author.name",
                icon="person"
              )
            v-flex(xs12)
              scenario-info-item(
                :label="$tc('common.trigger', 2)",
                :value="scenario.triggers.join(', ')",
                icon="bolt"
              )
            v-flex(v-if="hasDisableDuringPeriods", xs12)
              scenario-info-item(
                :label="$t('common.disableDuringPeriods')",
                :value="scenario.disable_during_periods.join(', ')",
                icon="highlight_off"
              )
            v-flex.mt-2(v-for="(action, index) in scenario.actions", :key="index", xs12)
              scenario-action-card(
                :action="action",
                :action-number="index + 1"
              )
</template>

<script>
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
    hasDisableDuringPeriods() {
      return this.scenario.disable_during_periods?.length;
    },
  },
};
</script>

<template lang="pug">
  div
    v-container
      v-container.pt-0(fluid, data-test="statTypeLayout")
        v-select(
          v-field="form.stat",
          :items="statsTypes",
          hide-details,
          return-object
        )
        v-text-field(
          v-field="form.title",
          v-validate="'required'",
          :placeholder="$t('common.title')",
          :error-messages="errors.collect('title')",
          name="title",
          data-test="statTitle"
        )
        v-card(v-if="form.stat.options.length || withTrend", color="secondary white--text", dark)
          v-card-title(data-test="statParameters") {{ $t('common.parameters') }}
          v-card-text
            v-switch(
              v-if="withTrend",
              v-field="form.trend",
              :label="$t('common.trend')",
              color="primary",
              data-test="statTrend",
              hide-details
            )
            template(v-for="option in form.stat.options")
              v-switch(
                v-if="option === $constants.STATS_OPTIONS.recursive",
                v-field="form.parameters.recursive",
                :label="$t('common.recursive')",
                color="primary",
                data-test="statRecursive",
                hide-details
              )
              v-layout(v-else-if="option === $constants.STATS_OPTIONS.states", data-test="statStates", row)
                v-select(
                  v-field="form.parameters.states",
                  :placeholder="$t('common.states')",
                  :items="stateTypes",
                  multiple,
                  chips,
                  hide-details
                )
              v-combobox(
                v-else-if="option === $constants.STATS_OPTIONS.authors",
                v-field="form.parameters.authors",
                :placeholder="$t('common.authors')",
                data-test="statAuthors",
                hide-details,
                chips,
                multiple
              )
              v-text-field(
                v-else-if="option === $constants.STATS_OPTIONS.sla",
                v-field="form.parameters.sla",
                v-validate="{ required: true, regex: /^(<|>|<=|>=)\\s*\\d+$/ }",
                :placeholder="$t('common.sla')",
                :error-messages="errors.collect('sla')",
                :hide-details="!errors.has('sla')",
                name="sla",
                data-test="statSla"
              )
                v-tooltip(slot="append", left)
                  v-icon(slot="activator", dark) help_outline
                  span {{ $t('modals.addStat.slaTooltip') }}
</template>

<script>
import { STATS_TYPES, ENTITIES_STATES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({ parameters: {} }),
    },
    withTrend: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    /**
     * Get stats different types from constant, and return an object with stat's value and stat's translated title
     */
    statsTypes() {
      return Object.values(STATS_TYPES)
        .map(item => ({ value: item.value, text: this.$t(`stats.types.${item.value}`), options: item.options }));
    },
    stateTypes() {
      return Object.keys(ENTITIES_STATES)
        .map(item => ({ value: ENTITIES_STATES[item], text: item }));
    },
  },
};
</script>

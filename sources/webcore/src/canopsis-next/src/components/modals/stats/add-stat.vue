<template lang="pug">
  v-card(data-test="addStatModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-form
      v-container
        v-container.pt-0(fluid, data-test="statTypeLayout")
          v-select(
            v-model="form.stat",
            hide-details,
            :items="statsTypes",
            return-object
          )
          v-text-field(
            data-test="statTitle",
            :placeholder="$t('common.title')",
            v-model="form.title",
            :error-messages="errors.collect('title')",
            v-validate="'required'",
            data-vv-name="title"
          )
          v-card(v-if="form.stat.options.length || config.withTrend", color="secondary white--text", dark)
            v-card-title(data-test="statParameters") {{ $t('common.parameters') }}
            v-card-text
              v-switch(
                data-test="statTrend",
                v-if="config.withTrend",
                :label="$t('common.trend')",
                v-model="form.trend",
                hide-details,
                color="primary"
              )
              template(v-for="option in form.stat.options")
                v-switch(
                  data-test="statRecursive",
                  v-if="option === $constants.STATS_OPTIONS.recursive",
                  :label="$t('common.recursive')",
                  v-model="form.parameters.recursive",
                  hide-details,
                  color="primary"
                )
                v-layout(v-else-if="option === $constants.STATS_OPTIONS.states", data-test="statStates", row)
                  v-select(
                    :placeholder="$t('common.states')",
                    :items="stateTypes",
                    v-model="form.parameters.states",
                    multiple,
                    chips,
                    hide-details
                  )
                v-combobox(
                  data-test="statAuthors",
                  v-else-if="option === $constants.STATS_OPTIONS.authors",
                  :placeholder="$t('common.authors')",
                  v-model="form.parameters.authors",
                  hide-details,
                  chips,
                  multiple
                )
                v-text-field(
                  data-test="statSla",
                  v-else-if="option === $constants.STATS_OPTIONS.sla",
                  v-model="form.parameters.sla",
                  v-validate="{ required: true, regex: /^(<|>|<=|>=)\\s*\\d+$/ }",
                  :placeholder="$t('common.sla')",
                  :error-messages="errors.collect('sla')",
                  :hide-details="!errors.has('sla')",
                  name="sla"
                )
                  v-tooltip(slot="append", left)
                    v-icon(slot="activator", dark) help_outline
                    span {{ $t('modals.addStat.slaTooltip') }}
        v-divider
        v-layout.py-1(justify-end)
          v-btn(data-test="addStatCancelButton", @click="hideModal", depressed, flat) {{ $t('common.cancel') }}
          v-btn.primary(data-test="addStatSubmitButton", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, pick } from 'lodash';

import { MODALS, STATS_TYPES, ENTITIES_STATES } from '@/constants';

import { setIn } from '@/helpers/immutable';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.addStat,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        stat: STATS_TYPES.alarmsCreated,
        title: '',
        trend: true,
        parameters: {
          recursive: true,
        },
      },
    };
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
  mounted() {
    if (!this.config.withTrend) {
      this.form.trend = false;
    }

    if (this.config.stat) {
      const selectedStat = Object.values(STATS_TYPES)
        .find(stat => stat.value === this.config.stat.stat.value) || STATS_TYPES.alarmsCreated;

      this.form = { ...cloneDeep(this.config.stat), stat: cloneDeep(selectedStat), title: this.config.statTitle };
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          const preparedForm = setIn(this.form, 'parameters', parameters => pick(parameters, this.form.stat.options));

          await this.config.action(preparedForm);
        }

        this.hideModal();
      }
    },
  },
};
</script>

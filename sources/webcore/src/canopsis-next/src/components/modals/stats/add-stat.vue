<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t(config.title) }}
    v-form
      v-container
        v-container.pt-0(fluid)
          v-select(
          v-model="form.stat",
          hide-details,
          :items="statsTypes",
          return-object,
          )
          v-text-field(
          :placeholder="$t('common.title')",
          v-model="form.title",
          :error-messages="errors.collect('title')",
          v-validate="'required'",
          data-vv-name="title",
          )
          v-switch(
          v-if="config.withTrend",
          :label="$t('common.trend')",
          v-model="form.trend",
          hide-details
          )
          v-card(v-if="form.stat.options.length", color="secondary white--text", dark)
            v-card-title {{ $t('common.parameters') }}
            v-card-text
              template(v-for="option in form.stat.options")
                v-switch(
                v-if="option === 'recursive'"
                :label="$t('common.recursive')",
                v-model="form.parameters.recursive",
                hide-details,
                color="primary"
                )
                v-select(
                v-if="option === 'states'"
                :placeholder="$t('common.states')",
                :items="stateTypes",
                v-model="form.parameters.states",
                multiple,
                chips,
                hide-details
                )
                v-combobox(
                v-if="option === 'authors'"
                :placeholder="$t('common.authors')",
                v-model="form.parameters.authors",
                hide-details,
                chips,
                multiple
                )
                v-text-field(
                v-if="option === 'sla'",
                :placeholder="$t('common.sla')",
                v-model="form.parameters.sla",
                hide-details
                )
        v-alert(:value="error", type="error") {{ error }}
        v-divider
        v-layout.py-1(justify-end)
          v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
          v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, STATS_TYPES, STATS_OPTIONS, ENTITIES_STATES } from '@/constants';

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
      error: '',
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
      let isFormValid = await this.$validator.validateAll();

      if (this.form.stat.options.find(option => option === STATS_OPTIONS.sla) && !this.form.parameters.sla) {
        isFormValid = false;
        this.error = this.$t('modals.addStat.slaRequired');
      }

      if (isFormValid && this.config.action) {
        await this.config.action(this.form);
        this.hideModal();
      }
    },
  },
};
</script>

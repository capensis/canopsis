<template lang="pug">
  v-card(data-test="addStatModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-card-text
      add-stat-form(v-model="form", :withTrend="config.withTrend")
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(data-test="addStatCancelButton", @click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(data-test="addStatSubmitButton", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, pick } from 'lodash';

import { MODALS, STATS_TYPES } from '@/constants';

import { setIn } from '@/helpers/immutable';

import modalInnerMixin from '@/mixins/modal/inner';

import AddStatForm from '@/components/other/stats/add-stat-form.vue';

export default {
  name: MODALS.addStat,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AddStatForm },
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

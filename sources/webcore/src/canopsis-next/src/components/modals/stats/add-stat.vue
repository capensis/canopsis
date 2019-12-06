<template lang="pug">
  v-form(data-test="addStatModal", @submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        add-stat-form(v-model="form", :withTrend="config.withTrend")
      template(slot="actions")
        v-btn(
          data-test="addStatCancelButton",
          depressed,
          flat,
          @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          data-test="addStatSubmitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, pick } from 'lodash';

import { MODALS, STATS_TYPES } from '@/constants';

import { setField } from '@/helpers/immutable';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import AddStatForm from '@/components/other/stats/add-stat-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.addStat,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AddStatForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
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
          const preparedForm = setField(this.form, 'parameters', parameters => pick(parameters, this.form.stat.options));

          await this.config.action(preparedForm);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

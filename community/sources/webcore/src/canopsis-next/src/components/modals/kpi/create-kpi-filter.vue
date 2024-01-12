<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <kpi-filter-form v-model="form" />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :disabled="isDisabled"
          :loading="submitting"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, PATTERNS_FIELDS, VALIDATION_DELAY } from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import KpiFilterForm from '@/components/other/kpi/filters/form/kpi-filter-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createKpiFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { KpiFilterForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { filter = {} } = this.modal.config;

    return {
      form: {
        name: filter.name ?? '',
        patterns: filterPatternsToForm(filter, [PATTERNS_FIELDS.entity]),
      },
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createFilter.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action({
            name: this.form.name,
            ...formFilterToPatterns(this.form.patterns, [PATTERNS_FIELDS.entity]),
          });
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

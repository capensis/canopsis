<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('common.filters') }}
      template(slot="text")
        remediation-instructions-filter-form(
          v-model="form",
          :filters="config.anotherFilters"
        )
      template(slot="actions")
        v-btn(
          :disabled="submitting",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import RemediationInstructionsFilterForm
  from '@/components/other/remediation/instructions-filter/remediation-instructions-filter-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRemediationInstructionsFilter,
  $_veeValidate: {
    validator: 'new',
  },
  components: { RemediationInstructionsFilterForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixin(),
  ],
  data() {
    const defaultForm = {
      with: true,
      all: false,
      instructions: [],
    };

    const { filter } = this.modal.config;

    return {
      form: filter ? cloneDeep(filter) : defaultForm,
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

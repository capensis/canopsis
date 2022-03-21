<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        c-patterns-field(v-model="form", alarm, entity, pbehavior)
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import RemediationPatternsPbehaviorTypesForm
  from '@/components/other/remediation/patterns/remediation-patterns-pbehavior-types-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.remediationPatterns,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    RemediationPatternsPbehaviorTypesForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: this.modal.config.patterns ? cloneDeep(this.modal.config.patterns) : {},
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.patterns.title');
    },
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

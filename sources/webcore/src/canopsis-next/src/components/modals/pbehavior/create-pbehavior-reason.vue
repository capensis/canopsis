<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.createPbehaviorReason.title') }}
      template(slot="text")
        create-pbehavior-reason-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { pbehaviorReasonToForm } from '@/helpers/forms/reason-pbehavior';

import validationErrorsMixin from '@/mixins/form/validation-errors';

import CreatePbehaviorReasonForm from '@/components/other/pbehavior/reasons/form/create-pbehavior-reason-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorReason,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    CreatePbehaviorReasonForm,
    ModalWrapper,
  },
  mixins: [
    validationErrorsMixin(),
  ],
  data() {
    return {
      form: pbehaviorReasonToForm(this.modal.config.pbehaviorReason),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(this.form);
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>

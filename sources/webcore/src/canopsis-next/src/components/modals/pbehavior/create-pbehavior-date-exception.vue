<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createPbehaviorDateException.title') }}
      template(slot="text")
        pbehavior-date-exception-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { formToPbehaviorDateException, pbehaviorDateExceptionToForm } from '@/helpers/forms/dates-exceptions-pbehavior';

import modalInnerMixin from '@/mixins/modal/inner';
import validationErrorsMixin from '@/mixins/form/validation-errors';

import PbehaviorDateExceptionForm from '@/components/other/pbehavior/dates-exceptions/form/pbehavior-date-exception-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorDateException,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PbehaviorDateExceptionForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    validationErrorsMixin(),
  ],
  data() {
    return {
      form: pbehaviorDateExceptionToForm(this.modal.config.pbehaviorDateException),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToPbehaviorDateException(this.form));
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

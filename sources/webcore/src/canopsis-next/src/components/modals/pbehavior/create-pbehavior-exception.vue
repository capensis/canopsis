<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.createPbehaviorException.title') }}
      template(slot="text")
        pbehavior-exception-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { formToPbehaviorException, pbehaviorExceptionToForm } from '@/helpers/forms/exceptions-pbehavior';

import validationErrorsMixin from '@/mixins/form/validation-errors';

import PbehaviorExceptionForm from '@/components/other/pbehavior/exceptions/form/pbehavior-exception-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorException,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PbehaviorExceptionForm,
    ModalWrapper,
  },
  mixins: [
    validationErrorsMixin(),
  ],
  data() {
    return {
      form: pbehaviorExceptionToForm(this.modal.config.pbehaviorException),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToPbehaviorException(this.form));
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

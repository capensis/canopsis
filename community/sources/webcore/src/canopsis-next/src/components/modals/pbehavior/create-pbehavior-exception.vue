<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createPbehaviorException.title') }}
      template(#text="")
        pbehavior-exception-form(v-model="form", with-type)
      template(#actions="")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { formToPbehaviorException, pbehaviorExceptionToForm } from '@/helpers/forms/exceptions-pbehavior';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import PbehaviorExceptionForm from '@/components/other/pbehavior/exceptions/form/pbehavior-exception-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorException,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    PbehaviorExceptionForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
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
        if (this.config.action) {
          await this.config.action(formToPbehaviorException(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

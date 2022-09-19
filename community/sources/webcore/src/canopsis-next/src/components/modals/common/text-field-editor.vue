<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ config.title }}
      template(#text="")
        v-text-field(
          v-model="text",
          v-validate="field.validationRules",
          v-bind="fieldProps",
          :error-messages="errors.collect(field.name)"
        )
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          type="submit",
          :disabled="isDisabled",
          :loading="submitting"
        ) {{ $t('common.submit') }}
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { validationErrorsMixinCreator } from '@/mixins/form/validation-errors';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.textFieldEditor,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator({ field: 'text' }),
    validationErrorsMixinCreator(),
  ],
  data() {
    const field = this.modal.config.field ?? {};

    return {
      text: field.value ?? '',
    };
  },
  computed: {
    field() {
      return this.config.field ?? { name: 'text', label: 'Text' };
    },

    fieldProps() {
      return omit(this.field, ['validationRules', 'value']);
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(this.text);
          }

          this.$modals.hide();
        } catch (err) {
          const { name } = this.field;

          this.setFormErrors({ [name]: err[name] || err.error || err.message });
        }
      }
    },
  },
};
</script>

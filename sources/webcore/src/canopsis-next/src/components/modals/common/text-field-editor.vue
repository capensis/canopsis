<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(data-test="textFieldEditorModal")
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        v-text-field(
          data-test="textField",
          v-model="text",
          v-validate="field.validationRules",
          :name="field.name",
          :label="field.label",
          :error-messages="errors.collect(field.name)"
        )
      template(slot="actions")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          type="submit",
          :disabled="isDisabled",
          :loading="submitting",
          data-test="submitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.textFieldEditor,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const field = this.modal.config.field || {};

    return {
      text: field.value || '',
    };
  },
  computed: {
    field() {
      return this.config.field || { name: 'text', label: 'Text' };
    },
  },
  methods: {
    async submit() {
      const formIsValid = await this.$validator.validateAll();

      if (formIsValid) {
        if (this.config.action) {
          await this.config.action(this.text);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>


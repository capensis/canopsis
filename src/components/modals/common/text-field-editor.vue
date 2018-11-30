<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-card-text
      v-text-field(
      v-model="text",
      v-validate="field.validationRules",
      :name="field.name",
      :label="field.label",
      :error-messages="errors.collect(field.name)"
      )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/modal-inner';

export default {
  name: MODALS.textFieldEditor,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
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

        this.hideModal();
      }
    },
  },
};
</script>


<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createDynamicInfoInformation.create.title') }}
    v-card-text
      v-form
        v-text-field(
          v-validate="'required'",
          v-model="form.name",
          :label="$t('modals.createDynamicInfoInformation.fields.name')",
          :error-messages="errors.collect('name')",
          name="name"
        )
        v-text-field(
          v-validate="'required'",
          v-model="form.value",
          :label="$t('modals.createDynamicInfoInformation.fields.value')",
          :error-messages="errors.collect('value')",
          name="value"
        )
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

/**
 * Modal to create dynamic info's information
 */
export default {
  name: MODALS.createDynamicInfoInformation,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        name: '',
        value: '',
      },
    };
  },
  mounted() {
    if (this.config.info) {
      this.form = { ...this.config.info };
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

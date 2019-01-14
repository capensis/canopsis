<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-card-text
      v-form
        v-text-field(
        :label="$t('common.name')",
        v-model="form.name",
        v-validate="'required|unique-name'",
        data-vv-name="name",
        :error-messages="errors.collect('name')"
        )
        v-text-field(
        :label="$t('common.description')",
        v-model="form.description",
        v-validate="'required'",
        data-vv-name="description",
        :error-messages="errors.collect('description')"
        )
        v-textarea(
        :label="$t('common.value')",
        v-model="form.value",
        v-validate,
        data-vv-rule="'required'",
        data-vv-name="value",
        :error-messages="errors.collect('value')"
        )
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn(@click="submit", color="primary") {{ $t('common.add') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.addEntityInfo,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        name: '',
        description: '',
        value: '',
      },
    };
  },
  created() {
    this.createUniqueValidationRule();
  },
  mounted() {
    if (this.config.editingInfo) {
      this.form = { ...this.config.editingInfo };
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(this.form);
        this.hideModal();
      }
    },

    createUniqueValidationRule() {
      this.$validator.extend('unique-name', {
        getMessage: () => this.$t('validator.unique'),
        validate: (value) => {
          if (this.config.editingInfo && this.config.editingInfo.name === value) {
            return true;
          }

          return !this.config.infos[value];
        },
      });
    },
  },
};
</script>

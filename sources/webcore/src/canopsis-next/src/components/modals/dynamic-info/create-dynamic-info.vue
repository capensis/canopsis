<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createDynamicInfo.create.title') }}
    v-card-text
      v-form
        dynamic-info-form(v-model="form")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import { dynamicInfoToForm, formToDynamicInfo } from '@/helpers/forms/dynamic-info';

import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';

/**
 * Modal to create dynamic information
 */
export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DynamicInfoForm },
  mixins: [modalInnerMixin],
  data() {
    const { dynamicInfo = {} } = this.modal.config;

    return {
      form: dynamicInfoToForm(dynamicInfo),
    };
  },
  computed: {
    patterns() {
      return this.form.patterns;
    },
  },
  watch: {
    patterns() {
      this.$validator.validate('patterns');
    },
  },
  methods: {
    async submit() {
      try {
        const isValid = await this.$validator.validateAll();

        if (isValid) {
          const preparedData = formToDynamicInfo(this.form);

          if (this.config.action) {
            await this.config.action(preparedData);
          }

          this.$modals.hide();
        }
      } catch (err) {
        this.$popups.error({ text: err.description });
      }
    },
  },
};
</script>

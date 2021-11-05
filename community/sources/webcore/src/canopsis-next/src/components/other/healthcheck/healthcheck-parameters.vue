<template lang="pug">
  v-layout.my-2(v-if="!form", justify-center)
    v-progress-circular(indeterminate, color="primary")
  v-flex(v-else, xs10, offset-xs1, md8, offset-md2, lg6, offset-lg3)
    v-form(@submit.prevent="submit")
      healthcheck-form(v-model="form")
      v-layout.mt-3(row, justify-end)
        v-btn.primary.mr-0(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { healthcheckParametersToForm } from '@/helpers/forms/healthcheck';

import { submittableMixin } from '@/mixins/submittable';
import { validationErrorsMixin } from '@/mixins/form/validation-errors';
import { entitiesHealthcheckParametersMixin } from '@/mixins/entities/healthcheck';

import HealthcheckForm from '@/components/other/healthcheck/form/healthcheck-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { HealthcheckForm },
  mixins: [
    submittableMixin(),
    validationErrorsMixin(),
    entitiesHealthcheckParametersMixin,
  ],
  data() {
    return {
      form: null,
    };
  },
  async mounted() {
    const healthcheckParameters = await this.fetchHealthcheckParametersWithoutStore();

    this.form = healthcheckParametersToForm(healthcheckParameters);
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          await this.updateHealthcheckParameters({ data: this.form });

          this.$popups.success({ text: this.$t('success.default') });
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>
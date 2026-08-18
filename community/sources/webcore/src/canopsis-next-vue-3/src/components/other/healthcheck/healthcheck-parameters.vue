<template>
  <div class="healthcheck-parameters">
    <div v-if="!form" class="my-2 d-flex justify-center">
      <v-progress-circular color="primary" indeterminate />
    </div>
    <div v-else class="healthcheck-parameters__form py-2">
      <v-form @submit.prevent="submit">
        <healthcheck-form :form="form" :errors="errors" @input="form = $event" />
        <div class="mt-3 d-flex justify-end">
          <v-btn
            :disabled="isDisabled"
            :loading="submitting"
            color="primary"
            type="submit"
          >
            {{ $t('common.submit') }}
          </v-btn>
        </div>
      </v-form>
    </div>
  </div>
</template>

<script>
import { HEALTHCHECK_ENGINES_NAMES } from '@/constants';

import { useHealthcheckParametersStore } from '@/stores/modules/healthcheck-parameters';

import { healthcheckParametersToForm } from '@/helpers/entities/healthcheck/form';

import HealthcheckForm from '@/components/other/healthcheck/form/healthcheck-form.vue';

export default {
  components: { HealthcheckForm },

  data() {
    return {
      form: null,
      errors: {},
      submitting: false,
    };
  },

  computed: {
    isDisabled() {
      return this.submitting;
    },

    engineNames() {
      return Object.values(HEALTHCHECK_ENGINES_NAMES);
    },
  },

  async mounted() {
    const healthcheckParameters = await useHealthcheckParametersStore().fetchItemWithoutStore();

    this.form = healthcheckParametersToForm(healthcheckParameters);
  },

  methods: {
    addError(field, message) {
      if (!this.errors[field]) {
        this.errors[field] = [];
      }

      this.errors[field].push(message);
    },

    validateLimit(path, value) {
      if (!value.enabled) {
        return;
      }

      if (!Number.isFinite(Number(value.limit)) || Number(value.limit) < 1) {
        this.addError(`${path}.limit`, this.$t('validation.min_value', [1]));
      }
    },

    validateEngine(engineName, value) {
      if (!value.enabled) {
        return;
      }

      const minimal = Number(value.minimal);
      const optimal = Number(value.optimal);

      if (!Number.isFinite(minimal) || minimal < 0) {
        this.addError(`${engineName}.minimal`, this.$t('validation.min_value', [0]));
      }

      if (!Number.isFinite(optimal) || optimal < 0) {
        this.addError(`${engineName}.optimal`, this.$t('validation.min_value', [0]));
      }

      if (Number.isFinite(minimal) && Number.isFinite(optimal) && minimal > optimal) {
        this.addError(`${engineName}.minimal`, this.$t('healthcheck.validation.max_value'));
        this.addError(`${engineName}.optimal`, this.$t('healthcheck.validation.min_value'));
      }
    },

    validate() {
      this.errors = {};

      this.validateLimit('queue', this.form.queue);
      this.validateLimit('messages', this.form.messages);
      this.engineNames.forEach(engineName =>
        this.validateEngine(engineName, this.form[engineName]),
      );

      return !Object.keys(this.errors).length;
    },

    async submit() {
      if (!this.validate()) {
        return;
      }

      try {
        this.submitting = true;

        await useHealthcheckParametersStore().update({ data: this.form });

        this.$popups.success({ text: this.$t('success.default') });
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.healthcheck-parameters {
  &__form {
    width: 83.3333%;
    max-width: 960px;
    margin-right: auto;
    margin-left: auto;
  }
}
</style>

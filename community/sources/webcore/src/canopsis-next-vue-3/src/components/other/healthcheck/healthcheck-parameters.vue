<template>
  <div class="healthcheck-parameters">
    <div v-if="!values" class="my-2 d-flex justify-center">
      <v-progress-circular color="primary" indeterminate />
    </div>
    <div v-else class="healthcheck-parameters__form py-2">
      <v-form @submit.prevent="submit">
        <healthcheck-form :form="values" :errors="errors" @input="values = $event" />
        <div class="mt-3 d-flex justify-end">
          <v-btn
            :loading="submitting"
            :disabled="submitting"
            color="primary"
            type="submit"
          >
            {{ t('common.submit') }}
          </v-btn>
        </div>
      </v-form>
    </div>
  </div>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { HEALTHCHECK_ENGINES_NAMES } from '@/constants';

import { healthcheckParametersToForm } from '@/helpers/entities/healthcheck/form';

import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useHealthcheckParameters } from '@/hooks/store/modules/healthcheck-parameters';

import HealthcheckForm from '@/components/other/healthcheck/form/healthcheck-form.vue';

export default {
  components: { HealthcheckForm },

  setup() {
    const { t } = useI18n();
    const popups = usePopups();
    const { fetchHealthcheckParametersWithoutStore, updateHealthcheckParameters } =
      useHealthcheckParameters();
    const { values, handleSubmit, submitting } = useForm({
      initialValues: null,
    });

    const errors = ref({});
    const engineNames = computed(() => Object.values(HEALTHCHECK_ENGINES_NAMES));

    const addError = (field, message) => {
      if (!errors.value[field]) {
        errors.value[field] = [];
      }

      errors.value[field].push(message);
    };

    const validateLimit = (path, value) => {
      if (!value.enabled) {
        return;
      }

      if (!Number.isFinite(Number(value.limit)) || Number(value.limit) < 1) {
        addError(`${path}.limit`, t('validation.min_value', [1]));
      }
    };

    const validateEngine = (engineName, value) => {
      if (!value.enabled) {
        return;
      }

      const minimal = Number(value.minimal);
      const optimal = Number(value.optimal);

      if (!Number.isFinite(minimal) || minimal < 0) {
        addError(`${engineName}.minimal`, t('validation.min_value', [0]));
      }

      if (!Number.isFinite(optimal) || optimal < 0) {
        addError(`${engineName}.optimal`, t('validation.min_value', [0]));
      }

      if (Number.isFinite(minimal) && Number.isFinite(optimal) && minimal > optimal) {
        addError(`${engineName}.minimal`, t('healthcheck.validation.max_value'));
        addError(`${engineName}.optimal`, t('healthcheck.validation.min_value'));
      }
    };

    const validateParameters = () => {
      errors.value = {};

      validateLimit('queue', values.value.queue);
      validateLimit('messages', values.value.messages);
      engineNames.value.forEach(engineName => validateEngine(engineName, values.value[engineName]),
      );

      return !Object.keys(errors.value).length;
    };

    const submit = handleSubmit(async () => {
      if (!validateParameters()) {
        return;
      }

      await updateHealthcheckParameters({ data: values.value });

      popups.success({ text: t('success.default') });
    });

    onMounted(async () => {
      const healthcheckParameters = await fetchHealthcheckParametersWithoutStore();

      values.value = healthcheckParametersToForm(healthcheckParameters);
    });

    return {
      values,
      errors,
      submitting,
      submit,
      t,
    };
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

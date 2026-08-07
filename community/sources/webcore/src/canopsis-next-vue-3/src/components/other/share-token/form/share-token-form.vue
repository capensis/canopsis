<template>
  <c-layout column>
    <c-name-field
      v-model="description"
      name="description"
      autofocus
      required
    />

    <c-information-block :title="t('shareToken.tokenExpiration')">
      <c-enabled-field v-model="durationEnabled" />

      <c-duration-field
        v-if="durationEnabled"
        v-model:duration="duration"
        name="duration"
        long
      />
    </c-information-block>
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();
    const {
      description,
      durationEnabled,
      duration,
    } = useFormModels({
      description: 'description',
      durationEnabled: 'duration_enabled',
      duration: 'duration',
    });

    return {
      t,
      description,
      durationEnabled,
      duration,
    };
  },
};
</script>

<template>
  <v-row no-gutters>
    <v-col cols="12">
      <c-name-field
        v-model="description"
        class="mb-4"
        name="description"
        autofocus
        required
      />
    </v-col>

    <v-col cols="12">
      <token-expiration-field
        v-model:enabled="durationEnabled"
        v-model:duration="duration"
        :label="t('shareToken.tokenExpiration')"
      />
    </v-col>
  </v-row>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import TokenExpirationField from '@/components/other/token/form/token-expiration-field.vue';

export default {
  components: { TokenExpirationField },

  setup() {
    const { t } = useI18n();
    const { description, durationEnabled, duration } = useFormModels({
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

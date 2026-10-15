<template>
  <v-row no-gutters>
    <v-col cols="12">
      <c-name-field
        v-model="name"
        :max-length="255"
        class="mb-4"
        autofocus
        required
      />
    </v-col>

    <v-col cols="12">
      <token-expiration-field
        v-model:enabled="durationEnabled"
        v-model:duration="duration"
        :label="t('apiToken.tokenExpiration')"
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
    const { name, durationEnabled, duration } = useFormModels({
      name: 'name',
      durationEnabled: 'duration_enabled',
      duration: 'duration',
    });

    return {
      t,
      name,
      durationEnabled,
      duration,
    };
  },
};
</script>

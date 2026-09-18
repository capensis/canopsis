<template>
  <v-row>
    <v-col cols="12">
      <c-text-field
        v-model="usernameField"
        :label="t('common.username')"
        color="primary"
        name="username"
        variant="outlined"
        density="comfortable"
        required
        autofocus
        clearable
      />
    </v-col>
    <v-col cols="12">
      <c-password-field
        v-model="passwordField"
        :label="t('common.password')"
        color="primary"
        name="password"
        variant="outlined"
        density="comfortable"
        required
        visibility
        clearable
      />
    </v-col>
  </v-row>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

import { useFormModels } from '@/hooks/form';

const { t } = useI18n();
const { usernameField, passwordField } = useFormModels({
  usernameField: 'username',
  passwordField: 'password',
});
</script>

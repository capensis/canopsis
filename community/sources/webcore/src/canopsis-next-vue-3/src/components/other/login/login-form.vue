<template>
  <c-layout column>
    <c-text-field
      v-model="usernameField"
      :label="t('common.username')"
      color="primary"
      name="username"
      variant="outlined"
      required
      autofocus
      clearable
    />

    <c-password-field
      v-model="passwordField"
      :label="t('common.password')"
      color="primary"
      name="password"
      variant="outlined"
      required
      visibility
      clearable
    />
  </c-layout>
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

<script setup>
import { computed } from 'vue';

import logo from '@/assets/canopsis.png';

import BasicLogin from './basic-login.vue';
import CasLogin from './cas-login.vue';
import SamlLogin from './saml-login.vue';
import OauthLogin from './oauth-login.vue';

const props = defineProps({
  basic: {
    type: Boolean,
    default: true,
  },
  cas: {
    type: Boolean,
    default: false,
  },
  saml: {
    type: Boolean,
    default: false,
  },
  oauth: {
    type: Boolean,
    default: false,
  },
});

const loginComponents = {
  basic: BasicLogin,
  cas: CasLogin,
  saml: SamlLogin,
  oauth: OauthLogin,
};

const components = computed(() => (
  Object.entries(loginComponents).reduce((acc, [key, component]) => {
    if (props[key]) {
      acc.push({ key, component });
    }

    return acc;
  }, [])
));
</script>

<template>
  <v-card class="login-card">
    <v-card-title class="bg-primary text-white">
      <div class="d-flex justify-space-between align-center w-100">
        <h3>{{ $t('common.login') }}</h3>

        <img
          :src="logo"
          class="login-card__logo"
          alt="Canopsis logo"
        >
      </div>
    </v-card-title>

    <v-card-text>
      <c-layout column>
        <component
          v-for="{ key, component } in components"
          :is="component"
          :key="key"
          class="login-card__item py-4"
        />
      </c-layout>
    </v-card-text>
  </v-card>
</template>

<style lang="scss" scoped>
.login-card {
  grid-area: form;
  width: 100%;
  display: flex;
  flex-flow: column;
  justify-content: space-between;

  &__logo {
    max-width: 35%;
    max-height: 4em;
    object-fit: scale-down;
  }

  &__item {
    &:not(:last-child) {
      border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
    }
  }
}
</style>

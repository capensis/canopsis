<template>
  <v-card class="login-card">
    <v-card-title class="primary white--text">
      <v-layout
        justify-space-between
        align-center
      >
        <h3>{{ $t('common.login') }}</h3>
        <img
          class="login-card__logo"
          src="@/assets/canopsis.png"
          alt=""
        >
      </v-layout>
    </v-card-title>
    <v-card-text>
      <v-layout column>
        <component
          v-for="{ key, component } in components"
          :is="component"
          :key="key"
          class="login-card__item py-4"
        />
      </v-layout>
    </v-card-text>
  </v-card>
</template>
<script>
import { computed } from 'vue';

import BasicLogin from './basic-login.vue';
import CasLogin from './cas-login.vue';
import SamlLogin from './saml-login.vue';
import OauthLogin from './oauth-login.vue';

export default {
  components: {
    BasicLogin,
    CasLogin,
    SamlLogin,
    OauthLogin,
  },
  props: {
    basic: {
      type: Boolean,
      default: true,
    },
    cas: {
      type: Boolean,
      default: true,
    },
    saml: {
      type: Boolean,
      default: true,
    },
    oauth: {
      type: Boolean,
      default: true,
    },
  },
  setup(props) {
    const components = computed(() => (
      ['basic', 'cas', 'saml', 'oauth'].reduce((acc, key) => {
        if (props[key]) {
          acc.push({ key, component: `${key}-login` });
        }

        return acc;
      }, [])
    ));

    return { components };
  },
};
</script>

<style lang="scss" scoped>
.login-card {
  grid-area: form;
  width: 100%;
  min-height: 500px;
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
      border-bottom: 1px solid var(--v-divider-border-color);
    }
  }
}
</style>

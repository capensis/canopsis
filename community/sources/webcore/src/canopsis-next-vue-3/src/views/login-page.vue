<script setup>
import { onMounted, onBeforeUnmount } from 'vue';

import { LOGIN_APP_INFO_POLLING_DELAY } from '@/constants';

import { useInfo } from '@/hooks/store/modules/info';

import LoginCard from '@/components/other/login/login-card.vue';
import LoginFooter from '@/components/other/login/login-footer.vue';

const {
  maintenance,
  description,
  isCASAuthEnabled,
  isSAMLAuthEnabled,
  isOauthAuthEnabled,
  fetchAppInfo,
} = useInfo();

let timeoutId = null;

const stopAppInfoPolling = () => {
  if (timeoutId) {
    clearTimeout(timeoutId);
  }
};

const startAppInfoPolling = () => {
  stopAppInfoPolling();

  timeoutId = setTimeout(async () => {
    await fetchAppInfo();

    startAppInfoPolling();
  }, LOGIN_APP_INFO_POLLING_DELAY);
};

onMounted(startAppInfoPolling);
onBeforeUnmount(stopAppInfoPolling);
</script>

<template>
  <div :class="{ 'login--maintenance': maintenance }" class="login">
    <div v-if="maintenance" class="login__image">
      <v-icon size="120" color="white">
        manufacturing
      </v-icon>
    </div>
    <div class="login__description">
      <c-compiled-template :template="description" />
    </div>
    <login-card
      :basic="true"
      :cas="isCASAuthEnabled"
      :saml="isSAMLAuthEnabled"
      :oauth="isOauthAuthEnabled"
    />
    <login-footer />
  </div>
</template>

<style lang="scss" scoped>
.login {
  min-width: 100%;
  min-height: 100vh;
  overflow-x: hidden;
  display: grid;
  align-items: center;

  grid-template-columns: 1fr 8fr 1fr;
  grid-template-rows: 5% auto auto 15% auto;
  background: var(--v-secondary-base);

  grid-template-areas:
    '. . .'
    '. image .'
    '. description .'
    '. form .'
    '. . .'
    'footer footer footer';

  &__image {
    grid-area: image;
    display: flex;
    justify-content: center;
    align-items: center;
    pointer-events: none;
  }

  &__description {
    grid-area: description;
    align-content: center;
    min-height: 500px;
    width: 100%;
    overflow-x: hidden;
    overflow-y: auto;
    color: white;
  }

  &__container {
    grid-area: form;
    width: 100%;
    min-height: 500px;
    display: flex;
    flex-flow: column;
    justify-content: space-between;
  }

  &--maintenance {
    background: #6a6a6a;
  }

  @media (min-width: 900px) {
    grid-template-columns: auto 40% 1% 40% auto;
    grid-template-rows: auto auto auto auto;

    grid-template-areas:
      '. . . . .'
      'image description . form .'
      '. . . . .'
      'footer footer footer footer footer';

    &__image {
      min-height: 500px;
      padding-top: 80px;
    }
  }

  @media (min-width: 1200px) {
    grid-template-columns: auto 30% 3% 30% auto;
  }
}
</style>

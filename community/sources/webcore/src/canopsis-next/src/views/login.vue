<template lang="pug">
  div.mainContainer.secondary
    div.description
      div(v-html="description")
    div.loginContainer
      base-login
      cas-login.mt-2(v-if="isCASAuthEnabled", key="cas")
      saml-login.mt-2(v-if="isSAMLAuthEnabled", key="saml")
    login-footer
</template>

<script>

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import BaseLogin from '@/components/other/login/base-login.vue';
import CasLogin from '@/components/other/login/cas-login.vue';
import SamlLogin from '@/components/other/login/saml-login.vue';
import LoginFooter from '@/components/other/login/login-footer.vue';

export default {
  components: {
    BaseLogin,
    CasLogin,
    SamlLogin,
    LoginFooter,
  },
  mixins: [authMixin, entitiesInfoMixin],
};
</script>

<style lang="scss" scoped>
  .mainContainer {
    min-width: 100%;
    min-height: 100vh;
    overflow-x: hidden;
    display: grid;
    align-items: center;

    grid-template-columns: 1fr 8fr 1fr;
    grid-template-rows: 5% auto auto 15% auto;

    grid-template-areas:
      ". . ."
      ". description ."
      ". form ."
      ". . ."
      "footer footer footer";

    @media (min-width: 900px) {
      grid-template-columns: auto 40% 1% 40% auto;
      grid-template-rows: auto auto auto auto;

      grid-template-areas:
        ". . . . ."
        ". description . form ."
        ". . . . ."
        "footer footer footer footer footer";
    }

    @media (min-width: 1200px) {
      grid-template-columns: auto 30% 3% 30% auto;
    }
  }

  .description {
    grid-area: description;
    align-content: center;
    min-height: 500px;
    width: 100%;
    overflow-x: hidden;
    overflow-y: auto;
    color: white;
  }

  .loginContainer {
    grid-area: form;
    width: 100%;
    min-height: 500px;
    display: flex;
    flex-flow: column;
    justify-content: space-between;
  }
</style>

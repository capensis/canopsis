<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        h3 {{ $t('common.login') }}
        img.secondaryLogo(src="@/assets/canopsis.png")
    v-card-text
      v-form.pa-2(data-test="loginForm", @submit.prevent.stop="submit")
        ldap-login-information(v-if="isLDAPAuthEnabled")
        login-form(v-field.model="form")
        v-flex
          v-layout.mb-1(justify-space-between, align-center)
            v-btn.ma-0(
              type="submit",
              color="primary",
              data-test="submitButton"
            ) {{ $t('common.connect') }}
            v-flex(v-if="hasServerError", xs9, data-test="errorLogin")
              v-alert.py-1.my-0.font-weight-bold(:value="hasServerError", type="error")
                span {{ $t('login.errors.incorrectEmailOrPassword') }}
          template(v-if="footer")
            v-divider.my-2
            v-layout
              div(v-html="footer")
</template>

<script>
import { ROUTES_NAMES, ROUTES } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import LdapLoginInformation from '@/components/other/login/ldap-login-information.vue';
import LoginForm from '@/components/other/login/form/login-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    LdapLoginInformation,
    LoginForm,
  },
  mixins: [authMixin, entitiesInfoMixin],
  data() {
    return {
      hasServerError: false,
      form: {
        username: '',
        password: '',
      },
    };
  },
  methods: {
    async submit() {
      try {
        this.hasServerError = false;

        const formIsValid = await this.$validator.validateAll();

        if (formIsValid) {
          const { redirect } = this.$route.query;
          const { defaultview: userDefaultView } = this.currentUser;

          await this.login(this.form);

          if (redirect && redirect !== ROUTES.home) {
            this.$router.push(redirect);
            return;
          }

          if (userDefaultView) {
            this.$router.push({
              name: ROUTES_NAMES.view,
              params: { id: userDefaultView._id },
            });
            return;
          }

          this.$router.push({ name: ROUTES_NAMES.home });
        }
      } catch (err) {
        this.hasServerError = true;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .secondaryLogo {
    max-width: 35%;
    max-height: 4em;
    object-fit: scale-down;
  }
</style>

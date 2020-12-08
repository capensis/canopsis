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
          v-divider
          v-layout(v-if="footer", data-test="loginFormFooter")
            div(v-html="footer")
</template>

<script>
import authMixin from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

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
          await this.login(this.form);
          await this.fetchAppInfos();

          if (this.$route.query.redirect && this.$route.query.redirect !== '/') {
            this.$router.push(this.$route.query.redirect);
          } else if (this.currentUser.defaultview) {
            this.$router.push({
              name: 'view',
              params: { id: this.currentUser.defaultview },
            });
          } else {
            this.$router.push({ name: 'home' });
          }
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

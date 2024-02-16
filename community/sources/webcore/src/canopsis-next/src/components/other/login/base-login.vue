<template>
  <v-card>
    <v-card-title class="primary white--text">
      <v-layout
        justify-space-between
        align-center
      >
        <h3>{{ $t('common.login') }}</h3>
        <img
          class="secondaryLogo"
          src="@/assets/canopsis.png"
          alt=""
        >
      </v-layout>
    </v-card-title>
    <v-card-text>
      <v-form
        class="pa-2"
        @submit.prevent.stop="submit"
      >
        <ldap-login-information v-if="isLDAPAuthEnabled" />
        <login-form v-field.model="form" />
        <v-flex>
          <v-layout
            class="mb-1"
            justify-space-between
            align-center
          >
            <v-btn
              type="submit"
              color="primary"
            >
              {{ $t('common.connect') }}
            </v-btn>
            <v-flex xs9>
              <c-alert
                class="py-1 my-0 font-weight-bold"
                :value="!!serverErrorMessage"
                type="error"
              >
                {{ serverErrorMessage }}
              </c-alert>
            </v-flex>
          </v-layout>
          <template v-if="footer">
            <v-divider class="my-2" />
            <v-layout>
              <c-compiled-template :template="footer" />
            </v-layout>
          </template>
        </v-flex>
      </v-form>
    </v-card-text>
  </v-card>
</template>

<script>
import { ROUTES_NAMES, ROUTES, RESPONSE_STATUSES } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import LdapLoginInformation from './partials/ldap-login-information.vue';
import LoginForm from './form/login-form.vue';

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
      serverErrorMessage: '',
      form: {
        username: '',
        password: '',
      },
    };
  },
  methods: {
    async submit() {
      try {
        this.serverErrorMessage = '';

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
        if (err?.status === RESPONSE_STATUSES.serviceUnavailable) {
          this.serverErrorMessage = this.$t('login.errors.underMaintenance');
        } else {
          this.serverErrorMessage = this.$t('login.errors.incorrectEmailOrPassword');
        }
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

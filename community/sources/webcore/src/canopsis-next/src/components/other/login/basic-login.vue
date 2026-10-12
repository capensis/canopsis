<template>
  <v-form
    class="gap-3"
    @submit.prevent.stop="submit"
  >
    <v-layout class="gap-2" column>
      <v-flex>
        <ldap-login-information v-if="isLDAPAuthEnabled" />
      </v-flex>
      <login-form v-field.model="form" />
      <v-flex>
        <v-layout
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
              :value="!!serverErrorMessage"
              class="py-1 my-0 font-weight-bold"
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
    </v-layout>
  </v-form>
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

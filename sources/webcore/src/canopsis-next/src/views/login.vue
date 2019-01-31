<template lang="pug">
  v-container.secondary(
  fill-height,
  fluid,
  d-flex,
  align-center
  )
    v-layout(justify-center, align-center, row)
      v-flex(xs11, md6, lg4)
        v-card
          v-tabs(v-model="activeTab", color="primary", dark, slider-color="secondary")
            v-tab Standard
            v-tab LDAP
            v-tab-item
              v-card
                v-card-text
                  v-layout(wrap)
                    v-flex(xs12)
                      v-layout(justify-center)
                        img.my-4(src="@/assets/canopsis-green.png")
                    v-flex(xs12)
                      login-form(v-model="standardForm", :hasServerError="hasServerError")
                  v-layout(justify-center)
                    v-btn Connect with WebSSO
                  v-layout(justify-end)
                    v-btn(@click="standardLogin", type="submit", color="primary") {{ $t('common.connect') }}
            v-tab-item
              v-card
                v-card-text
                  v-layout(wrap)
                    v-flex(xs12)
                      v-layout(justify-center)
                        img.my-4(src="@/assets/canopsis-green.png")
                    v-flex(xs12)
                      login-form(v-model="ldapForm", :hasServerError="hasServerError")
                  v-layout(justify-center)
                    v-btn Connect with WebSSO
                  v-layout(justify-end)
                    v-btn(@click="ldapLogin", type="submit", color="primary") {{ $t('common.connect') }}
</template>

<script>
import authMixin from '@/mixins/auth';

import LoginForm from '@/components/forms/login.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    LoginForm,
  },
  mixins: [authMixin],
  data() {
    return {
      standardForm: {
        username: '',
        password: '',
      },
      ldapForm: {
        username: '',
        password: '',
      },
      hasServerError: false,
      activeTab: 0,
    };
  },
  methods: {
    redirect() {
      if (
        this.$route.query.redirect &&
        this.$route.query.redirect !== '/'
      ) {
        this.$router.push(this.$route.query.redirect);
      } else if (this.currentUser.defaultview) {
        this.$router.push({
          name: 'view',
          params: { id: this.currentUser.defaultview },
        });
      } else {
        this.$router.push({ name: 'home' });
      }
    },

    async standardLogin() {
      try {
        this.hasServerError = false;

        const formIsValid = await this.$validator.validateAll();

        if (formIsValid) {
          await this.login(this.standardForm);

          this.redirect();
        }
      } catch (err) {
        this.hasServerError = true;
      }
    },

    async ldapLogin() {
      try {
        this.hasServerError = false;

        const formIsValid = await this.$validator.validateAll();

        if (formIsValid) {
          await this.login(this.form);

          this.redirect();
        }
      } catch (err) {
        this.hasServerError = true;
      }
    },
  },
};
</script>

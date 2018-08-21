<template lang="pug">
  v-container.blue.darken-4(
  fill-height,
  fluid,
  d-flex,
  align-center
  )
    v-layout(justify-center, align-center, row)
      v-flex(xs11, md6, lg5)
        v-card
          v-layout(row, wrap)
            v-flex(xs12)
              v-toolbar.green.darken-2.white--text
                v-toolbar-title {{ $t('common.login') }}
            v-flex(xs12, py-2)
              v-form(@submit.prevent="submit")
                v-flex(px-3)
                  v-alert(:value="hasServerError", type="error")
                    span {{ $t('login.errors.incorrectEmailOrPassword') }}
                v-flex(px-3)
                  v-text-field(
                  :label="$t('common.username')"
                  autofocus,
                  clearable,
                  color="blue darken-4",
                  v-model="form.username",
                  name="username",
                  v-validate="'required'",
                  data-vv-name="username",
                  :error-messages="errors.collect('username')",
                  )
                v-flex(px-3)
                  v-text-field(
                  :label="$t('common.password')"
                  clearable,
                  color="blue darken-4",
                  v-model="form.password",
                  type="password",
                  name="password",
                  v-validate="'required'",
                  data-vv-name="password",
                  :error-messages="errors.collect('password')",
                  )
                v-flex(xs2 px-2)
                  v-btn(
                  type="submit",
                  color="blue darken-4 white--text"
                  ) {{ $t('common.submit') }}
</template>

<script>
import authMixin from '@/mixins/auth';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [authMixin],
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

          if (this.$route.query.redirect) {
            this.$router.push(this.$route.query.redirect);
          } else {
            this.$router.push({
              name: 'home',
            });
          }
        }
      } catch (err) {
        this.hasServerError = true;
      }
    },
  },
};
</script>

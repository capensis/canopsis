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
          v-layout(row, wrap)
            v-flex(xs12)
              v-toolbar.primary.white--text
                v-toolbar-title {{ $t('common.login') }}
            v-flex(xs12)
              v-layout(justify-center)
                img.my-4(src="@/assets/canopsis-green.png")
            v-flex(xs12)
              v-form.py-2(@submit.prevent="submit")
                v-flex(px-3)
                  v-text-field(
                  :label="$t('common.username')",
                  :error-messages="errors.collect('username')",
                  v-model="form.username",
                  v-validate="'required'",
                  color="primary",
                  name="username",
                  autofocus,
                  clearable,
                  outline
                  )
                v-flex(px-3)
                  v-text-field(
                  :label="$t('common.password')",
                  :error-messages="errors.collect('password')",
                  v-model="form.password",
                  v-validate="'required'",
                  color="primary",
                  name="password",
                  type="password",
                  clearable,
                  outline
                  )
                v-flex.px-3.py-2
                  v-alert(:value="hasServerError", type="error")
                    span {{ $t('login.errors.incorrectEmailOrPassword') }}
                v-flex(xs2 px-2)
                  v-btn.primary(
                  @click="submit",
                  type="submit",
                  ) {{ $t('common.connect') }}
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

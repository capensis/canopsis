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
                  v-text-field(
                    autofocus,
                    clearable,
                    color="blue darken-4",
                    v-model="form.username",
                    name="username",
                    :label="$t('common.username')"
                  )
                v-flex(px-3)
                  v-text-field(
                    clearable,
                    color="blue darken-4",
                    v-model="form.password",
                    type="password",
                    name="password",
                    :label="$t('common.password')"
                  )
                v-flex(xs2 px-2)
                  v-btn(
                    type="submit",
                    color="blue darken-4 white--text"
                  ) {{ $t('common.submit') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('auth');

export default {
  name: 'login-page',
  data() {
    return {
      form: {
        username: '',
        password: '',
      },
    };
  },
  methods: {
    ...mapActions(['login']),

    async submit() {
      try {
        await this.login({
          username: this.form.username,
          password: this.form.password,
        });

        if (this.$route.query.redirect) {
          this.$router.push(this.$route.query.redirect);
        } else {
          this.$router.push({
            name: 'home', // TODO: fix it
          });
        }
      } catch (err) {
        console.error(err);

        // TODO: check response
      }
    },
  },
};
</script>

<style scoped>
</style>


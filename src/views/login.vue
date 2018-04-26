<template lang="pug">
  v-container(
    class="blue darken-4"
    fluid
    d-flex
    align-center
  )
    v-layout(justify-center row)
      v-flex(xs11 md6 lg5)
        v-card
          v-layout(row wrap)
            v-flex(xs12)
              v-toolbar(class="green darken-2 white--text")
                v-toolbar-title {{ $t('common.login') }}
            v-flex(xs12 py-2)
              v-form
                v-flex(px-3)
                  v-text-field(
                    autofocus
                    clearable
                    color="blue darken-4"
                    v-model="formData.username"
                    name="username"
                    :label="$t('common.username')"
                  )
                v-flex(px-3)
                  v-text-field(
                    clearable
                    color="blue darken-4"
                    v-model="formData.password"
                    type="password"
                    name="password"
                    :label="$t('common.password')"
                  )
                v-flex(xs2 px-2)
                  v-btn(
                    @click="handleSubmit"
                    depressed
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
      formData: {
        username: '',
        password: '',
      },
    };
  },
  methods: {
    ...mapActions(['login']),

    handleSubmit() {
      this.login({ username: this.formData.username, password: this.formData.password });
    },
  },
};
</script>

<style scoped>
</style>


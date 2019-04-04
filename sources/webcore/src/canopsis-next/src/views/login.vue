<template lang="pug">
  v-container.secondary(
  fill-height,
  fluid,
  d-flex,
  align-center
  )
    div
      v-layout(justify-center, align-center, row)
        v-flex(xs11, md6, lg4)
          v-card
            v-card-title.primary.white--text.elevation-3
              v-layout(justify-space-between, align-center)
                v-toolbar-title {{ $t('common.login') }} - {{ appTitle }}
                img.px-2(v-if="logo", src="@/assets/canopsis.png")
            v-layout(row, wrap)
              v-flex(xs12)
                v-layout(justify-center)
                  img.my-4.logo(:src="logo")
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
                    v-layout
                      v-btn.primary(type="submit") {{ $t('common.connect') }}
                v-divider
                v-runtime-template(:template="footer")
      div.version.pr-2.mb-2 {{ version }}
</template>

<script>
import VRuntimeTemplate from 'v-runtime-template';

import authMixin from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

import canopsisLogo from '@/assets/canopsis-green.png';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    VRuntimeTemplate,
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
  computed: {
    logo() {
      if (this.logo) {
        return this.logo;
      }

      return canopsisLogo;
    },
  },
  async mounted() {
    this.fetchLoginInfos();
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

<style lang="scss" scoped>
  .version {
    position: absolute;
    right: 0.5em;
    bottom: 0.5em;
    color: white;
    font-weight: bold;
  }

  .logo {
    width: auto;
    height: auto;
    max-width: 15em;
    max-height: 15em;
  }
</style>

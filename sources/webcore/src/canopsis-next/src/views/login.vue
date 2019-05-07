<template lang="pug">
  v-container.secondary(fluid)
    v-layout.container(fill-height, justify-center)
      v-flex(md6, v-if="$options.filters.mq($mq, { md: true, l: true })")
        v-layout(justify-center, align-center)
          img.mainLogo(:src="appLogo")
      v-flex(xs10 md4)
        v-card
          v-card-title.primary.white--text
            v-layout(justify-space-between, align-center)
              h3 {{ title }}
              img.secondaryLogo(v-if="logo", src="@/assets/canopsis.png")
          v-card-text
            v-form.mt-3(@submit.prevent="submit")
              v-flex
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
              v-flex
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
              v-flex.px-3(v-if="hasServerError")
                v-alert(:value="hasServerError", type="error")
                  span {{ $t('login.errors.incorrectEmailOrPassword') }}
              v-flex
                v-layout(justify-space-between, align-center)
                  v-btn.primary(type="submit") {{ $t('common.connect') }}
                  v-btn.my-4(
                  :href="casHref",
                  v-if="isCASAuthEnabled",
                  color="secondary",
                  small
                  ) {{ casConfig | get('title', null, $t('login.loginWithCAS')) }}
    div.version.pr-2.mb-2 {{ version }}
</template>

<script>
import authMixin from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

import canopsisLogo from '@/assets/canopsis.png';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [authMixin, entitiesInfoMixin],
  data() {
    return {
      hasServerError: false,
      form: {
        username: '',
        password: '',
      },
      activeTab: 0,
    };
  },
  computed: {
    casHref() {
      if (this.casConfig) {
        return `${this.casConfig.server}/login?service=${this.casConfig.service}/logged_in`;
      }

      return null;
    },
    appLogo() {
      if (this.logo) {
        return this.logo;
      }

      return canopsisLogo;
    },

    title() {
      return this.isLDAPAuthEnabled ? `${this.$t('login.standard')}/${this.$t('login.LDAP')}` : this.$t('common.login');
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
  .container {
    min-height: 100vh;
    width: 100%;
    margin: 0;
    padding: 0;
  }

  .version {
    position: absolute;
    right: 0.5em;
    bottom: 0.5em;
    color: white;
    font-weight: bold;
  }

  .mainLogo {
    max-width: 80%;
    max-height: 20em;
    object-fit: scale-down;
  }

  .secondaryLogo {
    max-width: 40%;
    max-height: 4em;
    object-fit: scale-down;
  }
</style>

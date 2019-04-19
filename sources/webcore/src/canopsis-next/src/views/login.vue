<template lang="pug">
  div.mainContainer.secondary
    div.contentContainer
      div User text
      div.loginContainer
        v-card
          v-card-title.primary.white--text
            v-layout(justify-space-between, align-center)
              h3 {{ title }}
              img.secondaryLogo(src="@/assets/canopsis.png")
          v-card-text
            v-form.mt-3.pa-3(@submit.prevent="submit")
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
                  v-btn(type="submit", color="primary") {{ $t('common.connect') }}
        v-card.mt-2
          v-card-text
            div.pa-3
              div.ml-2.mb-2.font-weight-bold {{ $t('login.loginWithCAS') }}
              v-btn(color="primary") {{ $t('common.connect') }}
    div.secondary.darken-1.footer
      a(href="https://doc.canopsis.net/") Documentation
      a(href="https://www.capensis.fr/canopsis/") Canopsis.com
      a(:href="changeLogHref") Notes de version
      div.version {{ version }}
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
    changeLogHref() {
      if (this.version) {
        return `https://doc.canopsis.net/notes-de-versions/${this.version}/`;
      }

      return 'https://doc.canopsis.net/';
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
  .mainContainer {
    min-height: 100vh;
    min-width: 100%;
    margin: 0;
    padding: 0;
    display: flex;
    flex-flow: column;
  }

  .contentContainer {
    min-height: 600px;
    width: 60%;
    margin: auto;
    // background-color: white;

    display: flex;
    justify-content: space-between;
  }

  .loginContainer {
    flex-grow: 0.5;
    display: flex;
    flex-flow: column;
  }

  .secondaryLogo {
    max-width: 30%;
    max-height: 3em;
    object-fit: scale-down;
  }

  .footer {
    position: relative;
    color: white;
    height: 7em;
    margin-top: auto;
    display: flex;
    justify-content: center;
    align-items: center;

    a {
      color: inherit;
      text-decoration: underline;
      padding: 0 2em;
    }

    .version {
      line-height: 7em;
      position: absolute;
      right: 0.5em;
      bottom: 0;
    }
  }
</style>

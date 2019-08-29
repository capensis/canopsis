<template lang="pug">
  div.mainContainer.secondary
    div.description
      div(v-html="description")
    div.loginContainer
      v-card
        v-card-title.primary.white--text
          v-layout(justify-space-between, align-center)
            h3 {{ $t('common.login') }}
            img.secondaryLogo(src="@/assets/canopsis.png")
        v-card-text
          v-form.pa-2(data-test="loginForm", @submit.prevent="submit")
            v-tooltip.ma-1(right, v-if="isLDAPAuthEnabled")
              v-icon(slot="activator") help
              v-layout(wrap)
                v-flex(xs12) {{ $t('login.connectionProtocols') }}
                ul
                  li {{ $t('login.standard') }}
                  li(v-if="isLDAPAuthEnabled") {{ $t('login.LDAP') }}
            v-flex.mt-1
              v-text-field(
              :label="$t('common.username')",
              :error-messages="errors.collect('username')",
              v-model="form.username",
              v-validate="'required'",
              color="primary",
              name="username",
              data-test="username",
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
              data-test="password",
              clearable,
              outline
              )
            v-flex
              v-layout.mb-1(justify-space-between, align-center)
                v-btn.ma-0(
                type="submit",
                color="primary",
                data-test="submitButton"
                ) {{ $t('common.connect') }}
                v-flex(v-if="hasServerError", xs9, data-test="errorLogin")
                  v-alert.py-1.my-0.font-weight-bold(:value="hasServerError", type="error")
                    span {{ $t('login.errors.incorrectEmailOrPassword') }}
              v-divider
              v-layout(v-if="footer", v-html="footer", data-test="loginFormFooter")
      v-card.mt-2(v-show="isCASAuthEnabled")
        v-card-text
          div.pa-3
            div.ml-2.mb-2.font-weight-bold {{ $t('login.loginWithCAS') }}
            v-btn.my-4(
            :href="casHref",
            color="primary"
            ) {{ casConfig | get('title', null, $t('login.loginWithCAS')) }}
    div.secondary.darken-1.footer
      a(:href="$constants.CANOPSIS_DOCUMENTATION", target="_blank") {{ $t('login.documentation') }}
      a(:href="$constants.CANOPSIS_WEBSITE", target="_blank") Canopsis.com
      a(:href="$constants.CANOPSIS_FORUM", target="_blank") {{ $t('login.forum') }}
      a.version(:href="changeLogHref", target="_blank") {{ version }}
</template>

<script>
import { CANOPSIS_DOCUMENTATION } from '@/constants';

import authMixin from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

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
        return `${CANOPSIS_DOCUMENTATION}/notes-de-version/${this.version}/`;
      }

      return CANOPSIS_DOCUMENTATION;
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
          await this.fetchAppInfos();

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
    min-width: 100%;
    min-height: 100vh;
    overflow-x: hidden;
    display: grid;
    align-items: center;

    grid-template-columns: 1fr 8fr 1fr;
    grid-template-rows: 5% auto auto 15% auto;

    grid-template-areas:
      ". . ."
      ". description ."
      ". form ."
      ". . ."
      "footer footer footer";

    @media (min-width: 900px) {
      grid-template-columns: auto 40% 1% 40% auto;
      grid-template-rows: auto auto auto auto;

      grid-template-areas:
        ". . . . ."
        ". description . form ."
        ". . . . ."
        "footer footer footer footer footer";
    }

    @media (min-width: 1200px) {
      grid-template-columns: auto 30% 3% 30% auto;
    }
  }

  .description {
    grid-area: description;
    align-content: center;
    min-height: 500px;
    width: 100%;
    overflow-x: hidden;
    overflow-y: auto;
    color: white;
  }

  .loginContainer {
    grid-area: form;
    width: 100%;
    min-height: 500px;
    display: flex;
    flex-flow: column;
    justify-content: space-between;
  }

  .mainLogo {
    max-width: 80%;
    max-height: 5em;
    object-fit: scale-down;
  }

  .secondaryLogo {
    max-width: 35%;
    max-height: 4em;
    object-fit: scale-down;
  }

  .footer {
    grid-area: footer;
    position: relative;
    color: white;
    min-height: 5em;
    margin-top: auto;
    display: flex;
    justify-content: center;
    align-items: center;

    a {
      color: inherit;
      text-decoration: none;
      padding: 0 2em;
    }

    .version {
      line-height: 5em;
      position: absolute;
      right: 0.5em;
      bottom: 0;
      text-decoration: underline;
      font-weight: bold;

      color: inherit;
      text-decoration: none;
    }
  }
</style>

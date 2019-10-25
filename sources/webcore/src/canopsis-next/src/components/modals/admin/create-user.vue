<template lang="pug">
  v-card(data-test="createUserModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      user-form(
        v-model="form",
        :isNew="isNew",
        :user="config.user",
        :onlyUserPrefs="config.onlyUserPrefs"
      )
    v-divider
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(data-test="submitButton", @click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS, GROUPS_NAVIGATION_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import UserForm from '@/components/other/user/user-form.vue';

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createUser,
  $_veeValidate: {
    validator: 'new',
  },
  components: { UserForm },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        _id: '',
        firstname: '',
        lastname: '',
        mail: '',
        password: '',
        role: null,
        ui_language: '',
        enable: true,
        defaultview: '',
        groupsNavigationType: GROUPS_NAVIGATION_TYPES.sideBar,
      },
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createUser.title');
    },

    isNew() {
      return !this.config.user;
    },
  },
  async mounted() {
    if (!this.isNew) {
      this.form = pick(this.config.user, [
        '_id',
        'firstname',
        'lastname',
        'mail',
        'password',
        'role',
        'ui_language',
        'enable',
        'defaultview',
        'groupsNavigationType',
      ]);

      if (!this.form.groupsNavigationType) {
        this.form.groupsNavigationType = GROUPS_NAVIGATION_TYPES.sideBar;
      }
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.hideModal();
      }
    },
  },
};
</script>

<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        user-form(
          v-model="form",
          :isNew="isNew",
          :user="config.user",
          :onlyUserPrefs="config.onlyUserPrefs"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          type="submit",
          data-test="submitButton"
        ) {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS, GROUPS_NAVIGATION_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import UserForm from '@/components/other/user/user-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createUser,
  $_veeValidate: {
    validator: 'new',
  },
  components: { UserForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
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

        this.$modals.hide();
      }
    },
  },
};
</script>

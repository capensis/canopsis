<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        playlist-form(v-model="form", :groups="availableGroups")
      template(slot="actions")
        v-btn(
          data-test="createPbehaviorCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit",
          data-test="createPbehaviorSubmitButton"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { ENTITIES_TYPES, MODALS } from '@/constants';
import { getDefaultPlaylist } from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';
import submittableMixin from '@/mixins/submittable';

import PlaylistForm from '@/components/other/playlist/playlist-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPlaylist,

  $_veeValidate: {
    validator: 'new',
  },

  components: { PlaylistForm, ModalWrapper },
  mixins: [
    authMixin,
    modalInnerMixin,
    entitiesViewGroupMixin,
    rightsEntitiesGroupMixin,
    submittableMixin(),
  ],
  data() {
    return {
      pending: false,
      form: this.modal.config.playlist ? cloneDeep(this.modal.config.playlist) : getDefaultPlaylist(),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createPlaylist.create.title');
    },
  },
  async mounted() {
    this.pending = true;

    if (!this.groupsPending) {
      await this.fetchGroupsList();
    }

    if (this.form.tabs_list.length) {
      const tabs = this.getList(ENTITIES_TYPES.viewTab, this.form.tabs_list, true);

      this.form.tabs_list = tabs.filter(tab =>
        tab[SCHEMA_EMBEDDED_KEY].parents.some(parent => this.checkReadAccess(parent.id)));
    }

    this.pending = false;
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

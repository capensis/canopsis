<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        playlist-form(v-model="form", :groups="availableGroups")
      template(slot="actions")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { ENTITIES_TYPES, MODALS } from '@/constants';
import { getDefaultPlaylist } from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';
import submittableMixin from '@/mixins/submittable';

import PlaylistForm from '@/components/other/playlists/playlist-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapGetters: mapEntitiesGetters } = createNamespacedHelpers('entities');

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
      form: this.modal.config.playlist
        ? {
          ...this.modal.config.playlist,
          tabs_list: [],
        }
        : getDefaultPlaylist(),
    };
  },
  computed: {
    ...mapEntitiesGetters(['getList']),

    title() {
      return this.config.title || this.$t('modals.createPlaylist.create.title');
    },

    playlist() {
      return this.modal.config.playlist;
    },
  },
  async mounted() {
    this.pending = true;

    if (!this.groupsPending) {
      await this.fetchGroupsList();
    }

    if (this.playlist && this.playlist.tabs_list.length) {
      const tabs = this.getList(ENTITIES_TYPES.viewTab, this.playlist.tabs_list, true);

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
          await this.config.action({
            ...this.form,
            tabs_list: this.form.tabs_list.map(({ _id }) => _id),
          });
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

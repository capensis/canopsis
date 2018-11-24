<template lang="pug">
  v-speed-dial(
  v-if="hasCreateAnyViewAccess || hasUpdateAnyViewAccess || hasDeleteAnyViewAccess",
  v-model="isVSpeedDialOpen",
  transition="slide-y-reverse-transition"
  direction="top"
  bottom,
  right,
  fixed,
  )
    v-tooltip(slot="activator", left)
      v-btn.primary(slot="activator", :input-value="isVSpeedDialOpen", fab, dark)
        v-icon settings
        v-icon close
      span {{ $t('layout.sideBar.buttons.settings') }}
    v-tooltip(v-if="hasUpdateAnyViewAccess || hasDeleteAnyViewAccess", left)
      v-btn(
      slot="activator",
      :input-value="isEditingMode",
      color="blue darken-4",
      small,
      dark,
      fab,
      @click.stop="$emit('toggleEditingMode')"
      )
        v-icon(dark) edit
        v-icon(dark) done
      span {{ $t('layout.sideBar.buttons.edit') }}
    v-tooltip(v-if="hasCreateAnyViewAccess", left)
      v-btn(
      slot="activator",
      color="green darken-4",
      small,
      dark,
      fab,
      @click.stop="showCreateViewModal"
      )
        v-icon(dark) add
      span {{ $t('layout.sideBar.buttons.create') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal/modal';
import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

export default {
  mixins: [modalMixin, rightsTechnicalViewMixin],
  props: {
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      isVSpeedDialOpen: false,
    };
  },
  methods: {
    showCreateViewModal() {
      this.showModal({
        name: MODALS.createView,
      });
    },
  },
};
</script>

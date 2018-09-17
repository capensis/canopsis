<template lang="pug">
v-speed-dial.d-inline-block(
  direction="left",
  :open-on-hover="true",
  transition="scale-transition"
  )
  v-btn(slot="activator", color="green darken-3", dark, fab, small)
    v-icon add
    v-icon close
  v-tooltip(v-for="button in buttons", :key="button.label", top)
    v-btn(slot="activator", :color="button.color", @click.prevent="button.action", fab, dark, small)
      v-icon {{ button.icon }}
    span {{ button.label }}
</template>

<script>
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';
import entityMixin from '@/mixins/entities/context-entity';

/**
 * Buttons to open the modal to add entities
 *
 * @module context
 */
export default {
  mixins: [
    modalMixin,
    entityMixin,
  ],
  data() {
    return {
      buttons: [
        {
          label: this.$t('entities.watcher'),
          icon: 'watch',
          color: 'indigo',
          action: this.showCreateWatcherModal,
        },
        {
          label: this.$t('entities.entities'),
          icon: 'perm_identity',
          color: 'red',
          action: this.showCreateEntityModal,
        },
      ],
    };
  },
  methods: {
    showCreateEntityModal() {
      this.showModal({
        name: MODALS.createEntity,
        config: {
          title: 'modals.createEntity.createTitle',
        },
      });
    },
    showCreateWatcherModal() {
      this.showModal({
        name: MODALS.createWatcher,
        config: {
          title: 'modals.createWatcher.createTitle',
        },
      });
    },
  },
};
</script>

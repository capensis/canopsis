<template lang="pug">
  v-speed-dial.d-inline-block(
    direction="left",
    transition="scale-transition"
  )
    v-btn.primary(slot="activator", dark, fab, small)
      v-icon add
      v-icon close
    v-tooltip(v-for="button in buttons", :key="button.label", top)
      v-btn(slot="activator", :color="button.color", @click.prevent="button.action", fab, dark, small)
        v-icon {{ button.icon }}
      span {{ button.label }}
</template>

<script>
import { MODALS } from '@/constants';

import entitiesWatcherMixin from '@/mixins/entities/watcher';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

/**
 * Buttons to open the modal to add entities
 *
 * @module context
 */
export default {
  mixins: [
    entitiesWatcherMixin,
    entitiesContextEntityMixin,
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
      this.$modals.show({
        name: MODALS.createEntity,
        config: {
          title: this.$t('modals.createEntity.createTitle'),
          action: entity => this.createContextEntityWithPopup(entity),
        },
      });
    },
    showCreateWatcherModal() {
      this.$modals.show({
        name: MODALS.createWatcher,
        config: {
          title: this.$t('modals.createWatcher.createTitle'),
          action: watcher => this.createWatcherWithPopup(watcher),
        },
      });
    },
  },
};
</script>

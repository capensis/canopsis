<template lang="pug">
v-speed-dial(
  direction="top",
  :open-on-hover="true",
  transition="scale-transition")
  v-btn(slot='activator', v-model="fab", color='green darken-3', dark, fab)
    v-icon add
    v-icon close
  v-tooltip(left)
    v-btn(slot="activator", fab, dark, small, color='indigo', @click.prevent="showCreateWatcherModal")
      v-icon watch
    span {{ $t(`entities.watcher`)}}
  v-tooltip(left)
    v-btn(slot="activator", fab, dark, small, color='red', @click.prevent="showCreateEntityModal")
      v-icon perm_identity
    span {{ $t(`entities.entities`)}}
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
      fab: false,
    };
  },
  methods: {
    showCreateEntityModal() {
      this.showModal({
        name: MODALS.createEntity,
        config: {
          title: this.$t('modals.createEntity.createTitle'),
        },
      });
    },
    showCreateWatcherModal() {
      this.showModal({
        name: MODALS.createWatcher,
        config: {
          title: 'modals.createWatcher.title',
        },
      });
    },
  },
};
</script>

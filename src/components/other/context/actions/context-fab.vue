<template lang="pug">
v-speed-dial(
  direction="top",
  :open-on-hover="true",
  transition="scale-transition")
  v-btn(slot='activator', v-model="fab", color='blue darken-2', dark, fab)
    v-icon add
    v-icon close
  v-tooltip(left)
    v-btn(slot="activator", fab, dark, small, color='indigo')
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

export default {
  mixins: [modalMixin],
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
          title: 'modals.createEntity.title',
        },
      });
    },
  },
};
</script>

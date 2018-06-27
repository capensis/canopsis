<template lang="pug">
  v-card
    v-list
      v-list-group(
      v-for="(value, key) in config.item.props"
      :key="key"
      no-action
    )
        v-list-tile(slot="activator")
          v-list-tile-content
            v-list-tile-title {{ key }}
        v-list-tile(@click="")
          v-list-tile-content
            v-list-tile-title {{ value }}
</template>

<script>
import ModalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';


export default {
  name: MODALS.contextInfos,
  mixins: [ModalInnerMixin],
  props: {
    template: {
      type: String,
    },
  },
};
</script>

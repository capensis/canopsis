<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.statsSelect.title')}}
    v-container
      v-btn(@click="openAddStatModal") Stat selector modal
      v-card(v-if="value.stat")
        v-card-title.green.darken-4.white--text {{ value.title }}
        v-card-text {{ value.stat.value }}
</template>

<script>
import modalMixin from '@/mixins/modal/modal';
import { STATS_TYPES, MODALS } from '@/constants';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  computed: {
    /**
     * Get stats different types from constant, and return an object with stat's value and stat's translated title
     */
    statsTypes() {
      return Object.values(STATS_TYPES)
        .map(item => ({ value: item.value, text: this.$t(`stats.types.${item.value}`), options: item.options }));
    },
  },
  methods: {
    openAddStatModal() {
      this.showModal({
        name: MODALS.addStat,
        config: {
          action: stat => this.$emit('input', stat),
        },
      });
    },
  },
};
</script>

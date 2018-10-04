<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.filterEditor') }}
    div.text-xs-center
      v-btn(@click="openFilterModal") {{ openFilterButtonText }}
      v-btn(@click="deleteFilter", icon)
        v-icon delete
</template>

<script>
import isEmpty from 'lodash/isEmpty';

import { MODALS, ENTITIES_TYPES } from '@/constants';
import modalMixin from '@/mixins/modal/modal';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Object,
    },
    entitiesType: {
      type: String,
      default: ENTITIES_TYPES.alarm,
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity].includes(value),
    },
  },
  computed: {
    openFilterButtonText() {
      if (isEmpty(this.value)) {
        return this.$t('modals.filter.create.title');
      }

      return this.$t('modals.filter.edit.title');
    },
  },
  methods: {
    openFilterModal() {
      this.showModal({
        name: MODALS.createFilter,
        config: {
          title: 'modals.filter.create.title',
          filter: this.value || {},
          entitiesType: this.entitiesType,
          action: newFilter => this.$emit('input', newFilter),
        },
      });
    },
    deleteFilter() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('input', {}),
        },
      });
    },
  },
};
</script>

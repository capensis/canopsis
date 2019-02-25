<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsColor.title') }}
    v-container(fluid)
      v-layout(v-for="(stat, key) in stats", align-center, :key="key")
        v-flex
          div {{ key }}:
        v-flex
          v-btn(@click="showColorPickerModal(key)") {{ $t('settings.statsColor.pickColor') }}
        v-flex
          div.py-2.text-xs-center(:style="{ backgroundColor: value[key] }") {{ value[key] }}
</template>

<script>
import { set } from 'lodash';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Object,
      default() {
        return {};
      },
    },
    stats: {
      type: Object,
    },
  },
  methods: {
    showColorPickerModal(key) {
      const newVal = { ...this.value };
      this.showModal({
        name: this.$constants.MODALS.colorPicker,
        config: {
          title: this.$t('modals.colorPicker.title'),
          action: color => this.$emit('input', set(newVal, key, color)),
        },
      });
    },
  },
};
</script>


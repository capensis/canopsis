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
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import formMixin from '@/mixins/form';

export default {
  mixins: [modalMixin, formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    stats: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    showColorPickerModal(key) {
      this.showModal({
        name: MODALS.colorPicker,
        config: {
          title: this.$t('modals.colorPicker.title'),
          color: this.value[key],
          action: color => this.updateField(key, color),
        },
      });
    },
  },
};
</script>


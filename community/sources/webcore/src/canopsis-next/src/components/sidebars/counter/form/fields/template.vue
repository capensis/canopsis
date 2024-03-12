<template>
  <widget-settings-flat-item :title="title">
    <template #actions>
      <v-btn
        class="primary"
        small
        @click="showTextEditorModal"
      >
        {{ $t('common.show') }}/{{ $t('common.edit') }}
      </v-btn>
    </template>
  </widget-settings-flat-item>
</template>

<script>
import { MODALS } from '@/constants';

import WidgetSettingsFlatItem from '@/components/sidebars/partials/widget-settings-flat-item.vue';

export default {
  components: { WidgetSettingsFlatItem },
  props: {
    value: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    variables: {
      type: Array,
      required: false,
    },
  },
  methods: {
    showTextEditorModal() {
      this.$modals.show({
        name: MODALS.textEditor,
        config: {
          title: this.title,
          text: this.value,
          variables: this.variables,
          action: value => this.$emit('input', value),
        },
      });
    },
  },
};
</script>

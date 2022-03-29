<template lang="pug">
  c-action-btn(
    :tooltip="$t('filterSelector.buttons.list')",
    icon="filter_list",
    small,
    @click="showFiltersListModal"
  )
</template>

<script>
import { MODALS } from '@/constants';

export default {
  props: {
    widgetId: {
      type: String,
      required: true,
    },
    private: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    showFiltersListModal() {
      this.$modals.show({
        name: MODALS.filtersList,
        config: {
          widgetId: this.widgetId,
          private: this.private,
          addable: this.addable,
          editable: this.editable,
          withAlarm: true,
          // TODO: uncomment
          // withEntity: true,
          // withPbehavior: true,
        },
      });
    },
  },
};
</script>

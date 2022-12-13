<template lang="pug">
  div
    v-subheader {{ $t('common.general') }}
    v-data-table(
      :headers="headers",
      :items="items",
      hide-actions
    )
      template(#items="{ item }")
        td.text-xs-left {{ item.v.state.a }}
        td.text-xs-left {{ item.v.connector }}
        td.text-xs-left {{ item.v.component }}
        td.text-xs-left {{ item.v.resource }}
</template>

<script>
/**
 * General table component for alarm entity
 *
 * @module alarm
 *
 * @prop {Object} item - alarm item
 */
export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      headers: [
        {
          text: this.$t('common.author'),
          sortable: false,
        },
        {
          text: this.$t('common.connector'),
          sortable: false,
        },
        {
          text: this.$t('common.component'),
          sortable: false,
        },
        {
          text: this.$t('common.resource'),
          sortable: false,
        },
      ],
    };
  },
};
</script>

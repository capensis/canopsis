<template lang="pug">
  div
    v-subheader {{ $t('common.general') }}
    v-data-table(
      :headers="headers",
      :items="items",
      hide-actions
    )
      template(slot="items", slot-scope="props")
        td.text-xs-left {{ props.item.v.state.a }}
        td.text-xs-left {{ props.item.v.connector }}
        td.text-xs-left {{ props.item.v.component }}
        td.text-xs-left {{ props.item.v.resource }}
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
      required: true,
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
          text: this.$t('alarm.connector'),
          sortable: false,
        },
        {
          text: this.$t('alarm.component'),
          sortable: false,
        },
        {
          text: this.$t('alarm.resource'),
          sortable: false,
        },
      ],
    };
  },
};
</script>

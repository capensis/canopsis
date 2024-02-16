<template>
  <div>
    <v-subheader>{{ $t('common.general') }}</v-subheader>
    <v-data-table
      :headers="headers"
      :items="items"
      hide-default-footer
    >
      <template #items="{ item }">
        <td class="text-left">
          {{ item.v.state.a }}
        </td>
        <td class="text-left">
          {{ item.v.connector }}
        </td>
        <td class="text-left">
          {{ item.v.component }}
        </td>
        <td class="text-left">
          {{ item.v.resource }}
        </td>
      </template>
    </v-data-table>
  </div>
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

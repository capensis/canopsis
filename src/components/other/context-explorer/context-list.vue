<template lang="pug">
div
    v-subheader {{ $t('tables.contextList.title') }}
    v-toolbar.toolbar(dense, flat)
      v-text-field(
      label="Search",
      v-model="searchingText",
      hide-details,
      single-line,
      @keyup.enter="submit",
      @keyup.delete="clear",
      )
      v-btn(icon, @click="submit")
        v-icon search
      v-btn(icon, @click="clear")
        v-icon clear
    v-data-table(
      :no-data-text="this.$t('tables.contextList.noDataText')",
      :headers="headers",
      :items="items",
      hide-actions,
      )
      template(slot="items", slot-scope="props")
        td.text-xs-left {{ props.item.name }}
        td.text-xs-left {{ props.item._id}}
</template>

<script>
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
          text: this.$t('tables.contextList.name'),
          sortable: false,
        },
        {
          text: this.$t('tables.contextList.id'),
          sortable: false,
        },
      ],
    };
  },
};
</script>


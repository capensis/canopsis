<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") Row
    v-container
      v-combobox(
      :value="row",
      @change="updateRow"
      @blur="blur"
      :items="items",
      label="Row",
      :search-input.sync="search",
      data-vv-name="row",
      v-validate="'required'",
      :error-messages="errors.collect('row')",
      item-text="title",
      item-value="title"
      )
        template(slot="no-data")
          v-list-tile
            v-list-tile-content
              v-list-tile-title(v-html="$t('modals.createView.noData')")
      template(v-if="row")
        v-slider(
        :value="columnSM",
        :max="row.availableColumns.sm"
        @input="$emit('update:columnSM', $event)"
        v-validate="'min_value:3'",
        data-vv-name="columnSM",
        :error-messages="errors.first('columnSM')",
        ticks="always"
        always-dirty,
        thumb-label
        )
        v-slider(
        :value="columnMD",
        :max="row.availableColumns.md"
        @input="$emit('update:columnMD', $event)"
        v-validate="'min_value:3'",
        data-vv-name="columnMD",
        :error-messages="errors.first('columnMD')",
        ticks="always"
        always-dirty,
        thumb-label,
        )
        v-slider(
        :value="columnLG",
        :max="row.availableColumns.lg"
        @input="$emit('update:columnLG', $event)"
        v-validate="'min_value:3'",
        data-vv-name="columnLG",
        :error-messages="errors.first('columnLG')",
        ticks="always"
        always-dirty,
        thumb-label,
        )
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    selectedRowId: {
      type: String,
      default: null,
    },
    availableRows: {
      type: Array,
      default: () => [],
    },
    columnSM: {
      type: Number,
      default: 3,
    },
    columnMD: {
      type: Number,
      default: 3,
    },
    columnLG: {
      type: Number,
      default: 3,
    },
  },
  data() {
    let selectedRow = null;

    if (this.selectedRowId) {
      selectedRow = this.availableRows.find(row => row._id === this.selectedRowId);
    }

    return {
      row: selectedRow,
      search: null,
      items: [...this.availableRows],
    };
  },
  methods: {
    updateRow(value) {
      if (value !== this.row) {
        if (typeof value === 'string') {
          let newRow = this.availableRows.find(v => v.title === value);

          if (!newRow) {
            newRow = { title: value, _id: 'asdasd', availableColumns: { sm: 12, md: 12, lg: 12 } };
          }

          this.row = newRow;
          this.items.push(newRow);
        } else {
          this.row = value;
        }

        this.$emit('update:selectedRow', { ...this.row });
      }
    },
    blur() {
      this.search = this.row ? this.row.title : '';
    },
  },
};
</script>

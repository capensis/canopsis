<template lang="pug">
  div
    v-layout(justify-end)
      v-btn.primary(fab, small, flat, @click="showAddInfoModal")
        v-icon add
    v-data-table(
      :headers="headers",
      :items="form",
      :no-data-text="$t('tables.noData')",
      item-key="key"
    )
      template(slot="items", slot-scope="{ item }")
        tr
          td {{ item.name }}
          td {{ item.value }}
          td
            v-layout
              v-btn(icon, small, @click="showEditInfoModal(item)")
                v-icon edit
              v-btn(icon, small, @click="removeInfo(item)")
                v-icon(color="error") delete
</template>

<script>
import { MODALS } from '@/constants';

import formArrayMixin from '@/mixins/form/array';

export default {
  mixins: [formArrayMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Array,
      required: true,
    },
  },
  computed: {
    headers() {
      return [
        { text: this.$t('modals.createDynamicInfoInformation.fields.name'), value: 'name' },
        { text: this.$t('modals.createDynamicInfoInformation.fields.value'), value: 'value' },
        { text: this.$t('common.actionsLabel'), value: 'actions' },
      ];
    },
  },
  methods: {
    findInfoIndex(info = {}) {
      return this.form.findIndex(({ name }) => name === info.name);
    },

    showAddInfoModal() {
      this.$modals.show({
        name: MODALS.createDynamicInfoInformation,
        config: {
          existingNames: this.form.map(info => info.name),
          action: newInfo => this.addItemIntoArray(newInfo),
        },
      });
    },

    showEditInfoModal(info) {
      this.$modals.show({
        name: MODALS.createDynamicInfoInformation,
        config: {
          info,

          existingNames: this.form.map(({ name }) => name),
          action: (newInfo) => {
            const index = this.findInfoIndex(info);

            this.updateItemInArray(index, newInfo);
          },
        },
      });
    },

    removeInfo(info) {
      const index = this.findInfoIndex(info);

      this.removeItemFromArray(index);
    },
  },
};
</script>

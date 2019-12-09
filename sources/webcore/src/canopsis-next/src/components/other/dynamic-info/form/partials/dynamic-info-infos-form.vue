<template lang="pug">
  div
    v-layout(justify-end)
      v-btn.primary(fab, small, flat, @click="showAddInfoModal")
        v-icon add
    v-data-table(
      :headers="headers",
      :items="form",
      :no-data-text="$t('tables.noData')"
    )
      template(slot="items", slot-scope="{ item, index }")
        tr
          td {{ item.name }}
          td {{ item.value }}
          td
            v-layout
              v-btn(icon, small, @click="showEditInfoModal(index)")
                v-icon edit
              v-btn(icon, small, @click="removeItemFromArray(index)")
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
    showAddInfoModal() {
      this.$modals.show({
        name: MODALS.createDynamicInfoInformation,
        config: {
          action: info => this.addItemIntoArray(info),
        },
      });
    },
    showEditInfoModal(index) {
      this.$modals.show({
        name: MODALS.createDynamicInfoInformation,
        config: {
          info: this.form[index],
          action: info => this.updateItemInArray(index, info),
        },
      });
    },
  },
};
</script>

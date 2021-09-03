<template lang="pug">
  div
    v-layout(justify-end)
      c-action-btn(
        :tooltip="$t('entity.manageInfos.createTitle')",
        icon="add",
        @click="showAddInfoModal"
      )
    v-data-table(
      :items="infos",
      :headers="tableHeaders",
      :no-data-text="$t('entity.manageInfos.emptyInfos')",
      item-key="name"
    )
      template(slot="items", slot-scope="props")
        td {{ props.item.name }}
        td {{ props.item.description }}
        td {{ props.item.value }}
        td
          v-layout(row)
            c-action-btn(
              type="edit",
              @click="showEditInfoModal(props.index, props.item)"
            )
            c-action-btn(
              type="delete",
              @click="removeItemFromArray(props.index)"
            )
</template>

<script>
import { MODALS } from '@/constants';

import { formArrayMixin } from '@/mixins/form';

export default {
  mixins: [
    formArrayMixin,
  ],
  model: {
    prop: 'infos',
    event: 'input',
  },
  props: {
    infos: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      tableHeaders: [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.description'), value: 'description' },
        { text: this.$t('common.value'), value: 'value' },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ],
    };
  },
  methods: {
    showAddInfoModal() {
      this.$modals.show({
        name: MODALS.createEntityInfo,
        config: {
          infos: this.infos,
          action: info => this.addItemIntoArray(info),
        },
      });
    },

    showEditInfoModal(index, info) {
      this.$modals.show({
        name: MODALS.createEntityInfo,
        config: {
          infos: this.infos,
          entityInfo: info,
          title: this.$t('modals.createEntityInfo.edit.title'),
          action: editedInfo => this.updateItemInArray(index, editedInfo),
        },
      });
    },
  },
};
</script>

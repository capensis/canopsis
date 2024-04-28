<template>
  <div>
    <v-layout justify-end>
      <c-action-btn
        :tooltip="$t('entity.addInformation')"
        icon="add"
        @click="showAddInfoModal"
      />
    </v-layout>
    <v-data-table
      :items="infos"
      :headers="tableHeaders"
      :no-data-text="$t('entity.emptyInfos')"
      item-key="name"
    >
      <template #item="{ item, index }">
        <tr>
          <td>{{ item.name }}</td>
          <td>{{ item.description }}</td>
          <td>{{ item.value }}</td>
          <td>
            <v-layout>
              <c-action-btn
                type="edit"
                @click="showEditInfoModal(index, item)"
              />
              <c-action-btn
                type="delete"
                @click="removeItemFromArray(index)"
              />
            </v-layout>
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
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

<template lang="pug">
  v-card
    v-card-text
      v-layout(justify-end)
        v-tooltip(left)
          v-btn(@click="showAddInfoModal", slot="activator", color="secondary", icon)
            v-icon add
          span {{ $t('modals.createEntity.manageInfos.addInfo') }}
      v-data-table(
        :items="items",
        item-key="name",
        :headers="tableHeaders",
        :no-data-text="$t('modals.createEntity.manageInfos.noInfos')"
      )
        template(slot="items", slot-scope="props")
          td {{ props.item.name }}
          td {{ props.item.description }}
          td {{ props.item.value }}
          td
            v-layout(row)
              c-action-btn(
                type="edit",
                @click="showEditInfoModal(props.item)"
              )
              c-action-btn(
                type="delete",
                @click="removeField(props.item.name)"
              )
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import formMixin from '@/mixins/form';

/**
 * Form to manipulation with infos
 *
 * @prop {Object} infos - infos from parent
 */
export default {
  mixins: [
    formMixin,
  ],
  model: {
    prop: 'infos',
    event: 'input',
  },
  props: {
    infos: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      tableHeaders: [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.description'),
          value: 'description',
        },
        {
          text: this.$t('common.value'),
          value: 'value',
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ],
    };
  },
  computed: {
    items() {
      return Object.keys(this.infos).map(key => ({
        name: key,
        description: this.infos[key].description,
        value: this.infos[key].value,
      }));
    },
  },
  methods: {
    showAddInfoModal() {
      this.$modals.show({
        name: MODALS.addEntityInfo,
        config: {
          infos: this.infos,
          title: this.$t('modals.addEntityInfo.addTitle'),
          action: info => this.updateField(info.name, omit(info, ['name'])),
        },
      });
    },

    showEditInfoModal(info) {
      this.$modals.show({
        name: MODALS.addEntityInfo,
        config: {
          infos: this.infos,
          editingInfo: info,
          title: this.$t('modals.addEntityInfo.editTitle'),
          action: editedInfo => this.updateAndMoveField(info.name, editedInfo.name, omit(editedInfo, ['name'])),
        },
      });
    },
  },
};
</script>

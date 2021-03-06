<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.dynamicInfoTemplatesList.title') }}
    template(slot="text")
      div
        v-layout(justify-end)
          v-btn.primary(fab, small, flat, @click="showAddTemplateModal")
            v-icon add
        v-data-table(
          :items="templates",
          :headers="headers",
          :loading="pending",
          item-key="_id",
          expand
        )
          template(slot="items", slot-scope="props")
            tr(@click="props.expanded = !props.expanded")
              td {{ props.item.title }}
              td
                v-layout(row)
                  c-action-btn(
                    :tooltip="$t('modals.createDynamicInfo.create.title')",
                    icon="assignment",
                    @click="selectTemplate(props.item)"
                  )
                  c-action-btn(
                    type="edit",
                    @click="showEditTemplateModal(props.item)"
                  )
                  c-action-btn(
                    type="delete",
                    @click="showDeleteTemplateModal(props.item._id)"
                  )
          template(slot="expand", slot-scope="props")
            v-container.secondary.lighten-2
              v-card
                v-card-text
                  v-data-iterator(:items="props.item.names")
                    v-flex(slot="item", slot-scope="nameProps")
                      v-card
                        v-card-title {{ nameProps.item }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { templateToDynamicInfoInfos } from '@/helpers/forms/dynamic-info-template';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('dynamicInfoTemplate');

export default {
  name: MODALS.dynamicInfoTemplatesList,
  components: { ModalWrapper },
  computed: {
    ...mapGetters(['pending', 'templates']),

    headers() {
      return [
        {
          text: this.$t('common.title'),
          sortable: false,
          value: 'title',
        },
        {
          text: this.$t('common.actionsLabel'),
          sortable: false,
        },
      ];
    },
  },
  mounted() {
    this.fetchTemplatesList();
  },
  methods: {
    ...mapActions({
      fetchTemplatesList: 'fetchList',
      createTemplate: 'create',
      updateTemplate: 'update',
      deleteTemplate: 'delete',
    }),

    showAddTemplateModal() {
      this.$modals.show({
        name: MODALS.createDynamicInfoTemplate,
        config: {
          action: newTemplate => this.createTemplate({ data: newTemplate }),
        },
      });
    },

    showEditTemplateModal(template) {
      this.$modals.show({
        name: MODALS.createDynamicInfoTemplate,
        config: {
          template,

          title: this.$t('modals.createDynamicInfoTemplate.edit.title'),
          action: newTemplate => this.updateTemplate({ id: template._id, data: newTemplate }),
        },
      });
    },

    showDeleteTemplateModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.deleteTemplate({ id }),
        },
      });
    },

    selectTemplate(template) {
      this.$modals.show({
        name: MODALS.createDynamicInfo,
        config: {
          dynamicInfo: {
            infos: templateToDynamicInfoInfos(template),
          },
          action: async (data) => {
            if (this.config.action) {
              await this.config.action(data);
            }

            this.$modals.hide();
          },
        },
      });
    },
  },
};
</script>

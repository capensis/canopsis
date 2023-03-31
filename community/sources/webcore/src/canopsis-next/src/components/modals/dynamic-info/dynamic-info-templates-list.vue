<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.dynamicInfoTemplatesList.title') }}
    template(#text="")
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
          template(#items="props")
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
          template(#expand="{ item }")
            v-container.secondary.lighten-2
              v-card
                v-card-text
                  v-data-iterator(:items="item.names")
                    template(#item="nameProps")
                      v-flex
                        v-card
                          v-card-title {{ nameProps.item }}
</template>

<script>
import { MODALS } from '@/constants';

import { templateToDynamicInfoInfos } from '@/helpers/forms/dynamic-info-template';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesDynamicInfoTemplatesMixin } from '@/mixins/entities/associative-table/dynamic-info-templates';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.dynamicInfoTemplatesList,
  components: { ModalWrapper },
  mixins: [modalInnerMixin, entitiesDynamicInfoTemplatesMixin],
  data() {
    return {
      pending: false,
      templates: [],
    };
  },
  computed: {
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
    this.fetchList();
  },
  methods: {
    showAddTemplateModal() {
      this.$modals.show({
        name: MODALS.createDynamicInfoTemplate,
        config: {
          action: async (newTemplate) => {
            this.templates = await this.createDynamicInfoTemplate({ data: newTemplate });
          },
        },
      });
    },

    showEditTemplateModal(template) {
      this.$modals.show({
        name: MODALS.createDynamicInfoTemplate,
        config: {
          template,

          title: this.$t('modals.createDynamicInfoTemplate.edit.title'),
          action: async (newTemplate) => {
            this.templates = await this.updateDynamicInfoTemplate({ id: template._id, data: newTemplate });
          },
        },
      });
    },

    showDeleteTemplateModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            this.templates = await this.removeDynamicInfoTemplate({ id });
          },
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

    async fetchList() {
      this.pending = true;
      this.templates = await this.fetchDynamicInfoTemplatesList();
      this.pending = false;
    },
  },
};
</script>

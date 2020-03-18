<template lang="pug">
  modal-wrapper
    template(slot="title")
      span Dynamic info templates
    template(slot="text")
      div
        v-layout(justify-end)
          v-btn.primary(fab, small, flat, @click="showAddTemplateModal")
            v-icon add
        v-data-table(
          :items="templates",
          :headers="headers",
          :loading="pending",
          expand
        )
          template(slot="items", slot-scope="props")
            tr(@click="props.expanded = !props.expanded")
              td {{ props.item.title }}
              td
                v-layout
                  v-btn(
                    icon,
                    small,
                    @click.stop="showEditTemplateModal(props.item)"
                  )
                    v-icon done
                  v-btn(
                    icon,
                    small,
                    @click.stop="showEditTemplateModal(props.item)"
                  )
                    v-icon edit
                  v-btn(
                    icon,
                    small,
                    @click.stop="showDeleteTemplateModal(props.item._id)"
                  )
                    v-icon.error--text delete
          template(slot="expand", slot-scope="props")
            v-container.secondary.lighten-2
              v-card
                v-card-text
                  v-data-iterator(:items="props.item.values")
                    v-flex(slot="item", slot-scope="valueProps")
                      v-card
                        v-card-title {{ valueProps.item.name }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('dynamicInfoTemplate');

/**
 * Modal to create dynamic info's information
 */
export default {
  name: MODALS.dynamicInfoTemplatesList,
  components: { ModalWrapper },
  mixins: [modalInnerMixin],
  computed: {
    ...mapGetters(['pending', 'templates']),

    headers() {
      return [
        {
          text: 'Title',
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
          action: async (newTemplate) => {
            await this.createTemplate({ data: newTemplate });

            this.fetchTemplatesList();
          },
        },
      });
    },
    showEditTemplateModal(template) {
      this.$modals.show({
        name: MODALS.createDynamicInfoTemplate,
        config: {
          template,

          action: async (newTemplate) => {
            await this.createTemplate({ data: newTemplate });

            this.fetchTemplatesList();
          },
        },
      });
    },
    showDeleteTemplateModal() {

    },
    selectTemplate() {

    },
  },
};
</script>

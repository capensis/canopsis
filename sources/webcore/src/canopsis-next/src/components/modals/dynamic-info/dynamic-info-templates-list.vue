<template lang="pug">
  modal-wrapper
    template(slot="title")
      span Dynamic info templates
    template(slot="text")
      div
        v-layout(justify-end)
          v-btn.primary(fab, small, flat, @click="showAddTemplateModal")
            v-icon add
        v-data-table(:items="templates", :headers="headers", expand)
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
import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

import modalInnerMixin from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create dynamic info's information
 */
export default {
  name: MODALS.dynamicInfoTemplatesList,
  components: { ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    return {
      templates: [
        {
          id: uid(),
          title: 'Test template',
          values: [
            { name: 'attr_dynamique1' },
            { name: 'attr_dynamique2' },
            { name: 'attr_statique' },
          ],
        },
      ],
    };
  },
  computed: {
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
  methods: {
    showAddTemplateModal() {
      this.$modals.show({
        name: MODALS.createDynamicInfoTemplate,
      });
    },
    showEditTemplateModal() {

    },
    showDeleteTemplateModal() {

    },
    selectTemplate() {

    },
  },
};
</script>

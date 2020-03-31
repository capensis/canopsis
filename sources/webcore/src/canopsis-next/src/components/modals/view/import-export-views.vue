<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('common.import') }}/{{ $t('common.export') }}
    template(slot="text")
      v-layout(row)
        v-flex(xs4)
          draggable-groups(v-model="importedGroups")
        v-flex(xs4)
          draggable-group-views(v-model="importedViews")
        v-flex(xs4)
          draggable-groups(v-model="currentGroups")
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewsGroupsMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';

import DraggableGroupViews from '@/components/layout/navigation/partial/groups-side-bar/draggable-group-views.vue';
import DraggableGroups from '@/components/layout/navigation/partial/groups-side-bar/draggable-groups.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.importExportViews,
  components: {
    DraggableGroupViews,
    DraggableGroups,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    entitiesViewsGroupsMixin,
    rightsEntitiesGroupMixin,
  ],
  data() {
    const { groups, views } = this.modal.config;

    return {
      importedGroups: [...groups],
      importedViews: [...views],
      currentGroups: [],
    };
  },
  methods: {
    async updateViews(tabId, viewId) {
      if (this.config.action) {
        await this.config.action({ tabId, viewId });
      }

      this.$modals.hide();
    },
  },
};
</script>

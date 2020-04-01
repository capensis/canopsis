<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('common.import') }}/{{ $t('common.export') }}
    template(slot="text")
      v-layout(row)
        v-flex.pl-1.pr-1(xs4)
          v-flex.text-xs-center.mb-2 Groups
          draggable-groups(v-model="importedGroups")
        v-flex.pl-1.pr-1(xs4)
          v-flex.text-xs-center.mb-2 Views
          draggable-group-views(v-model="importedViews")
        v-flex.pl-1.pr-1(xs4)
          v-flex.text-xs-center.mb-2 Result
          draggable-groups(v-model="currentGroups")
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(type="submit", @click="updateViews") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { cloneDeep } from 'lodash';

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
    return {
      importedViews: [],
      importedGroups: [],
      currentGroups: [],
    };
  },
  watch: {
    groupsOrdered: {
      immediate: true,
      handler() {
        this.setDefaultValues();
      },
    },
  },
  methods: {
    async updateViews() {
      this.$modals.hide();
    },
    setDefaultValues() {
      const { groups, views } = this.modal.config;

      this.importedViews = cloneDeep(views);
      this.importedGroups = cloneDeep(groups);
      this.currentGroups = cloneDeep(this.groupsOrdered);
    },
  },
};
</script>

<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('modals.importExportViews.title') }}</span>
    </template>
    <template #text="">
      <v-layout>
        <v-flex
          class="pl-1 pr-1"
          xs4
        >
          <v-flex class="text-center mb-2">
            {{ $t('modals.importExportViews.groups') }}
          </v-flex>
          <draggable-groups
            v-model="importedGroups"
            pull
            view-pull
            view-put
            @change:group="changeImportedGroupHandler"
          />
        </v-flex>
        <v-flex
          class="pl-1"
          xs4
        >
          <v-flex class="text-center mb-2">
            {{ $t('modals.importExportViews.views') }}
          </v-flex>
          <draggable-group-views
            v-model="importedViews"
            pull
          />
        </v-flex>
        <v-divider
          class="ml-1 mr-1 secondary"
          vertical
        />
        <v-flex xs4>
          <v-flex class="text-center mb-2">
            {{ $t('common.result') }}
          </v-flex>
          <draggable-groups
            v-model="currentGroups"
            put
            pull
            view-put
            @change:group="changeCurrentGroupHandler"
          />
        </v-flex>
      </v-layout>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
        class="primary"
        :loading="submitting"
        type="submit"
        @click="submit"
      >
        {{ $t('common.saveChanges') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import {
  prepareImportedViews,
  prepareImportedGroups,
  prepareCurrentGroupsForImporting,
  prepareViewGroupsForImportRequest,
} from '@/helpers/entities/view/form';

import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { authMixin } from '@/mixins/auth';
import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesViewMixin } from '@/mixins/entities/view';
import { submittableMixinCreator } from '@/mixins/submittable';

import DraggableGroupViews from '@/components/layout/navigation/partials/groups-side-bar/draggable-group-views.vue';
import DraggableGroups from '@/components/layout/navigation/partials/groups-side-bar/draggable-groups.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.importExportViews,
  components: {
    DraggableGroupViews,
    DraggableGroups,
    ModalWrapper,
  },
  mixins: [
    authMixin,
    modalInnerMixin,
    entitiesViewMixin,
    entitiesViewGroupMixin,
    submittableMixinCreator(),
  ],
  data() {
    return {
      importedViews: [],
      importedGroups: [],
      currentGroups: [],
    };
  },
  watch: {
    groups: {
      immediate: true,
      handler() {
        this.setDefaultValues();
      },
    },
  },
  methods: {
    changeImportedGroupHandler(groupIndex, group) {
      this.importedGroups.splice(groupIndex, 1, group);
    },

    changeCurrentGroupHandler(groupIndex, group) {
      this.currentGroups.splice(groupIndex, 1, group);
    },

    setDefaultValues() {
      const { importedGroups, importedViews } = this.config;

      this.importedViews = prepareImportedViews(importedViews);
      this.importedGroups = prepareImportedGroups(importedGroups);
      this.currentGroups = prepareCurrentGroupsForImporting(this.groups.filter(group => !group.is_private));
    },

    async submit() {
      try {
        await this.importViewsWithoutStore({
          data: prepareViewGroupsForImportRequest(this.currentGroups),
        });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.description || this.$t('errors.default') });
      } finally {
        await this.fetchAllGroupsListWithWidgetsWithCurrentUser();
        this.$modals.hide();
      }
    },
  },
};
</script>

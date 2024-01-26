<template lang="pug">
  div
    v-layout.d-inline-flex(v-if="serviceEntities.length", align-center, row)
      v-checkbox-functional.ml-3.pa-0(v-model="isAllSelected", :disabled="!entitiesWithActions.length")
      v-fade-transition(mode="out-in")
        span.font-italic(v-if="!selectedEntities.length") {{ $t('modals.service.massActionsDescription') }}
        service-entity-actions(
          v-else,
          :actions="actions",
          :entity="service",
          @apply="applyActionForSelected"
        )
    div.mt-2(v-for="serviceEntity in serviceEntities", :key="serviceEntity.key")
      service-entity(
        :service-id="service._id",
        :entity="serviceEntity",
        :last-action-unavailable="hasActionError(serviceEntity)",
        :entity-name-field="entityNameField",
        :widget-parameters="widgetParameters",
        :selected="isEntitySelected(serviceEntity)",
        :actions-requests="actionsRequests",
        @update:selected="updateSelected(serviceEntity, $event)",
        @remove:unavailable="removeEntityFromUnavailable(serviceEntity)",
        @apply:action="$listeners['apply:action']",
        @refresh="$listeners.refresh"
      )
    c-table-pagination.mt-1(
      v-if="totalItems > pagination.rowsPerPage",
      :total-items="totalItems",
      :rows-per-page="pagination.rowsPerPage",
      :page="pagination.page",
      @update:page="updatePage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import {
  getAvailableActionsByEntities,
  getAvailableEntityActionsTypes,
  isActionTypeAvailableForEntity,
} from '@/helpers/entities/entity';
import { filterById, mapIds } from '@/helpers/entities';
import { getPageForNewRecordsPerPage } from '@/helpers/pagination';

import { widgetActionPanelServiceEntityMixin } from '@/mixins/widget/actions-panel/service-entity';

import ServiceEntityActions from './service-entity-actions.vue';
import ServiceEntity from './service-entity.vue';

export default {
  components: {
    ServiceEntityActions,
    ServiceEntity,
  },
  mixins: [widgetActionPanelServiceEntityMixin],
  props: {
    service: {
      type: Object,
      required: true,
    },
    pagination: {
      type: Object,
      required: true,
    },
    serviceEntities: {
      type: Array,
      default: () => [],
    },
    entityNameField: {
      type: String,
      default: 'name',
    },
    widgetParameters: {
      type: Object,
      default: () => ({}),
    },
    totalItems: {
      type: Number,
      required: false,
    },
    actionsRequests: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      selectedEntities: [],
    };
  },
  computed: {
    entitiesWithActions() {
      return this.serviceEntities.filter(entity => getAvailableEntityActionsTypes(entity).length);
    },

    isAllSelected: {
      get() {
        return this.entitiesWithActions.length > 0
          && this.entitiesWithActions.every(({ _id: id }) => this.selectedEntitiesIds.includes(id));
      },
      set(checked) {
        if (checked) {
          this.selectedEntities = [...this.entitiesWithActions];
        } else {
          this.selectedEntities = [];
        }
      },
    },

    actions() {
      return getAvailableActionsByEntities(this.selectedEntities)
        .filter(this.actionsAccessFilterHandler)
        .map(action => ({
          ...action,
          disabled: !this.hasEntityWithoutAction(action.type),
        }));
    },

    selectedEntitiesIds() {
      return mapIds(this.selectedEntities);
    },
  },
  watch: {
    serviceEntities() {
      this.selectedEntities = [];
    },
  },
  methods: {
    hasActionError(entity) {
      return this.unavailableEntitiesAction[entity._id];
    },

    hasEntityWithoutAction(type) {
      return this.selectedEntities.filter(entity => isActionTypeAvailableForEntity(type, entity));
    },

    applyActionForSelected({ type }) {
      this.applyEntityAction(type, this.selectedEntities);
    },

    updateSelected(entity, checked) {
      if (checked) {
        this.selectedEntities.push(entity);
      } else {
        this.selectedEntities = filterById(this.selectedEntities, entity);
      }
    },

    isEntitySelected(entity) {
      return this.selectedEntitiesIds.includes(entity._id);
    },

    updatePage(page) {
      this.$emit('update:pagination', { ...this.pagination, page });
    },

    updateRecordsPerPage(rowsPerPage) {
      this.$emit('update:pagination', {
        ...this.pagination,

        rowsPerPage,
        page: getPageForNewRecordsPerPage(rowsPerPage, this.pagination.rowsPerPage, this.pagination.page),
      });
    },
  },
};
</script>

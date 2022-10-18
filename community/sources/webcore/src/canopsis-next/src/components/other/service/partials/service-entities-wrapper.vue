<template lang="pug">
  div
    v-layout.d-inline-flex(v-if="serviceEntities.length", align-center, row)
      v-checkbox-functional.ml-3.pa-0(v-model="isAllSelected")
      template(v-if="selectedEntities.length")
        service-entity-actions(:actions="actions", @apply="applyActionForSelected")
    div.mt-2(v-for="serviceEntity in serviceEntities", :key="serviceEntity.key")
      service-entity(
        :service-id="service._id",
        :entity="serviceEntity",
        :last-action-unavailable="hasActionError(serviceEntity)",
        :entity-name-field="entityNameField",
        :widget-parameters="widgetParameters",
        :selected="isEntitySelected(serviceEntity)",
        @select="updateSelected(serviceEntity, $event)",
        @remove-unavailable="removeEntityFromUnavailable(serviceEntity)",
        @add:action="$listeners['add:action']",
        @refresh="$listeners.refresh"
      )
    c-table-pagination(
      v-if="totalItems > pagination.rowsPerPage",
      :total-items="totalItems",
      :rows-per-page="pagination.rowsPerPage",
      :page="pagination.page",
      @update:page="updatePage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import { getAvailableActionsByEntities, isActionTypeAvailableForEntity } from '@/helpers/entities/entity';
import { filterById, mapIds } from '@/helpers/entities';

import { widgetActionPanelServiceEntityMixin } from '@/mixins/widget/actions-panel/service-entity';

import ServiceEntityActions from './service-entity-actions.vue';
import ServiceEntity from './service-entity.vue';

export default {
  inject: ['$actionsQueue'],
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
  },
  data() {
    return {
      selectedEntities: [],
    };
  },
  computed: {
    isAllSelected: {
      get() {
        return this.serviceEntities.every(({ _id: id }) => this.selectedEntitiesIds.includes(id));
      },
      set(checked) {
        this.selectedEntities = checked
          ? [...this.serviceEntities]
          : [];
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

    pendingEntitiesIdsByActionType() {
      return this.$actionsQueue.queue
        .reduce((acc, { actionType, entities }) => {
          const entitiesIds = mapIds(entities);

          if (acc[actionType]) {
            acc[actionType].push(...entitiesIds);
          } else {
            acc[actionType] = entitiesIds;
          }

          return acc;
        }, {});
    },
  },
  methods: {
    hasActionError(entity) {
      return this.unavailableEntitiesAction[entity._id];
    },

    isEntityActionPending(type, id) {
      const pendingEntitiesIds = this.pendingEntitiesIdsByActionType[type] || [];

      return pendingEntitiesIds.includes(id);
    },

    hasEntityWithoutAction(type) {
      return this.selectedEntities
        .filter(entity => isActionTypeAvailableForEntity(type, entity))
        .some(({ _id: id }) => !this.isEntityActionPending(type, id));
    },

    applyActionForSelected({ type }) {
      const availableEntities = this.selectedEntities
        .filter(entity => !this.isEntityActionPending(type, entity._id));

      this.addEntityAction(type, availableEntities);
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
      this.$emit('update:pagination', { ...this.pagination, rowsPerPage, page: 1 });
    },
  },
};
</script>
